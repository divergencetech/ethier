package solidity

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// CoverageCollector returns an EVMLogger that can be used to trace EVM
// operations for code coverage. The returned function returns an LCOV trace
// file at any time coverage is not being actively collected (i.e. it is not
// thread safe with respect to VM computation).
func CoverageCollector() (vm.EVMLogger, func() []byte) {
	lineHits := make(map[string]map[int]int)
	for file := range sourceMappers {
		lineHits[file] = make(map[int]int)
	}

	// TODO(aschlosberg) this is only a proxy for all "possible" lines in the
	// code (i.e. not blank lines and comments). It's imperfect because the
	// lines are derived from the op codes so "else lines" are missed. This will
	// need to be dealt with when tracing branches.
	var contracts []*compiledContract
	for _, c := range contractsByHash {
		contracts = append(contracts, c)
	}
	for _, m := range contractMatchers {
		contracts = append(contracts, m.compiledContract)
	}
	for _, c := range contracts {
		for _, l := range c.locations {
			if l.FileIdx >= len(c.sourceList) {
				continue
			}
			f := c.sourceList[l.FileIdx]
			lineHits[f][l.Line] = 0
		}
	}

	cc := &coverageCollector{
		lineHits: lineHits,
		last:     new(Location),
	}
	return cc, cc.lcovTraceFile
}

type coverageCollector struct {
	// contracts are a stack of contract addresses with the last entry of the
	// slice being the current contract, against which the pc is compared when
	// inspecting the source map. Without a stack (i.e. always using the
	// "bottom" contract, to which the tx is initiated) the returned source will
	// function incorrectly on library calls.
	contracts []common.Address

	lineHits map[string]map[int]int
	last     *Location
}

func (cc *coverageCollector) lcovTraceFile() []byte {
	var out bytes.Buffer
	linef := func(format string, a ...interface{}) {
		out.WriteString(fmt.Sprintf(format, a...))
		out.WriteRune('\n')
	}
	line := func(a ...interface{}) {
		out.WriteString(fmt.Sprint(a...))
		out.WriteRune('\n')
	}

	// It's important to range over a slice here and not the cc.lineHits map
	// because the map lacks a guaranteed order.
	var files []string
	for f := range sourceMappers {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, file := range files {
		mapper, ok := sourceMappers[file]
		if !ok || mapper.Len() == 0 {
			continue
		}
		linef("SF:%s", file)

		// TODO(aschlosberg) extend coverage to include functions and branches
		line("FNF:0")
		line("FNH:0")

		type count struct{ line, n int }
		lh := cc.lineHits[file]
		counts := make([]count, 0, len(lh))
		for l, n := range lh {
			counts = append(counts, count{l, n})
		}
		sort.Slice(counts, func(i, j int) bool {
			return counts[i].line < counts[j].line
		})

		var nonZero int
		for _, c := range counts {
			linef("DA:%d,%d", c.line, c.n)
			if c.n > 0 {
				nonZero++
			}
		}
		linef("LH:%d", nonZero)
		linef("LF:%d", len(counts))

		line("end_of_record")
	}

	return out.Bytes()
}

func (cc *coverageCollector) CaptureStart(evm *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	if create {
		RegisterDeployedContract(to, input)
	}
	cc.contracts = []common.Address{to}
}

func (cc *coverageCollector) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	// TODO(aschlosberg) this is adding spurious extra counts to lines.
	loc, ok := Source(cc.contracts[len(cc.contracts)-1], pc)
	if !ok || loc.Source == "" || loc.Line == 0 {
		return
	}
	if loc.FileIdx != cc.last.FileIdx || loc.Line != cc.last.Line {
		cc.lineHits[loc.Source][loc.Line]++
	}
	cc.last = loc
}

func (*coverageCollector) CaptureTxStart(gasLimit uint64) {}

func (*coverageCollector) CaptureTxEnd(restGas uint64) {}

func (*coverageCollector) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {}

func (cc *coverageCollector) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	cc.contracts = append(cc.contracts, to)
}

func (cc *coverageCollector) CaptureExit(output []byte, gasUsed uint64, err error) {
	cc.contracts = cc.contracts[:len(cc.contracts)-1]
}

func (*coverageCollector) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}
