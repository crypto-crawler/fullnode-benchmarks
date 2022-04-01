package pojo

import (
	"fmt"
	"math/big"
)

type BigInt struct {
	*big.Int
}

func NewBigInt(n *big.Int) *BigInt {
	m := big.NewInt(0).Set(n) // deep copy a big.Int
	return &BigInt{Int: m}
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}

	return []byte("\"0x" + b.Text(16) + "\""), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	s := string(p)
	if s == "null" {
		return nil
	}
	if s[0:3] != "\"0x" {
		return fmt.Errorf("%s does NOT start with \"0x", s)
	}
	s = s[3 : len(s)-1]

	z := big.NewInt(0)
	z, ok := z.SetString(s, 16)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	b.Int = z
	return nil
}
