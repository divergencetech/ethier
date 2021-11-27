package utils

import (
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/h-fam/errdiff"
)

func TestDynamicBuffer(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, 3)
	_, _, dynBuf, err := DeployTestableDynamicBuffer(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployTestableDynamicBuffer() error %v", err)
	}
	const (
		testStr        = "This is a really long test string that we want to use."
		testStrShort   = "This is a short string"
		testStr32      = "This test string is 32bytes long"
		outOfBoundsMsg = "DynamicBuffer: Appending out of bounds."
	)

	tests := []struct {
		name           string
		capacity       int64
		appendString   string
		repetitions    int64
		errDiffAgainst interface{}
	}{
		{
			name:         "Single append",
			capacity:     int64(len(testStr)),
			appendString: testStr,
			repetitions:  1,
		},
		{
			name:           "Double append out-of-bound",
			capacity:       int64(len(testStr)),
			appendString:   testStr,
			repetitions:    2,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:         "Mutliple append",
			capacity:     420 * int64(len(testStr)),
			appendString: testStr,
			repetitions:  420,
		},
		{
			name:           "Mutliple append out-of-bound",
			capacity:       420 * int64(len(testStr)),
			appendString:   testStr,
			repetitions:    421,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:         "Single append 32B",
			capacity:     int64(len(testStr32)),
			appendString: testStr32,
			repetitions:  1,
		},
		{
			name:           "Double append 32B out-of-bound",
			capacity:       int64(len(testStr32)),
			appendString:   testStr32,
			repetitions:    2,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:         "Mutliple append 32B",
			capacity:     420 * int64(len(testStr32)),
			appendString: testStr32,
			repetitions:  420,
		},
		{
			name:           "Mutliple append 32B out-of-bound",
			capacity:       420 * int64(len(testStr32)),
			appendString:   testStr32,
			repetitions:    421,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:         "Single append short",
			capacity:     int64(len(testStrShort)),
			appendString: testStrShort,
			repetitions:  1,
		},
		{
			name:           "Double append short out-of-bound",
			capacity:       int64(len(testStrShort)),
			repetitions:    2,
			appendString:   testStrShort,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:         "Mutliple append short",
			capacity:     420 * int64(len(testStrShort)),
			appendString: testStrShort,
			repetitions:  420,
		},
		{
			name:           "Mutliple append short out-of-bound",
			capacity:       420 * int64(len(testStrShort)),
			appendString:   testStrShort,
			repetitions:    421,
			errDiffAgainst: outOfBoundsMsg,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dynBuf.AllocateAndAppendRepeated(nil, big.NewInt(tt.capacity), tt.appendString, big.NewInt(tt.repetitions))

			if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
				t.Fatalf("AllocateAndAppendRepeated(%d, %q, %d) %s", tt.capacity, tt.appendString, tt.repetitions, diff)
			}

			if tt.errDiffAgainst != nil {
				return
			}

			want := ""
			for r := int64(0); r < tt.repetitions; r++ {
				want = want + tt.appendString
			}

			if got != want {
				t.Errorf("AllocateAndAppendRepeated(%d, %q, %d) got %q; want %q", tt.capacity, tt.appendString, tt.repetitions, got, want)
			}
		})
	}
}
