package pojo

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Minimal transaction interface.
//
// types.Transaction and pojo.BlocknativeTx have implemented this interface.
type TxData interface {
	Data() []byte
	Gas() uint64
	GasPrice() *big.Int
	Value() *big.Int
	Nonce() uint64
	To() *common.Address
	From() *common.Address
	Hash() common.Hash
	Source() string
}

// A *types.Transaction wrapper which implememts the TxData interface.
type RawTransaction struct {
	*types.Transaction
	source string
}

func NewRawTransaction(tx *types.Transaction, fullnodeUrl string) *RawTransaction {
	return &RawTransaction{
		Transaction: tx,
		source:      fullnodeUrl,
	}
}

func (tx *RawTransaction) From() *common.Address {
	// rawTx := (*types.Transaction)(tx)
	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
	if err != nil {
		return nil
	} else {
		from := msg.From()
		return &from
	}
}

func (tx *RawTransaction) Source() string {
	return tx.source
}
