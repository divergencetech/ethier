package random

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/go-cmp/cmp"
)

func deployPRNG(t *testing.T) (*ethtest.SimulatedBackend, *TestablePRNG) {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, 1)
	_, _, prng, err := DeployTestablePRNG(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployTestablePRNG() error %v", err)
	}

	return sim, prng
}

func deployCSPRNG(t *testing.T) (*ethtest.SimulatedBackend, *TestableCSPRNG) {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, 1)
	_, _, prng, err := DeployTestableCSPRNG(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployTestablePRNG() error %v", err)
	}

	return sim, prng
}

type TestingContract interface {
	Sample(opts *bind.CallOpts, seed [32]byte, bits uint16, n uint16) ([]*big.Int, error)
	BitLength(opts *bind.CallOpts, n *big.Int) (*big.Int, error)
	ReadLessThan(opts *bind.CallOpts, seed [32]byte, max *big.Int, n uint16) ([]*big.Int, error)
	TestStoreAndLoad(opts *bind.TransactOpts, seed [32]byte, bits uint16, beforeStore uint16) (*types.Transaction, error)
}

type Implementation struct {
	name string
	sim  *ethtest.SimulatedBackend
	prng TestingContract
}

func getImplementations(t *testing.T) []Implementation {
	var impl []Implementation
	{
		sim, prng := deployPRNG(t)
		impl = append(impl, Implementation{
			name: "PRNG",
			sim:  sim,
			prng: prng,
		})
	}
	{
		sim, prng := deployCSPRNG(t)
		impl = append(impl, Implementation{
			name: "CSPRNG",
			sim:  sim,
			prng: prng,
		})
	}

	return impl
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func normalCDF(z float64) float64 {
	return 0.5 + 0.5*math.Erf(z/math.Sqrt2)
}

func TestReadProbabilistic(t *testing.T) {
	tests := []struct {
		seedFrom                  string
		bitsPerSample, numSamples uint64
	}{
		{
			seedFrom:      "a",
			bitsPerSample: 8,
			numSamples:    5_000,
		},
		{
			seedFrom:      "b",
			bitsPerSample: 5,
			numSamples:    5_000,
		},
		{
			seedFrom:      "c",
			bitsPerSample: 7,
			numSamples:    10_000,
		},
		{
			seedFrom:      "d",
			bitsPerSample: 3,
			numSamples:    7_500,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("seed keccak256(%q) %d samples of %d bits", tt.seedFrom, tt.numSamples, tt.bitsPerSample)
		t.Run(name, func(t *testing.T) {
			var seed [32]byte
			copy(seed[:], crypto.Keccak256([]byte(tt.seedFrom)))

			for _, ii := range getImplementations(t) {
				t.Run(ii.name, func(t *testing.T) {

					samples, err := ii.prng.Sample(nil, seed, uint16(tt.bitsPerSample), uint16(tt.numSamples))
					if err != nil {
						t.Fatalf("Sample(%x, 8 bits, 1000 samples) error %v", seed, err)
					}

					var min, max, sum uint64
					var onesCount int

					numMax := (1 << tt.bitsPerSample) - 1
					distMean := float64(numMax) / 2.
					distVar := float64(numMax*numMax) / 12.

					min = math.MaxUint64

					var cusum, cusumMax int64

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

						onesCount += bits.OnesCount64(u)

						for j := uint64(0); j < tt.bitsPerSample; j++ {
							z := 2*int64((u>>j)&1) - 1
							cusum += z
							abs := abs(cusum)
							if abs > cusumMax {
								cusumMax = abs
							}
						}
					}

					t.Run("MinMax", func(t *testing.T) {
						// Extremely unlikely that neither the min nor the max are
						// returned given the high number of samples and the small
						// sampling spaces.
						if min != 0 {
							t.Errorf("Min = %d; want 0", min)
						}
						if want := uint64(1<<tt.bitsPerSample) - 1; max != want {
							t.Errorf("Max = %d; want %d", max, want)
						}
					})

					t.Run("Sum", func(t *testing.T) {

						expectedSum := distMean * float64(tt.numSamples)
						// The variance of as sum of uniform random variables in
						// [0,sampleMax] is given by the sum of the individual variances
						twoSigmaRelDev := 2 * math.Sqrt(distVar*float64(tt.numSamples)) / expectedSum

						if dev := 1 - float64(expectedSum)/float64(sum); math.Abs(dev) > twoSigmaRelDev {
							t.Errorf("Sum = %s deviates by %.2f%% (>%.2f%% abs = 2 sigma) from expected sum %s", humanize.Comma(int64(sum)), dev*100, twoSigmaRelDev*100, humanize.Comma(int64(expectedSum)))
						}
					})

					t.Run("Monobit Frequency", func(t *testing.T) {
						expectedMonocount := float64(tt.numSamples*tt.bitsPerSample) / 2
						// Follows a binomial distribution. Variance is nSamples/4
						twoSigmaRelDev := 2. / math.Sqrt(float64(tt.numSamples*tt.bitsPerSample))

						if dev := 1 - float64(expectedMonocount)/float64(onesCount); math.Abs(dev) > twoSigmaRelDev {
							t.Errorf("One count = %s deviates by %.2f%% (>%.2f%% abs = 2 sigma) from expected sum %s", humanize.Comma(int64(onesCount)), dev*100, twoSigmaRelDev*100, humanize.Comma(int64(expectedMonocount)))
						}
					})

					t.Run("Monobit CUSUM", func(t *testing.T) {
						// See Sect. 2.13 in https://nvlpubs.nist.gov/nistpubs/legacy/sp/nistspecialpublication800-22r1a.pdf
						p := float64(1)
						n := tt.numSamples * tt.bitsPerSample
						z := cusumMax
						sqrtn := math.Sqrt(float64(n))
						kStart := int64(0.25 * (-float64(n)/float64(z) + 1))
						kEnd := int64(0.25 * (float64(n)/float64(z) - 1))

						for k := kStart; k <= kEnd; k++ {
							p -= normalCDF(float64((4*k+1)*z)/sqrtn) - normalCDF(float64((4*k-1)*z)/sqrtn)
						}

						kStart = int64(0.25 * (-float64(n)/float64(z) - 3))
						for k := kStart; k <= kEnd; k++ {
							p += normalCDF(float64((4*k+3)*z)/sqrtn) - normalCDF(float64((4*k+1)*z)/sqrtn)
						}

						if p < 0.01 {
							t.Errorf("Cusum test failed with p=%.2f < 0.01", p)
						}
					})

				})

			}
		})
	}
}

func TestReadInternalStatePRNG(t *testing.T) {
	tests := []struct {
		seedFrom                  string
		bitsPerSample, numSamples uint64
	}{
		{
			seedFrom:      "a",
			bitsPerSample: 8,
			numSamples:    5_000,
		},
		{
			seedFrom:      "b",
			bitsPerSample: 5,
			numSamples:    5_000,
		},
		{
			seedFrom:      "c",
			bitsPerSample: 7,
			numSamples:    10_000,
		},
		{
			seedFrom:      "d",
			bitsPerSample: 3,
			numSamples:    7_500,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("seed keccak256(%q) %d samples of %d bits", tt.seedFrom, tt.numSamples, tt.bitsPerSample)
		t.Run(name, func(t *testing.T) {
			_, prng := deployPRNG(t)

			var seed [32]byte
			copy(seed[:], crypto.Keccak256([]byte(tt.seedFrom)))

			state, err := prng.SampleState(nil, seed, uint16(tt.bitsPerSample), uint16(tt.numSamples))
			if err != nil {
				t.Fatalf("Sample(%x, 8 bits, 1000 samples) error %v", seed, err)
			}

			// // This test is primarily in place to confirm that the assembly
			// // implementation works as intended when converted to a
			// // high-level, more readable equivalent.

			seedInt := new(big.Int).SetBytes(seed[:])
			two128 := new(big.Int).Lsh(big.NewInt(1), 128)

			// first 128 bits of seed
			carry := new(big.Int).Rsh(seedInt, 128)

			// last 128 bits of seed
			number := new(big.Int).Mod(seedInt, two128)

			factor := new(big.Int).Sub(two128, big.NewInt(10408))

			// Performing MWC updates
			nRefills := (tt.numSamples * tt.bitsPerSample) / 128
			for k := uint64(0); k < nRefills; k++ {
				tmp := new(big.Int).Mul(factor, number)
				tmp.Add(tmp, carry)

				carry.Rsh(tmp, 128)
				number.Mod(tmp, two128)
			}

			carry.Lsh(carry, 128)
			entropy := new(big.Int).Add(carry, number)

			wantState := TestablePRNGState{
				Entropy: entropy,
				Remain:  big.NewInt(128 - int64(tt.bitsPerSample*tt.numSamples)%128),
			}

			if diff := cmp.Diff(wantState, state, ethtest.Comparers()...); diff != "" {
				t.Errorf("After %d samples of %d bits each; internal state diff (-want +got):\n%s", tt.numSamples, tt.bitsPerSample, diff)
			}
		})
	}
}

func TestReadInternalStateCSPRNG(t *testing.T) {
	tests := []struct {
		seedFrom                  string
		bitsPerSample, numSamples uint64
	}{
		{
			seedFrom:      "a",
			bitsPerSample: 8,
			numSamples:    5_000,
		},
		{
			seedFrom:      "b",
			bitsPerSample: 5,
			numSamples:    5_000,
		},
		{
			seedFrom:      "c",
			bitsPerSample: 7,
			numSamples:    10_000,
		},
		{
			seedFrom:      "d",
			bitsPerSample: 3,
			numSamples:    7_500,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("seed keccak256(%q) %d samples of %d bits", tt.seedFrom, tt.numSamples, tt.bitsPerSample)
		t.Run(name, func(t *testing.T) {
			_, prng := deployCSPRNG(t)

			var seed [32]byte
			copy(seed[:], crypto.Keccak256([]byte(tt.seedFrom)))

			state, err := prng.SampleState(nil, seed, uint16(tt.bitsPerSample), uint16(tt.numSamples))
			if err != nil {
				t.Fatalf("Sample(%x, 8 bits, 1000 samples) error %v", seed, err)
			}

			// This test is primarily in place to confirm that the assembly
			// implementation works as intended when converted to a
			// high-level, more readable equivalent.
			wantState := TestableCSPRNGState{
				Seed:    new(big.Int).SetBytes(seed[:]),
				Counter: big.NewInt(int64(tt.bitsPerSample*tt.numSamples/256 + 1)),
				Remain:  big.NewInt(256 - int64(tt.bitsPerSample*tt.numSamples)%256),
			}

			// PRNG increments the counter and then fills the entropy pool
			// with keccak256(seed||counter).
			entropy := new(big.Int).Lsh(wantState.Seed, 256)
			entropy.Add(entropy, wantState.Counter)
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
	}
}

func TestBitLength(t *testing.T) {

	for _, ii := range getImplementations(t) {
		t.Run(ii.name, func(t *testing.T) {

			for _, in := range []uint64{0, 1, 2, 3, 4, 5, 63, 64, 127, 128, 255, 256, math.MaxUint64} {
				want := bits.Len64(in)

				got, err := ii.prng.BitLength(nil, new(big.Int).SetUint64(in))
				if err != nil {
					t.Errorf("BitLength(%d) error %v", in, err)
					continue
				}

				if got.Cmp(big.NewInt(int64(want))) != 0 {
					t.Errorf("BitLength(%d) got %d; want %d", in, got, want)
				}

			}
		})
	}
}

func TestReadLessThan(t *testing.T) {

	for _, ii := range getImplementations(t) {
		t.Run(ii.name, func(t *testing.T) {

			const n = uint16(1e4)

			for _, max := range []uint64{1, 2, 3, 7, 8, 127, 128, 1023, 1024} {
				t.Run(fmt.Sprintf("%d samples < %d", n, max), func(t *testing.T) {
					var seed [32]byte

					bigMax := new(big.Int).SetUint64(max)
					got, err := ii.prng.ReadLessThan(nil, seed, bigMax, n)
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
		})
	}
}

func TestStoreAndLoad(t *testing.T) {
	// The tests for store() and loadSource() are performed by TestablePRNG
	// itself, by asserting that a series of calls to read() are identical. This
	// Go test merely triggers the function and checks that there's no error.

	tests := []struct {
		bits, beforeStore uint16
	}{
		{
			// Used bits mod 256 == 0
			bits:        16,
			beforeStore: 32,
		},
		{
			// Used bits mod 256 != 0
			bits:        19,
			beforeStore: 17,
		},
		{
			bits:        126,
			beforeStore: 10,
		},
	}

	for i, tt := range tests {
		var seed [32]byte
		seed[31] = byte(i + 1)

		for _, ii := range getImplementations(t) {
			t.Run(ii.name, func(t *testing.T) {

				if _, err := ii.prng.TestStoreAndLoad(ii.sim.Acc(0), seed, tt.bits, tt.beforeStore); err != nil {
					t.Errorf("StoreAndLoad(%#x, %d, %d) error %v", seed, tt.bits, tt.beforeStore, err)
				}
			})
		}
	}
}
