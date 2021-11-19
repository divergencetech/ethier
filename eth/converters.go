package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

// Ether returns numerator/denominator ETH in wei.
func Ether(numerator, denominator int64) *big.Int {
	e := new(big.Int).Mul(big.NewInt(numerator), big.NewInt(params.Ether))
	return e.Div(e, big.NewInt(denominator))
}
