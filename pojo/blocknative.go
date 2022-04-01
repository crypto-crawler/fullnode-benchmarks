package pojo

import (
	"encoding/hex"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// ConnectResponse is the message we receive when opening a connection to the API
type ConnectResponse struct {
	ConnectionID  string `json:"connectionId"`
	ServerVersion string `json:"serverVersion"`
	ShowUX        bool   `json:"showUX"`
	Status        string `json:"status"`
	Reason        string `json:"reason"`
	Version       int    `json:"version"`
}

// BaseMessage is the base message required for all interactions with the websockets api
type BaseMessage struct {
	CategoryCode string    `json:"categoryCode"`
	EventCode    string    `json:"eventCode"`
	Timestamp    time.Time `json:"timeStamp"`
	DappID       string    `json:"dappId"` // api key
	Version      string    `json:"version"`
	Blockchain   `json:"blockchain"`
}

// Blockchain is a type fulfilling the blockchain params
type Blockchain struct {
	System  string `json:"system"`
	Network string `json:"network"`
}

// Blocknative message
type BlocknativeMsg struct {
	Version       int       `json:"version"`
	ServerVersion string    `json:"serverVersion"`
	TimeStamp     time.Time `json:"timeStamp"`
	ConnectionID  string    `json:"connectionId"`
	Status        string    `json:"status"`
	Event         struct {
		BaseMessage
		Transaction struct {
			TimeStamp        time.Time `json:"timeStamp"`
			Status           string    `json:"status"`
			MonitorID        string    `json:"monitorId"`
			MonitorVersion   string    `json:"monitorVersion"`
			TimePending      string    `json:"timePending"`
			PendingTimeStamp time.Time `json:"pendingTimeStamp"`
			BlocksPending    int       `json:"blocksPending"`
			Hash             string    `json:"hash"`
			From             string    `json:"from"`
			To               string    `json:"to"`
			Value            string    `json:"value"`
			Gas              int       `json:"gas"`
			GasPrice         string    `json:"gasPrice"`
			GasPriceGwei     float64   `json:"gasPriceGwei"`
			Nonce            int       `json:"nonce"`
			BlockHash        string    `json:"blockHash"`
			BlockNumber      int       `json:"blockNumber"`
			TransactionIndex int       `json:"transactionIndex"`
			Input            string    `json:"input"`
			GasUsed          string    `json:"gasUsed"`
			Asset            string    `json:"asset"`
			WatchedAddress   string    `json:"watchedAddress"`
			Direction        string    `json:"direction"`
			Counterparty     string    `json:"counterparty"`
		} `json:"transaction"`
	} `json:"event"`
}

// Implement the TxData interface for BlocknativeMsg.

func (tx *BlocknativeMsg) Data() []byte {
	data, _ := hex.DecodeString(tx.Event.Transaction.Input[2:])
	return data
}

func (tx *BlocknativeMsg) Gas() uint64 {
	return uint64(tx.Event.Transaction.Gas)
}

func (tx *BlocknativeMsg) GasPrice() *big.Int {
	gasPrice := big.NewInt(0)
	gasPrice.SetString(tx.Event.Transaction.GasPrice, 0)
	return gasPrice
}

func (tx *BlocknativeMsg) Value() *big.Int {
	amount := big.NewInt(0)
	amount.SetString(tx.Event.Transaction.Value, 0)
	return amount
}

func (tx *BlocknativeMsg) Nonce() uint64 {
	return uint64(tx.Event.Transaction.Nonce)
}

func (tx *BlocknativeMsg) To() *common.Address {
	address := common.HexToAddress(tx.Event.Transaction.To)
	return &address
}

func (tx *BlocknativeMsg) From() *common.Address {
	address := common.HexToAddress(tx.Event.Transaction.From)
	return &address
}

func (tx *BlocknativeMsg) Hash() common.Hash {
	return common.HexToHash(tx.Event.Transaction.Hash)
}

func (tx *BlocknativeMsg) Source() string {
	return "blocknative"
}
