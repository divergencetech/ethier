package solidity_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/bazelbuild/tools_jvm_autodeps/thirdparty/golang/parsers/util/offset"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/solidity"
	"github.com/divergencetech/ethier/solidity/srcmaptest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/google/go-cmp/cmp"
)

func TestSourceMap(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, 1)

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

	// As we need to test the SourceMap against actual program counters in the
	// VM, we plug it into a vm.EVMLogger that is used to trace the
	// SimulatedBackend on all calls to the tX test-contracts above.
	src, err := srcmaptest.SourceMap()
	if err != nil {
		t.Fatalf("srcmaptest.SourceMap() error %v", err)
	}

	cfg := sim.Blockchain().GetVMConfig()
	cfg.Debug = true
	spy := &chainIDInterceptor{
		src: src,
	}
	cfg.Tracer = spy

	for _, fn := range [](func(*bind.TransactOpts) (*types.Transaction, error)){
		t0.Id, t0.IdPlusOne, t0.FromLib, t1.Id, t2.Id,
	} {
		sim.Must(t, "")(fn(sim.Acc(0)))
	}

	type pos struct {
		File string
		// Line and Col are 1-indexed as this is how IDEs display them.
		Line, Col int
		Length    int
	}
	wantLen := len("chainid()")
	want := []pos{
		{
			File:   "solidity/srcmaptest/SourceMapTest.sol",
			Line:   25,
			Col:    24,
			Length: wantLen,
		},
		{
			File:   "solidity/srcmaptest/SourceMapTest.sol",
			Line:   32,
			Col:    31,
			Length: wantLen,
		},
		{
			File:   "solidity/srcmaptest/SourceMapTest.sol",
			Line:   14,
			Col:    24,
			Length: wantLen,
		},
		{
			File:   "solidity/srcmaptest/SourceMapTest.sol",
			Line:   49,
			Col:    24,
			Length: wantLen,
		},
		{
			File:   "solidity/srcmaptest/SourceMapTest2.sol",
			Line:   15,
			Col:    24,
			Length: wantLen,
		},
	}

	var got []pos
	for _, g := range spy.got {
		// TODO: add this functionality to the *solidity.SourceMap itself.
		m := offset.NewMapper(string(srcmaptest.ReadSourceFile(t, g.Source)))
		line, col, err := m.LineAndColumn(g.Start)
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, pos{
			File:   g.Source,
			Line:   line + 1, // 1-indexed like an IDE
			Col:    col + 1,
			Length: g.Length,
		})
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

// chainIDInterceptor is a vm.EVMLogger that listens for vm.CHAINID operations,
// recording the solidity.Pos associated with each call.
type chainIDInterceptor struct {
	src *solidity.SourceMap

	// contracts are a stack of contract addresses with the last entry of the
	// slice being the current contract, against which the pc is compared when
	// inspecting the source map. Without a stack (i.e. always using the
	// "bottom" contract, to which the tx is initiated) the returned source will
	// function incorrectly on library calls.
	contracts []common.Address
	got       []solidity.Location
}

func (i *chainIDInterceptor) CaptureStart(evm *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	i.contracts = []common.Address{to}
}

func (i *chainIDInterceptor) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	if op != vm.CHAINID {
		return
	}

	c := i.contracts[len(i.contracts)-1]
	if pos, ok := i.src.Source(c, pc); ok {
		i.got = append(i.got, pos)
	} else {
		i.got = append(i.got, solidity.Location{
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
