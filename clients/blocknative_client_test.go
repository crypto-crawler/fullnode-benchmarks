package clients

import (
	"os"
	"testing"

	"github.com/crypto-crawler/fullnode-benchmarks/pojo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestBlocknativeClientPancakeSwap(t *testing.T) {
	apiKey := os.Getenv("BLOCKNATIVE_API_KEY")
	if apiKey == "" {
		assert.FailNow(t, "Please provide the blocknative API key by setting the BLOCKNATIVE_API_KEY environment variable variable")
	}

	client, err := NewBlocknativeClient(apiKey, "ethereum", "bsc-main", nil, nil)
	assert.NoError(t, err)

	stopCh := make(chan struct{})
	txCh := make(chan pojo.TxData)
	err = client.Subscribe(stopCh, txCh)
	assert.NoError(t, err)

	tx := <-txCh
	assert.NotNil(t, tx)
}

func TestBlocknativeClientToAddress(t *testing.T) {
	apiKey := os.Getenv("BLOCKNATIVE_API_KEY")
	if apiKey == "" {
		assert.FailNow(t, "Please provide the blocknative API key by setting the BLOCKNATIVE_API_KEY environment variable variable")
	}

	// binance hot wallets
	toWhiteList := make(map[common.Address]bool)
	toWhiteList[common.HexToAddress("0x8894E0a0c962CB723c1976a4421c95949bE2D4E3")] = true
	toWhiteList[common.HexToAddress("0xdccF3B77dA55107280bd850ea519DF3705D1a75a")] = true

	client, err := NewBlocknativeClient(apiKey, "ethereum", "bsc-main", toWhiteList, nil)
	assert.NoError(t, err)

	stopCh := make(chan struct{})
	txCh := make(chan pojo.TxData)
	err = client.Subscribe(stopCh, txCh)
	assert.NoError(t, err)

	tx := <-txCh
	assert.NotNil(t, tx)
}
