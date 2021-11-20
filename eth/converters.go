package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

// EtherFraction returns numerator/denominator ETH in Wei.
func EtherFraction(numerator, denominator int64) *big.Int {
	e := new(big.Int).Mul(big.NewInt(numerator), big.NewInt(params.Ether))
	return e.Div(e, big.NewInt(denominator))
}

// Ether returns e in Wei.
func Ether(e int64) *big.Int {
	return EtherFraction(e, 1)
}
