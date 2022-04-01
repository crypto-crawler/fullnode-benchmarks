package utils

import (
	"math/big"
	"testing"

	"github.com/crypto-crawler/fullnode-benchmarks/pojo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestDecodeReturnedDataOfGetReserves(t *testing.T) {
	pair := common.HexToAddress("0x58F876857a02D6762E0101bb5C46A8c1ED44Dc16")
	data := "0x0000000000000000000000000000000000000000000d9364e40e581d2dfdc52f00000000000000000000000000000000000000000000409dd0fd22cd782430f50000000000000000000000000000000000000000000000000000000062413c6d"
	pairReserve, err := DecodeReturnedDataOfGetReserves(pair, data, 16448132)
	assert.NoError(t, err)

	reserve0, ok := big.NewInt(0).SetString("d9364e40e581d2dfdc52f", 16)
	assert.True(t, ok)
	reserve1, ok := big.NewInt(0).SetString("409dd0fd22cd782430f5", 16)
	assert.True(t, ok)
	assert.Equal(t, pojo.NewBigInt(reserve0), pairReserve.Reserve0)
	assert.Equal(t, pojo.NewBigInt(reserve1), pairReserve.Reserve1)
	assert.Equal(t, uint32(1648442477), pairReserve.BlockTimestampLast)
	assert.Equal(t, int64(16448132), pairReserve.BlockNumber)
}
