package crawler

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/singnet/reputation-adapter/adapter/database"
	"github.com/singnet/reputation-adapter/resources/contracts/mpe"
)

//Network config struct
type Network struct {
	RPCEndpoint     string
	DeployedAddress common.Address
	startingBlock   int64
	endingBlock     int64
}

//Escrow is escrow struct
type Escrow struct {
	ContractName     string
	ABI              abi.ABI
	ethClient        *ethclient.Client
	rawClient        *rpc.Client
	startingBlock    int64
	DeployedAddress  common.Address
	ContractInstance *mpe.Mpe
}

var networks = map[string]Network{
	"kovan": Network{
		"https://kovan.infura.io",
		common.HexToAddress("0x39f31ac7b393fe2c6660b95b878feb16ea8f3156"),
		9424242,
		0,
	},
	"ropsten": Network{
		"https://ropsten.infura.io",
		common.HexToAddress("0xAF5e3b8CF89815F24A12D45D4758D87257249778"),
		4429391,
		0,
	},
}

var channelLog = &database.ChannelLog{}

//New is a function to populate an instance of an Escrow adapter
func (e *Escrow) New(networkKey string) error {
	// Ethereum client
	currentNetwork := networks[networkKey]

	// Setup ethereum client
	if client, err := rpc.Dial(currentNetwork.RPCEndpoint); err != nil {
		return errors.Wrap(err, "error creating RPC client")
	} else {
		e.rawClient = client
		e.ethClient = ethclient.NewClient(client)
	}

	abiDefinition, err := abi.JSON(strings.NewReader(string(mpe.MpeABI)))
	if err != nil {
		return err
	}

	contractInstance, err := mpe.NewMpe(currentNetwork.DeployedAddress, e.ethClient)
	if err != nil {
		log.Fatal(err)
	}

	e.ContractName = "MultiPartyEscrow"
	e.ABI = abiDefinition
	e.startingBlock = currentNetwork.startingBlock
	e.DeployedAddress = currentNetwork.DeployedAddress
	e.ContractInstance = contractInstance

	return nil
}

//Start func
func (e *Escrow) Start() {
	//Start db
	channelLog.New()
	logs := e.getPastEvents(0)
	e.update(logs)

	//NOTICE: Will be removed if go-ethereum works with Kobvan in the future
	//https://github.com/ethereum/go-ethereum/pull/18166

	//channelLog.GetAll()

	/* 	lastBlockHex, err := e.CurrentBlockNumber()
	   	if err != nil {
	   		log.Fatal(err)
	   	}

	   	lastBlock, err := strconv.ParseUint(lastBlockHex[2:], 16, 64)
	   	if err != nil {
	   		log.Fatal(err)
	   	} */

	// Compare the block checkpoint in local database
	/* if lastBlock > channelLog.LastBlock {
	logs := e.getPastEvents(channelLog.LastBlock)
	} */
}

//GetInfo is a func
func (e *Escrow) GetInfo() {
	fmt.Println("Closed channels ", len(channelLog.Log))
	fmt.Println("Last Ethereum block ", channelLog.LastBlock)
}

//GetPastEvents func
func (e *Escrow) getPastEvents(startingBlock uint64) []types.Log {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startingBlock)),
		Addresses: []common.Address{e.DeployedAddress},
	}

	logs, err := e.ethClient.FilterLogs(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	return logs
}

func (e *Escrow) update(logs []types.Log) {
	for _, vLog := range logs {
		openTime := e.startingBlock
		closeTime := int64(vLog.BlockNumber)

		switch vLog.Topics[0].Hex() {
		case channelClaimSigHash.Hex():
			var channelClaimEvent ChannelClaim

			err := e.ABI.Unpack(&channelClaimEvent, "ChannelClaim", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			channelClaimEvent.ChannelId = vLog.Topics[1].Big()
			channelClaimEvent.Recipient = common.HexToAddress(vLog.Topics[2].Hex())

			channel, err := e.ContractInstance.Channels(nil, channelClaimEvent.ChannelId)

			if err != nil {
				log.Fatal(err)
			}

			nextChannel := &database.Channel{
				channelClaimEvent.ChannelId,
				channelClaimEvent.Nonce,
				channel.Sender,
				channelClaimEvent.Recipient,
				channelClaimEvent.ClaimAmount,
				openTime,
				closeTime,
			}

			channelLog.Update(nextChannel, false)

		case channelSenderClaimSigHash.Hex():
			var channelSenderClaimEvent ChannelSenderClaim
			err := e.ABI.Unpack(&channelSenderClaimEvent, "ChannelSenderClaim", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			channelSenderClaimEvent.ChannelId = vLog.Topics[1].Big()
			channel, err := e.ContractInstance.Channels(nil, channelSenderClaimEvent.ChannelId)
			if err != nil {
				log.Fatal(err)
			}

			nextChannel := &database.Channel{
				channelSenderClaimEvent.ChannelId,
				channelSenderClaimEvent.Nonce,
				channel.Sender,
				channel.Recipient,
				channelSenderClaimEvent.ClaimAmount,
				openTime,
				closeTime,
			}

			channelLog.Update(nextChannel, false)

		case channelExtendSigHash.Hex():
		case channelAddFundsSigHash.Hex():
		case channelOpenSigHash.Hex():
			var channelOpenEvent ChannelOpen
			err := e.ABI.Unpack(&channelOpenEvent, "ChannelOpen", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			channelOpenEvent.Sender = common.HexToAddress(vLog.Topics[1].Hex())
			channelOpenEvent.Recipient = common.HexToAddress(vLog.Topics[2].Hex())

			nextChannel := &database.Channel{
				channelOpenEvent.ChannelId,
				channelOpenEvent.Nonce,
				channelOpenEvent.Sender,
				channelOpenEvent.Recipient,
				big.NewInt(0),
				openTime,
				0,
			}

			channelLog.Insert(nextChannel, false)

		}

	}
}

// CurrentBlockNumber return current block
func (e *Escrow) CurrentBlockNumber() (currentBlockNumberHex string, err error) {
	// We have to do a raw call because the standard method of ethClient.HeaderByNumber(ctx, nil) errors on
	// unmarshaling the response currently. See https://github.com/ethereum/go-ethereum/issues/3230
	if err = e.rawClient.CallContext(context.Background(), &currentBlockNumberHex, "eth_blockNumber"); err != nil {
		return "", fmt.Errorf("error determining current block: %v", err)
	}

	return
}

// GetTimestampByBlockNumber is a func
func (e *Escrow) GetTimestampByBlockNumber(currentBlockNumberHex string) (timestamp int64, err error) {

	var currentBlockResponse struct{ Timestamp string }
	if err = e.rawClient.CallContext(context.Background(), &currentBlockResponse, "eth_getBlockByNumber", currentBlockNumberHex, true); err != nil {
		return 0, fmt.Errorf("error determining current block: %v", err)
	}
	timestamp, _ = strconv.ParseInt(currentBlockResponse.Timestamp[2:], 16, 64)

	return
}
