package utils

import (
	"bufio"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/crypto-crawler/fullnode-benchmarks/pojo"
	"github.com/ethereum/go-ethereum/common"
)

func Run[T any](inputCh <-chan T, stopCh <-chan struct{}, outputFile string) {
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bw := bufio.NewWriterSize(file, 32*1024) // 32KB buffer
	defer bw.Flush()

	mu := sync.Mutex{} // used between WriteString() and Flush()

	outputCh := make(chan string, 65536)
	go func() {
		for txt := range outputCh {
			mu.Lock()
			bw.WriteString(txt + "\n")
			mu.Unlock()
		}
	}()

	ticker := time.NewTicker(time.Second) // flush per second
	go func() {
		// Writing to disk is done in a separate goroutine to avoid blocking the main thread,
		// so that the received_at field is precise.
		for range ticker.C {
			mu.Lock()
			bw.Flush()
			mu.Unlock()
		}
	}()

	for x := range inputCh {
		bytes, _ := json.Marshal(x)
		jsonMap := make(map[string]interface{})
		json.Unmarshal(bytes, &jsonMap)
		jsonMap["received_at"] = time.Now().UnixMilli()
		bytes, _ = json.Marshal(jsonMap)
		outputCh <- string(bytes)
	}

	ticker.Stop()
	close(outputCh)
}

// Decode the data of ethOnBlock of GetReserves()
func DecodeReturnedDataOfGetReserves(pair common.Address, hexStr string, blockNumber int64) (*pojo.PairReserve, error) {
	bytes, err := hex.DecodeString(hexStr[2:])
	if err != nil {
		return nil, err
	}

	pairReserve := &pojo.PairReserve{
		Pair:               pair,
		Reserve0:           pojo.NewBigInt(big.NewInt(0).SetBytes(bytes[0:32])),
		Reserve1:           pojo.NewBigInt(big.NewInt(0).SetBytes(bytes[32:64])),
		BlockTimestampLast: uint32(big.NewInt(0).SetBytes(bytes[64:]).Int64()),
		BlockNumber:        blockNumber,
	}
	return pairReserve, nil
}

func ReadPairs(pairFile string) ([]common.Address, error) {
	file, err := os.Open(pairFile)
	if err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer gz.Close()

	pairs := make([]common.Address, 0)
	scanner := bufio.NewScanner(gz)
	for scanner.Scan() {
		line := scanner.Text()
		pairs = append(pairs, common.HexToAddress(line))
	}
	return pairs, nil
}
