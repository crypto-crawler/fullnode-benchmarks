package pojo

import (
	"crypto/md5"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type PairReserve struct {
	Pair               common.Address `json:"pair"`
	Reserve0           *BigInt        `json:"reserve0"`
	Reserve1           *BigInt        `json:"reserve1"`
	BlockTimestampLast uint32         `json:"block_timestamp_last"`
	BlockNumber        int64          `json:"block_number"`
}

func (p *PairReserve) Hash() uint64 {
	h := md5.New()
	h.Write(p.Pair.Bytes())
	h.Write(p.Reserve0.Bytes())
	h.Write(p.Reserve1.Bytes())
	{
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, p.BlockTimestampLast)
		h.Write(b)
	}
	{
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(p.BlockNumber))
		h.Write(b)
	}
	bs := h.Sum(nil)

	return big.NewInt(0).SetBytes(bs).Uint64()
}
