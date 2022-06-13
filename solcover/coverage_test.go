package solcover_test

import (
	"strings"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/solcover/srcmaptest"
	"github.com/google/go-cmp/cmp"
)

func TestCoverageCollector(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, 1)

	_, _, cov, err := srcmaptest.DeployCoverageTest(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployCoverageTest() error %v", err)
	}

	// sim.CollectCoverageTB(t, srcmaptest.SourceMap)
	sim.Must(t, "%T.Foo()", cov)(cov.Foo(sim.Acc(0)))
	sim.Must(t, "%T.Bar()", cov)(cov.Bar(sim.Acc(0)))

	// This has been confirmed visually by running the report through the LCOV
	// genhtml script.
	const want = `SF:solcover/srcmaptest/CoverageTest.sol
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
DA:29,3
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
LH:22
LF:30
end_of_record
SF:solcover/srcmaptest/SourceMapTest.sol
FNF:0
FNH:0
DA:11,0
DA:12,0
DA:14,0
DA:19,0
DA:23,0
DA:25,0
DA:27,0
DA:30,0
DA:32,0
DA:34,0
DA:35,0
DA:38,0
DA:39,0
DA:40,0
DA:44,0
DA:47,0
DA:49,0
DA:51,0
LH:0
LF:18
end_of_record
SF:solcover/srcmaptest/SourceMapTest2.sol
FNF:0
FNH:0
DA:9,0
DA:13,0
DA:15,0
DA:17,0
DA:22,0
DA:24,0
DA:29,0
DA:34,0
DA:36,0
DA:38,0
DA:41,0
DA:43,0
DA:44,0
DA:45,0
DA:48,0
DA:50,0
LH:0
LF:16
end_of_record
`

	lines := func(s string) []string {
		return strings.Split(s, "\n")
	}

	if diff := cmp.Diff(lines(want), lines(string(sim.CoverageReport()))); diff != "" {
		t.Errorf("sim.CoverageReport() diff (-want +got):\n%s", diff)
	}
}
