package pojo

import (
	"crypto/md5"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type PairReserve struct {
	// Index              int64          `json:"index"`
	Pair common.Address `json:"pair"`
	// Token0             common.Address `json:"token0"`
	// Token1             common.Address `json:"token1"`
	Reserve0           *BigInt `json:"reserve0"`
	Reserve1           *BigInt `json:"reserve1"`
	BlockTimestampLast uint32  `json:"block_timestamp_last"`
	BlockNumber        int64   `json:"block_number"`
	// CreatedAt          time.Time `json:"created_at"`
}

func (p *PairReserve) Hash() uint64 {
	h := md5.New()
	h.Write(p.Pair.Bytes())
	h.Write(p.Reserve0.Bytes())
	h.Write(p.Reserve1.Bytes())
	{
		n := big.NewInt(0)
		n.SetInt64(p.BlockNumber)
		h.Write(n.Bytes())
	}
	bs := h.Sum(nil)

	return big.NewInt(0).SetBytes(bs).Uint64()
}
