package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crypto-crawler/bloxroute-go/client"
	"github.com/crypto-crawler/bloxroute-go/types"
	"github.com/crypto-crawler/fullnode-benchmarks/utils"
	"github.com/ethereum/go-ethereum/common"
)

// Subscribe to pending transactions from the `newTxs` stream of bloXroute gateway or cloud API.
func main() {
	certFile := flag.String("cert", "external_gateway_cert.pem", "The cert file")
	keyFile := flag.String("key", "external_gateway_key.pem", "The key file")
	outputFile := flag.String("output", "bloxroute-pair-reserve-cloud.json", "The output file")
	gatewayUrl := flag.String("gateway", "", "The gateway url")
	header := flag.String("header", "", "The authorization header")
	flag.Parse()
	if *outputFile == "" {
		log.Println("-output is empty!")
		flag.Usage()
		return
	}
	if *gatewayUrl == "" {
		if *certFile == "" || *keyFile == "" {
			log.Println("If --gateway is absent, -cert and -key must be present!")
			flag.Usage()
			return
		}
	} else {
		if *header == "" {
			log.Println("-header must be present if -gateway is present!")
			flag.Usage()
			return
		}
	}

	// catch Ctrl+C
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	stopCh := make(chan struct{})

	var bloXrouteClient *client.BloXrouteClient = nil
	var err error = nil
	if *gatewayUrl == "" {
		log.Println("Connecting to bloXroute cloud")
		bloXrouteClient, err = client.NewBloXrouteClientToCloud("BSC-Mainnet", *certFile, *keyFile, stopCh)
	} else {
		log.Println("Connecting to bloXroute gateway")
		bloXrouteClient, err = client.NewBloXrouteClientToGateway(*gatewayUrl, *header, stopCh)
	}
	if err != nil {
		log.Fatal(err)
	}
	bloXrouteClientEx := client.NewBloXrouteClientExtended(bloXrouteClient, stopCh)

	outCh := make(chan *types.PairReserves)
	pairs := []common.Address{
		common.HexToAddress("0x58f876857a02d6762e0101bb5c46a8c1ed44dc16"),
		common.HexToAddress("0x7efaef62fddcca950418312c6c91aef321375a00"),
		common.HexToAddress("0x0ed7e52944161450477ee417de9cd3a859b14fd0"),
		common.HexToAddress("0x16b9a82891338f9ba80e2d6970fdda79d1eb0dae"),
		common.HexToAddress("0x2354ef4df11afacb85a5c7f98b624072eccddbb1"),
	}
	err = bloXrouteClientEx.SubscribePairReserves(pairs, outCh)
	if err != nil {
		log.Fatal(err)
	}

	go utils.Run(outCh, stopCh, *outputFile)

	<-signals
	log.Println("Ctrl+C detected, exiting...")
	close(stopCh)
	time.Sleep(1 * time.Second) // give some time for other goroutines to stop
}
