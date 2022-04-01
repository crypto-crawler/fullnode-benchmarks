package pojo

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNull(t *testing.T) {
	var n *BigInt = nil
	json, err := n.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "null", string(json))

	err = n.UnmarshalJSON([]byte("null"))
	assert.NoError(t, err)
}

func TestValidNumber(t *testing.T) {
	n := NewBigInt(big.NewInt(0))

	json, err := n.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "\"0x0\"", string(json))

	m, _ := big.NewInt(0).SetString("20d25b9dc731ea1e7a9", 16)
	n = NewBigInt(m)
	assert.NoError(t, err)
	assert.NotNil(t, n)

	json, err = n.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "\"0x20d25b9dc731ea1e7a9\"", string(json))
}
