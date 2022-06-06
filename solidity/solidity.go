// Package solidity provides a source map to determine original code locations
// from op codes in EVM traces. This package doesn't typically need to be used
// directly, and can be automatically supported by adding the source-map flag to
// `ethier gen`.
//
// See https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html
// for more information.
package solidity

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/vm"
)

// NewSourceMap returns a SourceMap based on solc --combined-json output,
// coupled with deployment addresses of the contracts. The sources are the
// sourceList array from the combined JSON, and the order must not be changed.
// The key of the contracts map must match the value of the deployed map to
// allow the SourceMap to find the correct mapping when only an address is
// available.
func NewSourceMap(sources []string, contracts map[string]*compiler.Contract, deployed map[common.Address]string) (*SourceMap, error) {
	// See https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html

	sm := &SourceMap{
		sources:   sources,
		contracts: make(map[string]*contractMap),
		deployed:  deployed,
	}

	for name, c := range contracts {
		cm := new(contractMap)
		if err := cm.parseCode(c); err != nil {
			return nil, fmt.Errorf("%q: %v", name, err)
		}
		if err := cm.parseSrcMap(c); err != nil {
			return nil, fmt.Errorf("%q: %v", name, err)
		}
		sm.contracts[name] = cm
	}

	return sm, nil
}

// A SourceMap maps a program counter from an EVM trace to the original position
// in the Solidity source code.
type SourceMap struct {
	sources   []string
	contracts map[string]*contractMap
	deployed  map[common.Address]string
}

// Source returns the code location that was compiled into the instruction at
// the specific program counter in the deployed contract.
func (sm *SourceMap) Source(contract common.Address, pc uint64) (Location, bool) {
	name, ok := sm.deployed[contract]
	if !ok {
		return Location{}, false
	}
	return sm.SourceByName(name, pc)
}

// SourceByName functions identically to Source but doesn't require that the
// contract has been deployed. NOTE that there isn't a one-to-one mapping
// between runtime byte code (i.e. program counter) and instruction number
// because the PUSH* instructions require additional bytes; see
// https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html.
func (sm *SourceMap) SourceByName(contract string, pc uint64) (Location, bool) {
	cm, ok := sm.contracts[contract]
	if !ok {
		return Location{}, false
	}
	p, ok := cm.codeLocation(pc)
	if !ok {
		return Location{}, false
	}

	if p.FileIdx < len(sm.sources) {
		p.Source = sm.sources[p.FileIdx]
	}
	return p, true
}

// A Location is an offset-based location in a Solidity file. Using notation
// described in https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html,
// s = Start, l = Length, f = FileIdx, j = Jump, and m = ModifierDepth.
type Location struct {
	Start, Length int

	// FileIdx refers to the index of the source file in the inputs to solc, as
	// passed to NewSourceMap(), but can generally be ignored in favour of the
	// Source, which is determined from the NewSourceMap() input.
	FileIdx int
	Source  string

	Jump          JumpType
	ModifierDepth int
}

// A JumpType describes the action of a JUMP instruction.
type JumpType string

// Possible JumpType values.
const (
	FunctionIn  = JumpType(`i`)
	FunctionOut = JumpType(`o`)
	RegularJump = JumpType(`-`)
)

// A contractMap is the contract-specific implementation of SourceMap as a
// SourceMap can be used on an arbitrary set of contracts.
type contractMap struct {
	locations       []Location
	pcToInstruction map[uint64]int
}

func (cm *contractMap) codeLocation(pc uint64) (Location, bool) {
	i, ok := cm.pcToInstruction[pc]
	if !ok {
		return Location{}, false
	}
	return cm.locations[i], true
}

// parseCode converts a Contract's runtime byte code into a mapping from
// program counter (position in byte code) to instruction number because the
// PUSH* OpCodes require additional byte code but the source-map is based on
// instruction number. It saves the mapping to the contractMap.
func (cm *contractMap) parseCode(c *compiler.Contract) error {
	code, err := hex.DecodeString(strings.TrimPrefix(c.RuntimeCode, "0x"))
	if err != nil {
		return fmt.Errorf("hex.DecodeString(%T.RuntimeCode): %v", c, err)
	}

	var instruction int
	pcToInstruction := make(map[uint64]int)

	for i, n := 0, len(code); i < n; i++ {
		pcToInstruction[uint64(i)] = instruction
		instruction++

		c := vm.OpCode(code[i])
		if c != vm.PUSH0 && c.IsPush() {
			i += int(c - vm.PUSH1 + 1)
		}
	}
	cm.pcToInstruction = pcToInstruction

	return nil
}

// parseSrcMap decompresses the Contract's runtime source map, storing each of
// the locations in the contractMap's `locations` slice. The indices in the
// slice correspond to the values in cm.pcToInstruction.
func (cm *contractMap) parseSrcMap(c *compiler.Contract) error {
	instructions := strings.Split(c.Info.SrcMapRuntime, ";")
	cm.locations = make([]Location, 0, len(instructions))

	const nFields = 5
	var last, curr [nFields]string
	for _, ins := range instructions {
		copy(curr[:], strings.Split(ins, ":"))
		for i, n := 0, len(curr); i < nFields; i++ {
			if i < n && curr[i] != "" {
				last[i] = curr[i]
			} else {
				curr[i] = last[i]
			}
		}

		start, err := strconv.Atoi(curr[0])
		if err != nil {
			return fmt.Errorf("parse `s`: %v", err)
		}
		length, err := strconv.Atoi(curr[1])
		if err != nil {
			return fmt.Errorf("parse `l`: %v", err)
		}
		fileIdx, err := strconv.Atoi(curr[2])
		if err != nil {
			return fmt.Errorf("parse `f`: %v", err)
		}
		modDepth, err := strconv.Atoi(curr[4])
		if err != nil {
			return fmt.Errorf("parse `m`: %v", err)
		}

		cm.locations = append(cm.locations, Location{
			Start:         start,
			Length:        length,
			FileIdx:       fileIdx,
			Jump:          JumpType(curr[3]),
			ModifierDepth: modDepth,
		})
	}

	return nil
}
