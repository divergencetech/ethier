package solidity_test

import (
	"os"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/solidity/srcmaptest"
)

func TestCoverageCollector(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, 1)

	_, _, cov, err := srcmaptest.DeployCoverageTest(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployCoverageTest() error %v", err)
	}

	sim.CollectCoverageTB(t, srcmaptest.SourceMap)

	sim.Must(t, "%T.Foo()", cov)(cov.Foo(sim.Acc(0)))
	sim.Must(t, "%T.Bar()", cov)(cov.Bar(sim.Acc(0)))

	// This has been confirmed visually by running the report through the LCOV
	// genhtml script.
	const want = `SF:solidity/srcmaptest/CoverageTest.sol
	FNF:0
	FNH:0
	DA:11,2
	DA:15,2
	DA:16,1
	DA:17,1
	DA:18,2
	DA:19,1
	DA:21,0
	DA:22,0
	DA:24,1
	DA:25,1
	DA:28,2
	DA:29,2
	DA:30,0
	DA:31,2
	DA:32,0
	DA:33,1
	DA:34,1
	DA:35,1
	DA:36,2
	DA:37,1
	DA:39,0
	DA:43,1
	DA:44,0
	DA:47,1
	DA:48,1
	DA:49,1
	DA:50,0
	DA:53,2
	DA:54,1
	DA:56,0
	LH:30
	LF:61
	end_of_record
	SF:solidity/srcmaptest/SourceMapTest.sol
	FNF:0
	FNH:0
	LH:0
	LF:54
	end_of_record
	SF:solidity/srcmaptest/SourceMapTest2.sol
	FNF:0
	FNH:0
	LH:0
	LF:20
	end_of_record
	`
)
