// Package solidity provides mapping from EVM-trace program counters to original
// Solidity source code, including automated coverage collection for tests.
//
// This package doesn't typically need to be used directly, and is automatically
// supported by adding the source-map flag to `ethier gen` of the
// github.com/divergencetech/ethier/ethier binary for generating Go bindings.
//
// See https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html
// for more information.
package solidity

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/bazelbuild/tools_jvm_autodeps/thirdparty/golang/parsers/util/offset"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/vm"
)

// A Location is an offset-based location in a Solidity file. Using notation
// described in https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html,
// s = Start, l = Length, f = FileIdx, j = Jump, and m = ModifierDepth.
type Location struct {
	Start, Length int

	// FileIdx refers to the index of the source file in the inputs to solc, as
	// returned in solc output, but can generally be ignored in favour of the
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

// A compiledContract couples a *compiler.Contract with the fully qualified name
// of the Solidity contract it represents, and the solc source-list of files
// that were used to compile it. A fully qualified name is the concatenation of
// the source file name, a colon :, and the name of the contract.
type compiledContract struct {
	*compiler.Contract
	name       string
	sourceList []string

	instructions pcToInstruction
	locations    []*Location
}

// location converts the program counter into an instruction number, and returns
// the corresponding Location and true. If the pc was not in the runtime source
// map of the contract, or if cc == nil, then location returns (nil, false).
func (cc *compiledContract) location(pc uint64) (*Location, bool) {
	if cc == nil {
		// This allows direct calls on values returned from maps when the key
		// didn't exist, avoiding an extra check of the ok.
		return nil, false
	}

	i, ok := cc.instructions[pc]
	if !ok {
		return nil, false
	}
	return cc.locations[i], true
}

// A contractMatcher allows matching of deployed contracts against solc output
// that includes libraries, allowing the library address to change.
type contractMatcher struct {
	re *regexp.Regexp
	*compiledContract
}

// A sourceFile holds the input to solc for a particular file in a compilation's
// source list.
type sourceFile struct {
	contents   string
	mapper     *offset.Mapper
	isExternal bool
}

var (
	sourceCode = make(map[string]*sourceFile)
	// contractsByName is keyed by the fully qualified name of the contract, its
	// source file and Solditiy name; e.g. path/to/file.sol:MyContract.
	contractsByName = make(map[string]*compiledContract)
	// contractsByHash identifies contracts by the SHA256 hash of their byte
	// code (creation, not runtime). The sha256 package is convenient because it
	// returns a fixed-size array instead of a slice (like crypto.Keccak256),
	// which can be used as a map key.
	contractsByHash = make(map[[sha256.Size]byte]*compiledContract)
	// contractMatchers are stored by the hash of their regex pattern to avoid
	// duplication.
	contractMatchers = make(map[[sha256.Size]byte]contractMatcher)
	// deployedContracts maps contracts by their deployed addresses iff the
	// contract has already been registered.
	deployedContracts = make(map[common.Address]*compiledContract)
)

// RegisterSourceCode registers the contents of source files passed in source
// lists to RegisterContract. This allows op codes in contracts, deployed or
// otherwise, to be mapped back to the specific Solidity code that resulted in
// their compilation.
//
// RegisterSourceCode SHOULD be called before all calls to RegisterContract that
// include fileName in the sourceList otherwise Location values will not contain
// line and column numbers. External source code, e.g. OpenZeppelin contracts,
// won't be monitored in coverage collection, but will be available for
// Etherscan verification.
func RegisterSourceCode(fileName, contents string, isExternal bool) {
	sourceCode[fileName] = &sourceFile{
		contents:   contents,
		mapper:     offset.NewMapper(contents),
		isExternal: isExternal,
	}
}

// RegisterContract registers a compiled contract by its fully qualified name.
// This allows inspection of EVM traces to watch for contract deployments that
// are matched against already-registered contracts, and then for mapping each
// step in the trace back to original source code via the program counter. See
// SourceByName() documentation re fully qualified names.
//
// The order of the sourceList MUST match the solc output from which the
// *Contract was parsed. The Contract's Info.SrcMapRuntime represents files as
// indices into this slice so a change in order will result in invalid results.
// All files included in sourceList SHOULD be registered via RegisterSourceCode
// before calling RegisterContract.
//
// RegisterContract MUST be called before deployment otherwise
// RegisterDeployedContract will fail to match the byte code. This is typically
// done as part on an init() function, and `ethier gen` generated code performs
// this step automatically.
func RegisterContract(name string, c *compiler.Contract, sourceList []string) {
	instructions, err := parseCode(c)
	if err != nil {
		panic(fmt.Sprintf("Parsing RuntimeCode of %q: %v", name, err))
	}

	locations, err := parseSrcMap(c, sourceList)
	if err != nil {
		panic(fmt.Sprintf("Parsing SrcMap of %q: %v", name, err))
	}

	cc := &compiledContract{
		Contract:     c,
		name:         name,
		sourceList:   sourceList,
		instructions: instructions,
		locations:    locations,
	}
	contractsByName[name] = cc

	if libraryPlaceholder.MatchString(c.Code) {
		registerContractByRegexp(cc)
	} else {
		registerContractByHash(cc)
	}
}

func registerContractByRegexp(cc *compiledContract) {
	pattern := libraryPlaceholder.ReplaceAllString(
		strings.TrimPrefix(cc.Code, "0x"),
		"73[[:xdigit:]]{40}",
	)
	contractMatchers[sha256.Sum256([]byte(pattern))] = contractMatcher{
		re:               regexp.MustCompile(pattern),
		compiledContract: cc,
	}
}

func registerContractByHash(cc *compiledContract) {
	code := strings.TrimPrefix(cc.Code, "0x")
	bin, err := hex.DecodeString(code)
	if err != nil {
		// panic is used instead of returning an error because the expected
		// usage of RegisterContract() is in init() functions of generated code,
		// so it's not possible to propagate an error. log.Fatal() wouldn't
		// provide enough context but panic gives the code location.
		panic(fmt.Sprintf("solidity.RegisterContract(): hex.DecodeString(%q): %v", code, err))
	}
	contractsByHash[sha256.Sum256(bin)] = cc
}

// RegisterDeployedContract matches the code against contracts registered with
// RegisterContract, allowing future calls to Source(addr, …) to return
// data pertaining to the correct contract.
//
// RegisterDeployedContract should be called by EVMLogger.CaptureStart when the
// create flag is true, passing the deployment address and the input code bytes,
// th
func RegisterDeployedContract(addr common.Address, code []byte) {
	c, ok := contractsByHash[sha256.Sum256(code)]
	if ok {
		deployedContracts[addr] = c
		return
	}

	for _, m := range contractMatchers {
		if m.re.MatchString(hex.EncodeToString(code)) {
			deployedContracts[addr] = m.compiledContract
			return
		}
	}
}

// Source returns the code location that was compiled into the instruction at
// the specific program counter in the deployed contract. The contract's address
// MUST have been registered with RegisterDeployedContract().
func Source(contract common.Address, pc uint64) (*Location, bool) {
	return deployedContracts[contract].location(pc)
}

// SourceByName functions identically to Source but doesn't require that the
// contract has been deployed. The contract MUST have been registered with
// RegisterContract().
//
// NOTE that there isn't a one-to-one mapping between runtime byte code (i.e.
// program counter) and instruction number because the PUSH* instructions
// require additional bytes as documented in:
// https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html.
func SourceByName(contract string, pc uint64) (*Location, bool) {
	return contractsByName[contract].location(pc)
}

var (
	// libraryPlaceHolder finds all places in which bind.Bind has inserted a
	// string identifying a library address to be pushed (PUSH20 == 0x73). In
	// actual deployment this value is replaced, by the generated code, with the
	// deployed library's address, but for source mapping it can be ignored
	// because the PUSH20 means that contractMap.parseCode() will skip the 20
	// bytes. We can therefore replace it with a push of anything, so use the
	// zero address for simplicity.
	libraryPlaceholder = regexp.MustCompile(`73__\$[[:xdigit:]]{34}\$__`)
	pushZeroAddress    = fmt.Sprintf("73%x", common.Address{})
)

// pcToInstruction maps a program counter to an instruction number. As PUSH<N>
// op codes use N+1 bytes, the program counter within a specific contract's
// runtime code does not match to the instruction number. A pcToInstruction map
// is therefore contract-specific.
type pcToInstruction map[uint64]int

// parseCode converts a Contract's runtime byte code into a mapping from
// program counter (position in byte code) to instruction number because the
// PUSH* OpCodes require additional byte code but the source-map is based on
// instruction number. It saves the mapping to the contractMap.
func parseCode(c *compiler.Contract) (pcToInstruction, error) {
	rawCode := strings.TrimPrefix(c.RuntimeCode, "0x")
	rawCode = libraryPlaceholder.ReplaceAllString(rawCode, pushZeroAddress)

	code, err := hex.DecodeString(rawCode)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString(%T.RuntimeCode): %v", c, err)
	}

	var instruction int
	instructions := make(pcToInstruction)

	for i, n := 0, len(code); i < n; i++ {
		instructions[uint64(i)] = instruction
		instruction++

		c := vm.OpCode(code[i])
		if c.IsPush() {
			i += int(c - vm.PUSH0)
		}
	}
	return instructions, nil
}

// nMappingFields is the number of fields in the solc source mapping: s,l,f,j,m.
const nMappingFields = 5

// srcMapNode captures the s:l:f:j:m pattern of solc source mapping. If
// compressed, empty values mean that the previous value from the last node must
// be copied across.
type srcMapNode [nMappingFields]string

// parseSrcMap decompresses the Contract's runtime source map, returning a slice
// of Locations, indexed by instruction number; i.e. the indices in the slice
// correspond to the values a pcToInstruction map.
func parseSrcMap(c *compiler.Contract, sourceList []string) ([]*Location, error) {
	instructions := strings.Split(c.Info.SrcMapRuntime, ";")
	locations := make([]*Location, len(instructions))

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

		loc, err := locationFromNode(curr, sourceList)
		if err != nil {
			return nil, err
		}
		locations[i] = loc
	}

	return locations, nil
}

// locationFromNode parses an s:l:f:j:m node from a solc source map, returning
// the corresponding Location.
func locationFromNode(node srcMapNode, sourceList []string) (*Location, error) {
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

	if start < 0 || length < 0 || fileIdx < 0 {
		// TODO(aschlosberg) investigate the meaning of -1 values in the
		// compressed output. For now, just make it an impossible file index.
		fileIdx = math.MaxInt64
		loc.FileIdx = fileIdx
	}

	if fileIdx >= len(sourceList) {
		return loc, nil
	}

	loc.Source = sourceList[fileIdx]
	c := sourceCode[loc.Source]
	m := c.mapper
	if m.Len() == 0 {
		return loc, nil
	}

	for offset, set := range map[int][2]*int{
		start:          {&loc.Line, &loc.Col},
		start + length: {&loc.EndLine, &loc.EndCol},
	} {
		l, c, err := m.LineAndColumn(offset)
		if err != nil {
			return nil, fmt.Errorf("%T.LineAndColumn(%d) of %q: %v", m, offset, sourceList[fileIdx], err)
		}
		*set[0] = l + 1
		*set[1] = c + 1
	}
	return loc, nil
}
