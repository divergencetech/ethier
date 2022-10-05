package utils

import (
	"math/big"
	"strings"
	"testing"

	b64 "encoding/base64"

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

func base64Len(data string) int64 {
	return int64((len(data) + 2) / 3 * 4)
}

func base64LenUnpadded(data string) int64 {
	return base64Len(data) - int64(len(data)%3)
}

func TestDynamicBufferBase64(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, 3)
	_, _, dynBuf, err := DeployTestableDynamicBuffer(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployTestableDynamicBuffer() error %v", err)
	}
	const (
		testStrWithoutPadding = "This is a really long test string that with padding"
		testStrWithPadding    = "This is a short string"
		outOfBoundsMsg        = "DynamicBuffer: Appending out of bounds."
	)

	if base64Len(testStrWithoutPadding) != base64LenUnpadded(testStrWithoutPadding) {
		t.Fatalf("Test string has incorrect length: %s", testStrWithoutPadding)
	}

	if base64Len(testStrWithPadding) == base64LenUnpadded(testStrWithPadding) {
		t.Fatalf("Test string has incorrect length: %s", testStrWithPadding)
	}

	tests := []struct {
		name           string
		capacity       int64
		appendString   string
		repetitions    int64
		fileSafe       bool
		noPadding      bool
		errDiffAgainst interface{}
	}{
		{
			name:         "Single append",
			capacity:     base64Len(testStrWithoutPadding),
			appendString: testStrWithoutPadding,
			repetitions:  1,
		},
		{
			name:           "Double append out-of-bound",
			capacity:       base64Len(testStrWithoutPadding),
			appendString:   testStrWithoutPadding,
			repetitions:    2,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:         "Mutliple append",
			capacity:     420 * base64Len(testStrWithoutPadding),
			appendString: testStrWithoutPadding,
			repetitions:  420,
		},
		{
			name:           "Mutliple append out-of-bound",
			capacity:       420 * base64Len(testStrWithoutPadding),
			appendString:   testStrWithoutPadding,
			repetitions:    421,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:         "Single append short",
			capacity:     base64Len(testStrWithPadding),
			appendString: testStrWithPadding,
			repetitions:  1,
		},
		{
			name:           "Double append short out-of-bound",
			capacity:       base64Len(testStrWithPadding),
			repetitions:    2,
			appendString:   testStrWithPadding,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:         "Mutliple append short",
			capacity:     420 * base64Len(testStrWithPadding),
			appendString: testStrWithPadding,
			repetitions:  420,
		},
		{
			name:           "Mutliple append short out-of-bound",
			capacity:       420 * base64Len(testStrWithPadding),
			appendString:   testStrWithPadding,
			repetitions:    421,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:           "Mutliple append short out-of-bound",
			capacity:       420 * base64Len(testStrWithPadding),
			appendString:   testStrWithPadding,
			repetitions:    421,
			errDiffAgainst: outOfBoundsMsg,
		},
		{
			name:         "Single (unpadded) append with noPadding ",
			capacity:     base64Len(testStrWithoutPadding),
			appendString: testStrWithoutPadding,
			repetitions:  1,
			noPadding:    true,
		},
		{
			name:         "Double (unpadded) append with noPadding ",
			capacity:     2 * base64Len(testStrWithoutPadding),
			appendString: testStrWithoutPadding,
			repetitions:  2,
			noPadding:    true,
		},
		{
			name:         "Single (padded) append with noPadding ",
			capacity:     base64LenUnpadded(testStrWithPadding),
			appendString: testStrWithPadding,
			repetitions:  1,
			noPadding:    true,
		},
		{
			name:         "Mutliple (padded) append with noPadding ",
			capacity:     42 * base64LenUnpadded(testStrWithPadding),
			appendString: testStrWithPadding,
			repetitions:  42,
			noPadding:    true,
		},
		{
			name:         "Mutliple append file safe",
			capacity:     42 * base64Len(testStrWithoutPadding),
			appendString: testStrWithoutPadding,
			repetitions:  42,
			fileSafe:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dynBuf.AllocateAndAppendRepeatedBase64(nil, big.NewInt(tt.capacity), []byte(tt.appendString), big.NewInt(tt.repetitions), tt.fileSafe, tt.noPadding)

			if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
				t.Fatalf("AllocateAndAppendRepeatedBase64(%d, %q, %d, %v, %v) %s", tt.capacity, tt.appendString, tt.repetitions, tt.fileSafe, tt.noPadding, diff)
			}

			if tt.errDiffAgainst != nil {
				return
			}

			want := ""
			for r := int64(0); r < tt.repetitions; r++ {
				want = want + b64.StdEncoding.EncodeToString([]byte(tt.appendString))
			}

			if tt.noPadding {
				want = strings.ReplaceAll(want, "=", "")
			}

			if tt.fileSafe {
				want = strings.ReplaceAll(want, "+", "-")
				want = strings.ReplaceAll(want, "/", "_")
			}

			if got != want {
				t.Errorf("AllocateAndAppendRepeatedBase64(%d, %q, %d, %v, %v) got %q; want %q", tt.capacity, tt.appendString, tt.repetitions, tt.fileSafe, tt.noPadding, got, want)
			}
		})
	}
}
