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
	"regexp"
	"strconv"
	"strings"

	"github.com/bazelbuild/tools_jvm_autodeps/thirdparty/golang/parsers/util/offset"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/vm"
)

// A Source is a source file used by solc when compiling contracts. The File is
// the string output from solc's source list in the combined JSON, and Code is
// the contents of File.
type Source struct {
	File, Code string
}

// NewSourceMap returns a SourceMap based on solc --combined-json output,
// coupled with deployment addresses of the contracts. The sources are the
// sourceList array from the combined JSON, and the order must not be changed
// The key of the contracts map must match the value of the deployed map to
// allow the SourceMap to find the correct mapping when only an address is
// available.
func NewSourceMap(sources []*Source, contracts map[string]*compiler.Contract, deployed map[common.Address]string) (*SourceMap, error) {
	// See https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html

	sm := &SourceMap{
		sources:   sources,
		contracts: make(map[string]*contractMap),
		deployed:  deployed,
	}

	for _, src := range sources {
		m := offset.NewMapper(src.Code)
		sm.mappers = append(sm.mappers, m)
	}

	isDeployed := make(map[string]bool)
	for _, name := range deployed {
		isDeployed[name] = true
	}

	for name, c := range contracts {
		if !isDeployed[name] {
			continue
		}

		cm := new(contractMap)
		if err := cm.parseCode(c); err != nil {
			return nil, fmt.Errorf("%q: %v", name, err)
		}
		if err := cm.parseSrcMap(c, sm); err != nil {
			return nil, fmt.Errorf("%q: %v", name, err)
		}
		sm.contracts[name] = cm
	}

	return sm, nil
}

// A SourceMap maps a program counter from an EVM trace to the original position
// in the Solidity source code.
type SourceMap struct {
	sources   []*Source
	contracts map[string]*contractMap
	deployed  map[common.Address]string
	mappers   []*offset.Mapper
}

// Source returns the code location that was compiled into the instruction at
// the specific program counter in the deployed contract.
func (sm *SourceMap) Source(contract common.Address, pc uint64) (*Location, bool) {
	name, ok := sm.deployed[contract]
	if !ok {
		return nil, false
	}
	return sm.SourceByName(name, pc)
}

// SourceByName functions identically to Source but doesn't require that the
// contract has been deployed. NOTE that there isn't a one-to-one mapping
// between runtime byte code (i.e. program counter) and instruction number
// because the PUSH* instructions require additional bytes; see
// https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html.
func (sm *SourceMap) SourceByName(contract string, pc uint64) (*Location, bool) {
	cm, ok := sm.contracts[contract]
	if !ok {
		return nil, false
	}
	p, ok := cm.codeLocation(pc)
	if !ok {
		return nil, false
	}

	if p.FileIdx < len(sm.sources) {
		p.Source = sm.sources[p.FileIdx].File
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

	// Line and Col are both 1-indexed as this is typical behaviour of IDEs and
	// coverage reports.
	Line, Col int
	// EndLine and EndCol are Length bytes after Line and Col, also 1-indexed.
	EndLine, EndCol int

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
	locations       []*Location
	pcToInstruction map[uint64]int
}

func (cm *contractMap) codeLocation(pc uint64) (*Location, bool) {
	i, ok := cm.pcToInstruction[pc]
	if !ok {
		return nil, false
	}
	return cm.locations[i], true
}

var (
	// libraryPlaceHolder finds all places in which bind.Bind has inserted a
	// string identifying a library address to be pushed (PUSH20 == 0x73). In
	// actual deployment this value is replaced, by the generated code, with the
	// deployed library's address, but for source mapping it can be ignored
	// because the PUSH20 means that contractMap.parseCode() will skip the 20
	// bytes. We can therefore replace it with a push of anything, so use the
	// zero address for simplicity.
	libraryPlaceholder = regexp.MustCompile(`73__\$[[:xdigit:]]+\$__`)
	pushZeroAddress    = fmt.Sprintf("73%x", common.Address{})
)

// parseCode converts a Contract's runtime byte code into a mapping from
// program counter (position in byte code) to instruction number because the
// PUSH* OpCodes require additional byte code but the source-map is based on
// instruction number. It saves the mapping to the contractMap.
func (cm *contractMap) parseCode(c *compiler.Contract) error {
	rawCode := strings.TrimPrefix(c.RuntimeCode, "0x")
	rawCode = libraryPlaceholder.ReplaceAllString(rawCode, pushZeroAddress)

	code, err := hex.DecodeString(rawCode)
	if err != nil {
		return fmt.Errorf("hex.DecodeString(%T.RuntimeCode): %v", c, err)
	}

	var instruction int
	pcToInstruction := make(map[uint64]int)

	for i, n := 0, len(code); i < n; i++ {
		pcToInstruction[uint64(i)] = instruction
		instruction++

		c := vm.OpCode(code[i])
		if c.IsPush() {
			i += int(c - vm.PUSH0)
		}
	}
	cm.pcToInstruction = pcToInstruction

	return nil
}

// nMappingFields is the number of fields in the solc source mapping: s,l,f,j,m.
const nMappingFields = 5

// srcMapNode captures the s:l:f:j:m pattern of solc source mapping. If
// compressed, empty values mean that the previous value from the last node must
// be copied across.
type srcMapNode [nMappingFields]string

// parseSrcMap decompresses the Contract's runtime source map, storing each of
// the locations in the contractMap's `locations` slice. The indices in the
// slice correspond to the values in cm.pcToInstruction.
func (cm *contractMap) parseSrcMap(c *compiler.Contract, sm *SourceMap) error {
	instructions := strings.Split(c.Info.SrcMapRuntime, ";")
	cm.locations = make([]*Location, len(instructions))

	var last, curr srcMapNode
	for i, instr := range instructions {
		copy(curr[:], strings.Split(instr, ":"))
		for j, n := 0, len(curr); j < nMappingFields; j++ {
			if j < n && curr[j] != "" {
				last[j] = curr[j]
			} else {
				curr[j] = last[j]
			}
		}

		loc, err := locationFromNode(curr, sm)
		if err != nil {
			return err
		}
		cm.locations[i] = loc
	}

	return nil
}

func locationFromNode(node srcMapNode, sm *SourceMap) (*Location, error) {
	start, err := strconv.Atoi(node[0])
	if err != nil {
		return nil, fmt.Errorf("parse `s`: %v", err)
	}
	length, err := strconv.Atoi(node[1])
	if err != nil {
		return nil, fmt.Errorf("parse `l`: %v", err)
	}
	fileIdx, err := strconv.Atoi(node[2])
	if err != nil {
		return nil, fmt.Errorf("parse `f`: %v", err)
	}
	modDepth, err := strconv.Atoi(node[4])
	if err != nil {
		return nil, fmt.Errorf("parse `m`: %v", err)
	}

	loc := &Location{
		Start:         start,
		Length:        length,
		FileIdx:       fileIdx,
		Jump:          JumpType(node[3]),
		ModifierDepth: modDepth,
	}

	if fileIdx >= len(sm.mappers) {
		return loc, nil
	}

	m := sm.mappers[fileIdx]
	if m.Len() == 0 {
		return loc, nil
	}

	for offset, set := range map[int][2]*int{
		start:          {&loc.Line, &loc.Col},
		start + length: {&loc.EndLine, &loc.EndCol},
	} {
		l, c, err := m.LineAndColumn(offset)
		if err != nil {
			return nil, fmt.Errorf("%T.LineAndColumn(%d) of %q: %v", m, offset, sm.sources[fileIdx], err)
		}
		*set[0] = l + 1
		*set[1] = c + 1
	}
	return loc, nil
}
