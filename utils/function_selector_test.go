package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureToFunctionSelector(t *testing.T) {
	methodId := SignatureToFunctionSelector("addLiquidityETH(address,uint256,uint256,uint256,address,uint256)")
	// see https://www.4byte.directory/signatures/?bytes4_signature=0xf305d719
	assert.Equal(t, "0xf305d719", methodId.String())
	jsonStr, _ := methodId.MarshalJSON()
	assert.Equal(t, "\"0xf305d719\"", string(jsonStr))

	methodId = SignatureToFunctionSelector("getReserves()")
	// see https://www.4byte.directory/signatures/?bytes4_signature=0x0902f1ac
	assert.Equal(t, "0x0902f1ac", methodId.String())
	jsonStr, _ = methodId.MarshalJSON()
	assert.Equal(t, "\"0x0902f1ac\"", string(jsonStr))
}
