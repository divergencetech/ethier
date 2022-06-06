package solidity

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/compiler"
)

func (cm *contractMap) uncompressed() string {
	var parts []string
	for _, p := range cm.locations {
		parts = append(parts, fmt.Sprintf("%d:%d:%d:%s:%d", p.Start, p.Length, p.FileIdx, p.Jump, p.ModifierDepth))
	}
	return strings.Join(parts, ";")
}

func TestMapDecompression(t *testing.T) {
	const (
		// From the example at
		// https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html
		// but including modifier depth and jump type.
		compressed   = `1:2:1:-:0;:9;2:1:2;;`
		uncompressed = `1:2:1:-:0;1:9:1:-:0;2:1:2:-:0;2:1:2:-:0;2:1:2:-:0`
	)

	in := map[string]*compiler.Contract{
		"dummy": {
			RuntimeCode: "0x00",
			Info: compiler.ContractInfo{
				SrcMapRuntime: compressed,
			},
		},
	}
	sm, err := NewSourceMap(nil, in, nil)
	if err != nil {
		t.Fatalf("NewSourceMap(%+v): %v", in, err)
	}

	if got, want := sm.contracts["dummy"].uncompressed(), uncompressed; got != want {
		t.Errorf("NewSourceMap(%+v) got uncompressed mapping %q; want %q", in, got, want)
	}
}
