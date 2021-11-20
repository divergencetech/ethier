package ethtest

import (
	"math/big"
	"testing"
)

func TestFastForward(t *testing.T) {
	sim := NewSimulatedBackendTB(t, 1)

	if got, want := sim.BlockNumber().Int64(), int64(1); got != want {
		t.Fatalf("%T.BlockNumber() when fresh; got %d; want %d", sim, got, want)
	}

	// Tests are deliberately not hermetic as we want to confirm when
	// FastForward() returns false.
	tests := []struct {
		ffToBlock      int64
		want           bool
		wantBlockAfter int64
	}{
		{10, true, 10},
		{10, false, 10},
		{5, false, 10},
		{11, true, 11},
		{20, true, 20},
		{1, false, 20},
	}

	for _, tt := range tests {
		before := sim.BlockNumber()
		if got := sim.FastForward(big.NewInt(tt.ffToBlock)); got != tt.want {
			t.Errorf("%T.FastForward(%d) when on block %d; got %t want %t", sim, tt.ffToBlock, before, got, tt.want)
		}

		if got := sim.BlockNumber().Int64(); got != tt.wantBlockAfter {
			// Fatal because all future tests are invalid as the state is
			// incorrect.
			t.Fatalf("%T.BlockNumber() after FastForward(%d) when on already at block %d; got %d want %d", sim, tt.ffToBlock, before, got, tt.wantBlockAfter)
		}
	}
}
