// Package srcmaptest is a test-only package of generated Solidity bindings used
// to test the ethier/solidity package.
package srcmaptest

//go:generate ethier gen --experimental_src_map SourceMapTest.sol SourceMapTest2.sol CoverageTest.sol
