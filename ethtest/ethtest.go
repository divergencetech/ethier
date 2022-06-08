// The ethtest package provides helpers for testing Ethereum smart contracts.
package ethtest

import (
	"math/big"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/google/go-cmp/cmp"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/divergencetech/ethier/eth"
)

// ExecutionErrData checks if err is both an rpc.Error and rpc.DataError, and
// returns err.ErrorData() iff err.ErrorCode()==3 (i.e. an Execution error under
// the JSON RPC error codes IP).
func ExecutionErrData(err error) (interface{}, bool) {
	type rpcError interface {
		rpc.Error
		rpc.DataError
	}

	switch err := err.(type) {
	case rpcError:
		if err.ErrorCode() != 3 {
			return nil, false
		}
		return err.ErrorData(), true
	default:
		return nil, false
	}
}

// GasPrice is the assumed gas price, in GWei, when logging transaction costs.
var GasPrice uint64 = 200

// LogGas logs the amount and cost of the Transaction's gas. See GasPrice.
func LogGas(tb testing.TB, tx *types.Transaction, prefix string) {
	tb.Helper()

	cost := big.NewRat(int64(tx.Gas()*GasPrice), 1e9)
	tb.Logf("[%s] %s = %s%s @ %d gwei", prefix, humanize.Comma(int64(tx.Gas())), cost.FloatString(4), eth.Symbol, GasPrice)
}

// Comparers returns `extra`, appended with common comparison Options for
// cmp.Diff(); e.g. for big.Int,
func Comparers(extra ...cmp.Option) []cmp.Option {
	return append(
		extra,
		cmp.Comparer(func(a, b *big.Int) bool {
			switch {
			case a == nil && b == nil:
				return true
			case (a == nil) != (b == nil):
				return false
			default:
				return a.Cmp(b) == 0
			}
		}),
	)
}
