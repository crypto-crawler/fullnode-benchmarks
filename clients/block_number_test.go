package clients

import (
	"os"
	"testing"
	"time"

	"github.com/crypto-crawler/bloxroute-go/client"
	"github.com/stretchr/testify/assert"
)

func TestBlockNumberOnBloXroute(t *testing.T) {
	certFile := os.Getenv("BLOXROUTE_CERT_FILE")
	keyFile := os.Getenv("BLOXROUTE_KEY_FILE")
	if certFile == "" || keyFile == "" {
		assert.FailNow(t, "Please provide the bloXroute cert and key files path in the environment variable variable")
	}

	stopCh := make(chan struct{})
	client, err := client.NewBloXrouteClientToCloud("BSC-Mainnet", certFile, keyFile, stopCh)
	assert.NoError(t, err)

	blockNumber, err := NewBlockNumberOnBloXroute(client, stopCh)
	assert.NoError(t, err)

	number1 := blockNumber.Get()
	time.Sleep(time.Second * 4)
	number2 := blockNumber.Get()

	assert.Greater(t, number2.Uint64(), number1.Uint64())

	close(stopCh)
}

func TestBlockNumberOnFullnode(t *testing.T) {
	fullnodeUrl := os.Getenv("FULLNODE_URL")
	if fullnodeUrl == "" {
		assert.FailNow(t, "Please provide the fullnode URL in the FULLNODE_URL environment variable")
	}

	stopCh := make(chan struct{})
	blockNumber, err := NewBlockNumberOnFullnode(fullnodeUrl, stopCh)
	assert.NoError(t, err)

	number1 := blockNumber.Get()
	time.Sleep(time.Second * 4)
	number2 := blockNumber.Get()

	assert.Greater(t, number2.Uint64(), number1.Uint64())

	close(stopCh)
}
