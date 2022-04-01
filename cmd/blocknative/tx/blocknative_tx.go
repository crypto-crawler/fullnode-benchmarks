package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crypto-crawler/fullnode-benchmarks/clients"
	"github.com/crypto-crawler/fullnode-benchmarks/pojo"
	"github.com/crypto-crawler/fullnode-benchmarks/utils"
)

// Subscribe to pending transactions from blocknative.com
func main() {
	apiKey := flag.String("apikey", "", "blocknative API key")
	outputFile := flag.String("output", "blocknative-tx.json", "The output file")
	flag.Parse()
	if *apiKey == "" || *outputFile == "" {
		flag.Usage()
		return
	}

	client, err := clients.NewBlocknativeClient(*apiKey, "ethereum", "bsc-main", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	// catch Ctrl+C
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	stopCh := make(chan struct{})

	pendingTxCh := make(chan pojo.TxData)
	err = client.Subscribe(stopCh, pendingTxCh)
	if err != nil {
		log.Fatal(err)
	}

	go utils.Run(pendingTxCh, stopCh, *outputFile)

	<-signals
	log.Println("Ctrl+C detected, exiting...")
	close(stopCh)
	time.Sleep(1 * time.Second) // give some time for other goroutines to stop
}
