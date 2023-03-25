package random

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/google/go-cmp/cmp"
)

func TestNextShuffler(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, 1)

	tests := []struct {
		total, seed uint64
	}{
		{20, 0},
		{20, 1},
		{20, 2},
		{20, 3},
		{15, 4},
		{25, 5},
		{50, 6},
	}

	// A common error in Fisher–Yates implementations is not including the
	// current index in the set of possibilities for the next shuffle, resulting
	// in Sattolo's algorithm. Although a test for this is probabilistic, the
	// chance of a false-negative on a stationary index is very low.
	var gotStationary bool

	for _, tt := range tests {
		t.Run(fmt.Sprintf("shuffling %d items with seed %d", tt.total, tt.seed), func(t *testing.T) {
			_, _, shuffler, err := DeployTestableNextShuffler(sim.Acc(0), sim, new(big.Int).SetUint64(tt.total))
			if err != nil {
				t.Fatalf("DeployTestableNextShuffler() error %v", err)
			}

			// Capture the random values used in the shuffle to (a) check invariants;
			// and (b) reimplement regular Fisher–Yates shuffle for comparison with the
			// contract output.
			rand := make(chan *TestableNextShufflerSwappedWith)
			shuffler.TestableNextShufflerFilterer.WatchSwappedWith(nil, rand)
			defer close(rand)

			if _, err := shuffler.Permute(sim.Acc(0), tt.seed); err != nil {
				t.Fatalf("Permute(%d) error %v", tt.seed, err)
			}

			runShuffling := func() []uint64 {
				var got []uint64
				for i := uint64(0); i < tt.total; i++ {
					n, err := shuffler.Permutation(nil, new(big.Int).SetUint64(i))
					if err != nil {
						t.Fatalf("Permutation(%d) error %v", i, err)
					}
					if !n.IsUint64() {
						t.Fatalf("Permutation(%d).IsUint64() = false; want true", i)
					}
					got = append(got, n.Uint64())
				}
				return got
			}
			got := runShuffling()

			gotShuffles := make([]int, tt.total)
			for i := uint64(0); i < tt.total; i++ {
				swap := <-rand
				gotShuffles[int(swap.Current.Uint64())] = int(swap.With.Uint64())
			}

			var want []uint64
			for i := uint64(0); i < uint64(tt.total); i++ {
				want = append(want, i)
			}
			for i, j := range gotShuffles {
				want[i], want[j] = want[j], want[i]
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Permutation diff compared to regular Fisher–Yates (-want +got):\n%s", diff)
			}

			for i, j := range gotShuffles {
				gotStationary = gotStationary || i == j
				if i > j {
					t.Errorf("Index %d Swapped with earlier index %d; want lookahead only", i, j)
				}
				if uint64(j) >= tt.total {
					t.Errorf("Index %d Swapped with out-of-range index %d; want within list of length %d", i, j, tt.total)
				}
			}

			shuffler.Reset(sim.Acc(0))
			got = runShuffling()
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Permutation diff compared to regular Fisher–Yates (-want +got):\n%s", diff)
			}

			inOrder := true
			for i, j := range gotShuffles {
				if i != j {
					inOrder = false
					break
				}
			}
			if inOrder {
				n := new(big.Int).MulRange(1, int64(tt.total))
				p := new(big.Float).Quo(big.NewFloat(1), new(big.Float).SetInt(n))
				t.Errorf("Shuffle was in order; this is very unlikely to happen by change (p=%.2e)", p)
			}

		})
	}

	if !gotStationary {
		t.Error("No shuffles kept the index in place; likely off-by-one error")
	}
}
