package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crypto-crawler/fullnode-benchmarks/clients"
	"github.com/crypto-crawler/fullnode-benchmarks/utils"
	"github.com/ethereum/go-ethereum/common"
)

// Subscribe to pending transactions from a standard fullnode.
func main() {
	fullNodeUrl := flag.String("fullnode", os.Getenv("FULLNODE_URL"), "The fullnode URL")
	outputFile := flag.String("output", "fullnode-block.json", "The output file")
	flag.Parse()
	if *fullNodeUrl == "" || *outputFile == "" {
		flag.Usage()
		return
	}

	// catch Ctrl+C
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	stopCh := make(chan struct{})

	blockHashCh, err := clients.SubscribeBlockHash(*fullNodeUrl, stopCh)
	if err != nil {
		log.Fatal(err)
	}

	jsonCh := make(chan map[string]common.Hash)
	go func() {
		// put into a map
		for blockHash := range blockHashCh {
			jsonMap := make(map[string]common.Hash)
			jsonMap["hash"] = blockHash
			jsonCh <- jsonMap
		}
		close(jsonCh)
	}()

	go utils.Run(jsonCh, stopCh, *outputFile)

	<-signals
	log.Println("Ctrl+C detected, exiting...")
	close(stopCh)
	time.Sleep(1 * time.Second) // give some time for other goroutines to stop
}
