package solcover

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/compiler"
)

func TestMapDecompression(t *testing.T) {
	const (
		// From the example at
		// https://docs.soliditylang.org/en/v0.8.14/internals/source_mappings.html
		// but including modifier depth and jump type.
		compressed   = `1:2:1:-:0;:9;2:1:2;;`
		uncompressed = `1:2:1:-:0;1:9:1:-:0;2:1:2:-:0;2:1:2:-:0;2:1:2:-:0`
	)

	in := &compiler.Contract{
		Info: compiler.ContractInfo{
			SrcMapRuntime: compressed,
		},
	}

	locations, err := parseSrcMap(in, nil)
	if err != nil {
		t.Fatalf("parseSrcMap(%+v, nil) error %v", in, err)
	}

	var gotParts []string
	for _, l := range locations {
		gotParts = append(gotParts, fmt.Sprintf("%d:%d:%d:%s:%d", l.Start, l.Length, l.FileIdx, l.Jump, l.ModifierDepth))
	}
	if got, want := strings.Join(gotParts, ";"), uncompressed; got != want {
		t.Errorf("parseSrcMap(%+v) got reconstructed nodes %q; want %q", in, got, want)
	}
}
