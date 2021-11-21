package random_test

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/go-cmp/cmp"
)

func deploy(t *testing.T) *TestablePRNG {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, 1)
	_, _, prng, err := DeployTestablePRNG(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployTestablePRNG() error %v", err)
	}

	return prng
}

func TestRead(t *testing.T) {
	prng := deploy(t)
	tests := []struct {
		seedFrom                  string
		bitsPerSample, numSamples uint64
	}{
		{
			seedFrom:      "a",
			bitsPerSample: 8,
			numSamples:    5000,
		},
		{
			seedFrom:      "b",
			bitsPerSample: 5,
			numSamples:    5000,
		},
		{
			seedFrom:      "c",
			bitsPerSample: 7,
			numSamples:    5000,
		},
		{
			seedFrom:      "d",
			bitsPerSample: 3,
			numSamples:    7500,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("seed keccak256(%q) %d samples of %d bits", tt.seedFrom, tt.numSamples, tt.bitsPerSample)
		t.Run(name, func(t *testing.T) {
			var seed [32]byte
			copy(seed[:], crypto.Keccak256([]byte(tt.seedFrom)))

			samples, state, err := prng.Sample(nil, seed, uint16(tt.bitsPerSample), uint16(tt.numSamples))
			if err != nil {
				t.Fatalf("Sample(%x, 8 bits, 1000 samples) error %v", seed, err)
			}

			t.Run("probabilistic tests", func(t *testing.T) {
				var min, max, sum uint64
				min = math.MaxUint64
				for i, s := range samples {
					if !s.IsUint64() {
						t.Fatalf("sample[%d].IsUint64() got false; want true", i)
					}

					u := s.Uint64()
					sum += u
					if u < min {
						min = u
					}
					if u > max {
						max = u
					}
				}

				// Extremely unlikely that neither the min nor the max are
				// returned given the high number of samples and the small
				// sampling spaces.
				if min != 0 {
					t.Errorf("Min = %d; want 0", min)
				}
				if want := uint64(1<<tt.bitsPerSample) - 1; max != want {
					t.Errorf("Max = %d; want %d", max, want)
				}

				expectedSum := (1<<(tt.bitsPerSample-1))*tt.numSamples - tt.numSamples/2
				if dev := 1 - float64(expectedSum)/float64(sum); math.Abs(dev) > 0.01 {
					t.Errorf("Sum = %s deviates by %.2f%% (>1%% abs) from expected sum %s", humanize.Comma(int64(sum)), dev*100, humanize.Comma(int64(expectedSum)))
				}
			})

			t.Run("internal state", func(t *testing.T) {
				// This test is primarily in place to confirm that the assembly
				// implementation works as intended when converted to a
				// high-level, more readable equivalent.
				wantState := TestablePRNGState{
					Seed:    new(big.Int).SetBytes(seed[:]),
					Counter: big.NewInt(int64(tt.bitsPerSample*tt.numSamples/256 + 1)),
					Remain:  big.NewInt(256 - int64(tt.bitsPerSample*tt.numSamples)%256),
				}

				// PRNG fills the entropy pool with keccak256(seed||counter) and
				// then increments the counter. We must therefore decrease the
				// counter when regenerating the expected entropy.
				entropy := new(big.Int).Lsh(wantState.Seed, 256)
				entropy.Add(entropy, new(big.Int).Sub(wantState.Counter, big.NewInt(1)))
				// It's important not to use entropy.Bytes() here as we MUST
				// have exactly 64 bytes to be hashed.
				buf := make([]byte, 64)
				entropy.FillBytes(buf)
				// Some of the entropy has depleted; we know how much from
				// Remain.
				wantState.Entropy = new(big.Int).Rsh(
					new(big.Int).SetBytes(crypto.Keccak256(buf)),
					uint(256-wantState.Remain.Uint64()),
				)

				if diff := cmp.Diff(wantState, state, ethtest.Comparers()...); diff != "" {
					t.Errorf("After %d samples of %d bits each; internal state diff (-want +got):\n%s", tt.numSamples, tt.bitsPerSample, diff)
				}
			})
		})
	}
}

func TestBitLength(t *testing.T) {
	prng := deploy(t)

	for _, in := range []uint64{0, 1, 2, 3, 4, 5, 63, 64, 127, 128, 255, 256, math.MaxUint64} {
		want := bits.Len64(in)

		got, err := prng.BitLength(nil, new(big.Int).SetUint64(in))
		if err != nil {
			t.Errorf("BitLength(%d) error %v", in, err)
			continue
		}

		if got.Cmp(big.NewInt(int64(want))) != 0 {
			t.Errorf("BitLength(%d) got %d; want %d", in, got, want)
		}

	}
}

func TestReadLessThan(t *testing.T) {
	prng := deploy(t)

	const n = uint16(1e4)

	for _, max := range []uint64{1, 2, 3, 7, 8, 127, 128, 1023, 1024} {
		t.Run(fmt.Sprintf("%d samples < %d", n, max), func(t *testing.T) {
			var seed [32]byte

			bigMax := new(big.Int).SetUint64(max)
			got, err := prng.ReadLessThan(nil, seed, bigMax, n)
			if err != nil {
				t.Fatalf("ReadLessThan() error %v", err)
			}

			for _, s := range got {
				if s.Cmp(bigMax) != -1 {
					t.Errorf("ReadLessThan(%d) returned sample %d out of range", max, s)
				}
			}
		})
	}
}
