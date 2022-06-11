package solidity_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/solidity"
	"github.com/divergencetech/ethier/solidity/srcmaptest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSourceMap(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, 1)

	// As we need to test the SourceMap against actual program counters in the
	// VM, we plug it into a vm.EVMLogger that is used to trace the
	// SimulatedBackend on all calls to the contracts deployed below. Note that
	// source maps are all registered automatically by the generated code that
	// ethier adds on top of abigen, and that the EVMLogger only needs to call
	// RegisterDeployedContract() on all calls to CreateStart with create==true.
	cfg := sim.Blockchain().GetVMConfig()
	cfg.Debug = true
	spy := &chainIDInterceptor{}
	cfg.Tracer = spy

	// Deploying contracts ensures that we test whether their addresses are
	// passed to the SourceMap when it's constructed.
	_, _, t0, err := srcmaptest.DeploySourceMapTest0(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeploySourceMapTest0() error %v", err)
	}
	_, _, t1, err := srcmaptest.DeploySourceMapTest1(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeploySourceMapTest1() error %v", err)
	}
	_, _, t2, err := srcmaptest.DeploySourceMapTest2(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeploySourceMapTest2() error %v", err)
	}

	for _, fn := range [](func(*bind.TransactOpts) (*types.Transaction, error)){
		t0.Id, t0.IdPlusOne, t0.FromLib,
		t1.Id,
		t2.Id, t2.With1Modifier, t2.With2Modifiers,
	} {
		sim.Must(t, "")(fn(sim.Acc(0)))
	}

	wantLen := len("chainid()")
	want := []*solidity.Location{
		// All values were manually determined by inspection in an IDE.
		{
			Source:  "solidity/srcmaptest/SourceMapTest.sol",
			Start:   706,
			Length:  wantLen,
			Line:    25,
			EndLine: 25,
			Col:     24,
			EndCol:  24 + wantLen,
		},
		{
			Source:  "solidity/srcmaptest/SourceMapTest.sol",
			Start:   872,
			Length:  wantLen,
			Line:    32,
			EndLine: 32,
			Col:     31,
			EndCol:  31 + wantLen,
		},
		{
			Source:  "solidity/srcmaptest/SourceMapTest.sol",
			Start:   489,
			Length:  wantLen,
			Line:    14,
			EndLine: 14,
			Col:     24,
			EndCol:  24 + wantLen,
		},
		{
			Source:  "solidity/srcmaptest/SourceMapTest.sol",
			Start:   1208,
			Length:  wantLen,
			Line:    49,
			EndLine: 49,
			Col:     24,
			EndCol:  24 + wantLen,
		},
		{
			Source:  "solidity/srcmaptest/SourceMapTest2.sol",
			Start:   375,
			Length:  wantLen,
			Line:    15,
			EndLine: 15,
			Col:     24,
			EndCol:  24 + wantLen,
		},
		{
			Source:  "solidity/srcmaptest/SourceMapTest2.sol",
			Start:   490,
			Length:  wantLen,
			Line:    22,
			EndLine: 22,
			Col:     18,
			EndCol:  18 + wantLen,
		},
		{
			Source:        "solidity/srcmaptest/SourceMapTest2.sol",
			Start:         749,
			Length:        wantLen,
			Line:          36,
			EndLine:       36,
			Col:           24,
			EndCol:        24 + wantLen,
			ModifierDepth: 1,
		},
		{
			Source:  "solidity/srcmaptest/SourceMapTest2.sol",
			Start:   490,
			Length:  wantLen,
			Line:    22,
			EndLine: 22,
			Col:     18,
			EndCol:  18 + wantLen,
		},
		{
			Source:        "solidity/srcmaptest/SourceMapTest2.sol",
			Start:         595,
			Length:        wantLen,
			Line:          29,
			EndLine:       29,
			Col:           18,
			EndCol:        18 + wantLen,
			ModifierDepth: 1,
		},
		{
			Source:        "solidity/srcmaptest/SourceMapTest2.sol",
			Start:         958,
			Length:        wantLen,
			Line:          48,
			EndLine:       48,
			Col:           24,
			EndCol:        24 + wantLen,
			ModifierDepth: 2,
		},
	}

	ignore := []string{"FileIdx", "Jump"}
	opt := cmpopts.IgnoreFields(solidity.Location{}, ignore...)
	if diff := cmp.Diff(want, spy.got, opt); diff != "" {
		t.Error(diff)
	}
}

// chainIDInterceptor is a vm.EVMLogger that listens for vm.CHAINID operations,
// recording the solidity.Pos associated with each call.
type chainIDInterceptor struct {
	// contracts are a stack of contract addresses with the last entry of the
	// slice being the current contract, against which the pc is compared when
	// inspecting the source map. Without a stack (i.e. always using the
	// "bottom" contract, to which the tx is initiated) the returned source will
	// function incorrectly on library calls.
	contracts []common.Address
	got       []*solidity.Location
}

func (i *chainIDInterceptor) CaptureStart(evm *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	if create {
		solidity.RegisterDeployedContract(to, input)
	}
	i.contracts = []common.Address{to}
}

func (i *chainIDInterceptor) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	if op != vm.CHAINID {
		return
	}

	c := i.contracts[len(i.contracts)-1]
	if pos, ok := solidity.Source(c, pc); ok {
		i.got = append(i.got, pos)
	} else {
		i.got = append(i.got, &solidity.Location{
			Source: fmt.Sprintf("pc %d not found in contract %v", pc, c),
		})
	}
}

func (*chainIDInterceptor) CaptureTxStart(gasLimit uint64) {}

func (*chainIDInterceptor) CaptureTxEnd(restGas uint64) {}

func (*chainIDInterceptor) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {}

func (i *chainIDInterceptor) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	i.contracts = append(i.contracts, to)
}

func (i *chainIDInterceptor) CaptureExit(output []byte, gasUsed uint64, err error) {
	i.contracts = i.contracts[:len(i.contracts)-1]
}

func (*chainIDInterceptor) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}
