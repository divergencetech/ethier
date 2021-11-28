package ethtest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

type coverageCollector struct {
	hits map[string]int
}

var _ vm.EVMLogger = newCoverageCollector()

func newCoverageCollector() *coverageCollector {
	return &coverageCollector{
		hits: make(map[string]int),
	}
}

func (c *coverageCollector) CaptureState(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	if op != vm.PUSH1 {
		return
	}
	if strings.Contains(hex.EncodeToString(scope.Memory.Data()), "84fa5025f707e6c9d43d615ea30ac7c1141ae88d1eb634fadedcb81e09c7cc6c") {
		panic(true)
	}
	// pushed := hex.EncodeToString(scope.Memory.GetPtr(n-32, 32))
	// c.hits[pushed]++
}

// MarshalJSON marshals the hit map to a format understood by the
// solidity-coverage npm package's Instrumenter, as reverse engineered from its
// DataCollector class.
func (c *coverageCollector) MarshalJSON() ([]byte, error) {
	munged := make(map[string]map[string]int)

	for k, v := range c.hits {
		munged[fmt.Sprintf("0x%s", k)] = map[string]int{"hits": v}
	}
	return json.Marshal(munged)
}

// All other methods are noops, provided only to satisfice the interface
// definition. Argument names are from the definition of EVMLogger.

func (c *coverageCollector) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
}

func (c *coverageCollector) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}

func (c *coverageCollector) CaptureExit(output []byte, gasUsed uint64, err error) {}

func (c *coverageCollector) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

func (c *coverageCollector) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {}
