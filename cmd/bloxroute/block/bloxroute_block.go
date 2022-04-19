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
)

// doc: https://docs.bloxroute.com/streams/bdnblocks

// Subscribe to pending blocks from the `bdnBlocks` stream of bloXroute gateway or cloud API.
func main() {
	certFile := flag.String("cert", "external_gateway_cert.pem", "The cert file")
	keyFile := flag.String("key", "external_gateway_key.pem", "The key file")
	outputFile := flag.String("output", "bloxroute-block-cloud.json", "The output file")
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

	pendingBlockCh := make(chan *types.Block)
	_, err = bloXrouteClient.SubscribeBdnBlocks([]string{"hash"}, pendingBlockCh)
	if err != nil {
		log.Fatal(err)
	}

	go utils.Run(pendingBlockCh, stopCh, *outputFile)

	<-signals
	log.Println("Ctrl+C detected, exiting...")
	close(stopCh)
	time.Sleep(1 * time.Second) // give some time for other goroutines to stop
}
