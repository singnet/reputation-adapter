package main

import (
	"flag"
	"log"

	crawler "github.com/singnet/reputation-adapter/adapter/crawler"
	"github.com/singnet/reputation-adapter/adapter/server"
)

func main() {

	//Flags config
	networkKey := flag.String("network", "kovan", "network. One of {mainnet, ropsten, kovan}")

	//New Escrow crawler Instance
	c := &crawler.Escrow{}
	err := c.New(*networkKey)
	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	// Grpc Server
	server.Start()

}
