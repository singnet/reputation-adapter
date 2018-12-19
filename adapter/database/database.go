package database

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"os"

	"database/sql"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

//ChannelLog type
type ChannelLog struct {
	DB        *sql.DB
	Log       []*Channel
	Index     map[*big.Int]int
	LastBlock uint64
}

//Channel type
type Channel struct {
	ChannelId   *big.Int
	Nonce       *big.Int
	Sender      common.Address
	Recipient   common.Address
	ClaimAmount *big.Int
	OpenTime    int64
	CloseTime   int64
}

//GetByTimeRange func
func (l *ChannelLog) GetByTimeRange(start int64, end int64) {
	queryStmt, err := l.DB.Prepare("select * from channel where is_closed=true and close_time between $1 and $2 order by close_time desc;")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := queryStmt.Query(start, end)
	defer rows.Close()

	var (
		channelId int64
		nonce     int64
		sender    string
		recipient string
		amount    int64
		openTime  int64
		closeTime int64
		isClosed  bool
	)
	for rows.Next() {

		err := rows.Scan(&channelId, &nonce, &sender, &recipient, &amount, &openTime, &closeTime, &isClosed)
		if err != nil {
			log.Fatal(err)
		}

		nextChannel := &Channel{
			big.NewInt(channelId),
			big.NewInt(nonce),
			common.HexToAddress(sender),
			common.HexToAddress(recipient),
			big.NewInt(amount),
			openTime,
			closeTime,
		}

		l.Append(nextChannel, uint64(closeTime))

	}

	return
}

//Insert is a func
func (l *ChannelLog) Insert(nc *Channel, dryRun bool) error {

	const qry = "INSERT INTO channel (channel_id, nonce, sender, recipient, amount, open_time, close_time) VALUES ($1,$2,$3,$4,$5,$6,$7);"

	if !dryRun {
		_, err := l.DB.Exec(qry, nc.ChannelId.Int64(), nc.Nonce.Int64(), nc.Sender.String(), nc.Recipient.String(), nc.ClaimAmount.Int64(), nc.OpenTime, nc.CloseTime)
		if err != nil {
			return err
		}
	}

	fmt.Printf("INSERT INTO channel (channel_id, nonce, sender, recipient, amount, open_time, close_time) VALUES (%s,%s,%s,%s,%s,%v,%v);\n",
		nc.ChannelId, nc.Nonce, nc.Sender.String(), nc.Recipient.String(), nc.ClaimAmount, nc.OpenTime, nc.CloseTime)

	return nil

}

//Update is a func
func (l *ChannelLog) Update(nc *Channel, dryRun bool) error {
	const qry = "UPDATE channel SET amount = $1, close_time = $2, is_closed = true WHERE channel_id = $3 AND nonce = $4 AND is_closed <> true;"

	if !dryRun {
		rows, err := l.DB.Exec(qry, nc.ClaimAmount.Int64(), nc.CloseTime, nc.ChannelId.Int64(), nc.Nonce.Int64()-1)

		if err != nil {
			return err
		}

		rowsAffected, err := rows.RowsAffected()
		if rowsAffected < 1 {
			return errors.New("Already present in local storage")
		}
	}

	fmt.Printf("UPDATE channel SET amount = %s, close_time = %v, is_closed = true WHERE channel_id = %v AND nonce = %v; \n", nc.ClaimAmount, nc.CloseTime, nc.ChannelId, nc.Nonce)

	return nil
}

//New function
func (l *ChannelLog) New() {

	user := os.Getenv("POSTGRES_USR")
	password := os.Getenv("POSTGRES_PSW")
	dbname := os.Getenv("POSTGRES_DB")
	if user == "" || password == "" || dbname == "" {
		log.Fatal("Missing POSTGRES enviroment variables")
	}
	var buffer bytes.Buffer
	buffer.WriteString("user=")
	buffer.WriteString(user)
	buffer.WriteString(" password=")
	buffer.WriteString(password)
	buffer.WriteString(" dbname=")
	buffer.WriteString(dbname)
	buffer.WriteString(" sslmode=disable")
	connStr := buffer.String()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		err = errors.Wrapf(err, "Couldn't ping postgre database")
		return
	}

	l.DB = db
	l.Log = []*Channel{}
	l.Index = make(map[*big.Int]int)
	l.LastBlock = uint64(0)

	return
}

//Append function
func (l *ChannelLog) Append(nextChannel *Channel, blockNumber uint64) {
	channelID := nextChannel.ChannelId
	position := len(l.Log)
	l.Log = append(l.Log, nextChannel)
	l.Index[channelID] = position
	l.LastBlock = blockNumber
}
