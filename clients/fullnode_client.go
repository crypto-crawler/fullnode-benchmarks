package clients

import (
	"context"
	"log"

	"github.com/crypto-crawler/fullnode-benchmarks/abi"
	"github.com/crypto-crawler/fullnode-benchmarks/pojo"
	"github.com/crypto-crawler/fullnode-benchmarks/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/fxfactorial/defi-abigen/contracts/uniswap/pair"
)

func SubscribePendingTxHash(fullNodeUrl string, stopCh <-chan struct{}) (<-chan common.Hash, error) {
	ctx := context.Background()
	rpcClient, err := rpc.DialContext(ctx, fullNodeUrl)
	if err != nil {
		return nil, err
	}
	gethClient := gethclient.New(rpcClient)

	txHashCh := make(chan common.Hash)

	sub, err := gethClient.SubscribePendingTransactions(ctx, txHashCh)
	if err != nil {
		return nil, err
	}

	go func() {
		<-stopCh
		sub.Unsubscribe()
		rpcClient.Close()
		ctx.Done()
		close(txHashCh)
	}()

	return txHashCh, nil
}

// Subscribe pending transactions from the fullnode.
func SubscribePendingTx(fullNodeUrl string, fromWhiteList map[common.Address]bool, toWhiteList map[common.Address]bool, stopCh <-chan struct{}, txCh chan<- pojo.TxData) error {
	ctx := context.Background()
	rpcClient, err := rpc.DialContext(ctx, fullNodeUrl)
	if err != nil {
		return err
	}

	ethClient, err := ethclient.DialContext(ctx, fullNodeUrl)
	if err != nil {
		return err
	}
	gethClient := gethclient.New(rpcClient)

	txHashCh := make(chan common.Hash, 1024)
	sub, err := gethClient.SubscribePendingTransactions(ctx, txHashCh)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			select {
			case <-stopCh:
				sub.Unsubscribe()
				rpcClient.Close()
				return
			case txnHash := <-txHashCh:
				go func() {
					tx, isPending, err := utils.TransactionByHashWithRetry(ethClient, txnHash, 11)
					if err != nil {
						// Usually happens when eth.syncing is not false
						// log.Printf("TransactionByHashWithRetry(%s) failed, error: %v", txnHash.Hex(), err)
						return
					}
					// only care about pending transactions
					if tx == nil || !isPending {
						return
					}
					// Only care about transactions that interact with smart contracts
					interactWithContract := tx.To() != nil && len(tx.Data()) > 0
					if !interactWithContract {
						return
					}

					txData := pojo.TxData(pojo.NewRawTransaction(tx, fullNodeUrl))
					// if both are empty, there is no filtering at all
					if len(toWhiteList) == 0 && len(fromWhiteList) == 0 {
						txCh <- txData
					} else if len(toWhiteList) > 0 && toWhiteList[*tx.To()] {
						// transactions sent to addresses in `toWhiteList`
						txCh <- txData
					} else if len(fromWhiteList) > 0 {
						// or transactions sent from addresses in `toWhiteList`
						msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
						if err == nil {
							if fromWhiteList[msg.From()] {
								txCh <- txData
							}
						}
					}
				}()
			}
		}
	}()

	return nil
}

func SubscribeNewHead(fullNodeUrl string, stopCh <-chan struct{}) (<-chan *types.Header, error) {
	ctx := context.Background()
	ethClient, err := ethclient.DialContext(ctx, fullNodeUrl)
	if err != nil {
		return nil, err
	}

	headerCh := make(chan *types.Header)

	sub, err := ethClient.SubscribeNewHead(ctx, headerCh)
	if err != nil {
		return nil, err
	}

	go func() {
		<-stopCh
		sub.Unsubscribe()
		ethClient.Close()
		ctx.Done()
		close(headerCh)
	}()

	return headerCh, nil
}

func SubscribeBlockHash(fullNodeUrl string, stopCh <-chan struct{}) (<-chan common.Hash, error) {
	blockHashCh := make(chan common.Hash)
	headerCh, err := SubscribeNewHead(fullNodeUrl, stopCh)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-stopCh:
				close(blockHashCh)
				return
			case header := <-headerCh:
				blockHashCh <- header.Hash()
			}
		}
	}()

	return blockHashCh, nil
}

// Poll GetReserves() periodically from the fullnode.
func PullPairReserves(fullNodeUrl string, pairs []common.Address, stopCh <-chan struct{}) (<-chan *pojo.PairReserve, error) {
	ctx := context.Background()
	ethClient, err := ethclient.DialContext(ctx, fullNodeUrl)
	if err != nil {
		return nil, err
	}

	blockNumber, err := NewBlockNumberOnFullnode(fullNodeUrl, stopCh)
	if err != nil {
		return nil, err
	}

	pairInstances := make([]*pair.Pair, 0)
	for _, pairAddress := range pairs {
		pairInstance, err := pair.NewPair(pairAddress, ethClient)
		if err != nil {
			return nil, err
		}
		pairInstances = append(pairInstances, pairInstance)
	}

	outCh := make(chan *pojo.PairReserve)

	go func() {
		visited := make(map[uint64]bool)
		for {
			select {
			case <-stopCh:
				close(outCh)
				return
			default:
				for i, pairInstance := range pairInstances {
					ret, err := pairInstance.GetReserves(nil)
					if err != nil {
						panic(err)
					}

					pairReserve := &pojo.PairReserve{
						Pair:               pairs[i],
						Reserve0:           pojo.NewBigInt(ret.Reserve0),
						Reserve1:           pojo.NewBigInt(ret.Reserve1),
						BlockTimestampLast: ret.BlockTimestampLast,
						BlockNumber:        blockNumber.Get().Int64(),
					}
					hash := pairReserve.Hash()
					if !visited[hash] {
						outCh <- pairReserve
						visited[hash] = true
					}
				}
			}
		}
	}()

	return outCh, nil
}

// BulkReader
func PullPairReservesBulk(fullNodeUrl string, pairs []common.Address, stopCh <-chan struct{}) (<-chan *pojo.PairReserve, error) {
	ctx := context.Background()
	ethClient, err := ethclient.DialContext(ctx, fullNodeUrl)
	if err != nil {
		return nil, err
	}

	blockNumber, err := NewBlockNumberOnFullnode(fullNodeUrl, stopCh)
	if err != nil {
		return nil, err
	}

	router := common.HexToAddress("0x45974B68d81Be55E71F7ACD5c1378a9d52CF02Be")
	bulkReader, err := abi.NewBulkReader(router, ethClient)
	if err != nil {
		return nil, err
	}

	outCh := make(chan *pojo.PairReserve)

	go func() {
		visited := make(map[uint64]bool)
		for {
			select {
			case <-stopCh:
				close(outCh)
				return
			default:
				arr, err := bulkReader.GetReservesForBenchmark(nil, pairs)
				if err != nil {
					panic(err)
				}
				for i := 0; i < len(pairs); i++ {
					pairReserve := &pojo.PairReserve{
						Pair:               pairs[i],
						Reserve0:           pojo.NewBigInt(arr[i][0]),
						Reserve1:           pojo.NewBigInt(arr[i][1]),
						BlockTimestampLast: uint32(arr[i][2].Int64()),
						BlockNumber:        blockNumber.Get().Int64(),
					}
					hash := pairReserve.Hash()
					if !visited[hash] {
						outCh <- pairReserve
						visited[hash] = true
					}
				}
			}
		}
	}()

	return outCh, nil
}

// BulkReader + header
func PullPairReservesBulkHeader(fullNodeUrl string, pairs []common.Address, stopCh <-chan struct{}) (<-chan *pojo.PairReserve, error) {
	ctx := context.Background()
	ethClient, err := ethclient.DialContext(ctx, fullNodeUrl)
	if err != nil {
		return nil, err
	}

	blockNumber, err := NewBlockNumberOnFullnode(fullNodeUrl, stopCh)
	if err != nil {
		return nil, err
	}

	headerCh, err := SubscribeNewHead(fullNodeUrl, stopCh)
	if err != nil {
		return nil, err
	}

	router := common.HexToAddress("0x45974B68d81Be55E71F7ACD5c1378a9d52CF02Be")
	bulkReader, err := abi.NewBulkReader(router, ethClient)
	if err != nil {
		return nil, err
	}

	outCh := make(chan *pojo.PairReserve)

	go func() {
		visited := make(map[uint64]bool)
		for {
			select {
			case <-stopCh:
				close(outCh)
				return
			case <-headerCh:
				arr, err := bulkReader.GetReservesForBenchmark(nil, pairs)
				if err != nil {
					panic(err)
				}
				for i := 0; i < len(pairs); i++ {
					pairReserve := &pojo.PairReserve{
						Pair:               pairs[i],
						Reserve0:           pojo.NewBigInt(arr[i][0]),
						Reserve1:           pojo.NewBigInt(arr[i][1]),
						BlockTimestampLast: uint32(arr[i][2].Int64()),
						BlockNumber:        blockNumber.Get().Int64(),
					}
					hash := pairReserve.Hash()
					if !visited[hash] {
						outCh <- pairReserve
						visited[hash] = true
					}
				}
			}
		}
	}()

	return outCh, nil
}
