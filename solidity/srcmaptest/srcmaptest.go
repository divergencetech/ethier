// Package srcmaptest is a test-only package of generated Solidity bindings used
// to test the ethier/solidity package.
package srcmaptest

import (
	"embed"
	"strings"
	"testing"
)

//go:generate ethier gen --experimental_src_map SourceMapTest.sol SourceMapTest2.sol CoverageTest.sol

//go:embed *.sol
var sourceFiles embed.FS

// ReadSourceFile returns the Solidity source of the specified file.
func ReadSourceFile(t *testing.T, file string) []byte {
	t.Helper()

	const prefix = "solidity/srcmaptest/"

	if !strings.HasPrefix(file, prefix) {
		t.Fatalf("srcmaptest.ReadSourceFile(%q) file must have prefix %q", file, prefix)
	}

	buf, err := sourceFiles.ReadFile(strings.TrimPrefix(file, prefix))
	if err != nil {
		t.Fatalf("srcmaptest.%T.ReadFile(%q with prefix trimmed) error %v", sourceFiles, file, err)
	}
	return buf
}
