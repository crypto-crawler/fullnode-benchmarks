package clients

import (
	"context"
	"log"
	"math/big"
	"sync"

	"github.com/crypto-crawler/bloxroute-go/client"
	bloXrouteTypes "github.com/crypto-crawler/bloxroute-go/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Get current block number
type BlockNumber struct {
	blockNumber *big.Int
	rw          *sync.RWMutex
}

func NewBlockNumberOnBloXroute(bloXrouteClient *client.BloXrouteClient, stopCh <-chan struct{}) (*BlockNumber, error) {
	outCh := make(chan *bloXrouteTypes.EthOnBlockResponse)
	callParams := make([]map[string]string, 0)
	callParams = append(callParams, map[string]string{"name": "block_number", "method": "eth_blockNumber"})
	_, err := bloXrouteClient.SubscribeEthOnBlock(nil, callParams, outCh)
	if err != nil {
		return nil, err
	}

	rw := sync.RWMutex{}
	blockNumber := big.NewInt(0)
	go func() {
		for {
			select {
			case <-stopCh:
				close(outCh)
				return
			case resp := <-outCh:
				if resp.Name == "block_number" {
					if number, ok := big.NewInt(0).SetString(resp.Response, 0); ok {
						rw.Lock()
						blockNumber.Set(number)
						rw.Unlock()
					} else {
						log.Fatal(resp)
					}
				}
			}
		}
	}()

	// Initialize the first block number
	for resp := range outCh {
		if resp.Name == "block_number" {
			if number, ok := big.NewInt(0).SetString(resp.Response, 0); ok {
				blockNumber.Set(number)
				break
			} else {
				log.Fatal(resp)
			}
		}
	}

	return &BlockNumber{
		blockNumber: blockNumber,
		rw:          &rw,
	}, nil
}

func NewBlockNumberOnFullnode(fullnodeUrl string, stopCh <-chan struct{}) (*BlockNumber, error) {
	ctx := context.Background()
	ethClient, err := ethclient.DialContext(ctx, fullnodeUrl)
	if err != nil {
		return nil, err
	}

	headCh := make(chan *types.Header)
	sub, err := ethClient.SubscribeNewHead(ctx, headCh)
	if err != nil {
		return nil, err
	}

	blockNumber := big.NewInt(0)
	{
		header, err := ethClient.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, err
		}
		blockNumber = header.Number
	}

	rw := sync.RWMutex{}
	go func() {
	DONE:
		for {
			select {
			case <-stopCh:
				break DONE
			case err := <-sub.Err():
				log.Println(err)
				break DONE
			case head := <-headCh:
				// log.Println("New block:", head.Number.Uint64())
				rw.Lock()
				blockNumber.Set(head.Number)
				rw.Unlock()
			}
		}

		sub.Unsubscribe()
		close(headCh)
		ctx.Done()
		ethClient.Close()
	}()

	return &BlockNumber{
		blockNumber: blockNumber,
		rw:          &rw,
	}, nil
}

// Get current block number without network IO, so this method is super fast.
func (b *BlockNumber) Get() *big.Int {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return big.NewInt(0).Set(b.blockNumber) // deep copy
}
