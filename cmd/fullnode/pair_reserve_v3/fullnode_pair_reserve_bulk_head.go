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

// Use PullPairReservesBulk().
func main() {
	fullNodeUrl := flag.String("fullnode", os.Getenv("FULLNODE_URL"), "The fullnode URL")
	outputFile := flag.String("output", "fullnode-pair-reserve-bulk-header.json", "The output file")
	flag.Parse()
	if *fullNodeUrl == "" || *outputFile == "" {
		flag.Usage()
		return
	}

	// catch Ctrl+C
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	stopCh := make(chan struct{})

	pairs := []common.Address{
		common.HexToAddress("0x58f876857a02d6762e0101bb5c46a8c1ed44dc16"),
		common.HexToAddress("0x7efaef62fddcca950418312c6c91aef321375a00"),
		common.HexToAddress("0x0ed7e52944161450477ee417de9cd3a859b14fd0"),
		common.HexToAddress("0x16b9a82891338f9ba80e2d6970fdda79d1eb0dae"),
		common.HexToAddress("0x2354ef4df11afacb85a5c7f98b624072eccddbb1"),
	}

	pairReserveCh, err := clients.PullPairReservesBulkHeader(*fullNodeUrl, pairs, stopCh)
	if err != nil {
		log.Fatal(err)
	}

	go utils.Run(pairReserveCh, stopCh, *outputFile)

	<-signals
	log.Println("Ctrl+C detected, exiting...")
	close(stopCh)
	time.Sleep(1 * time.Second) // give some time for other goroutines to stop
}
