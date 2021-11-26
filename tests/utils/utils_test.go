package utils

import (
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/google/go-cmp/cmp"
	"github.com/h-fam/errdiff"
)

func TestOwnerPausable(t *testing.T) {
	const (
		owner = iota
		vandal
		newOwner
	)

	sim := ethtest.NewSimulatedBackendTB(t, 3)
	_, _, op, err := DeployOwnerPausable(sim.Acc(owner), sim)
	if err != nil {
		t.Fatalf("DeployOwnerPausable() error %v", err)
	}

	got, err := op.Owner(nil)
	if err != nil {
		t.Fatalf("Owner() error %v", err)
	}
	if want := sim.Acc(owner).From; !cmp.Equal(got, want) {
		t.Fatalf("Owner() got %s; want %s", got, want)
	}

	const notOwnerMsg = "Ownable: caller is not the owner"

	// Tests are deliberately not hermetic and some actions aren't tests per se
	// but setup for the next action (e.g. transferring ownership).
	tests := []struct {
		desc           string
		action         func() error
		errDiffAgainst interface{}
	}{
		{
			desc: "Pause() as vandal; ",
			action: func() error {
				_, err := op.Pause(sim.Acc(vandal))
				return err
			},
			errDiffAgainst: notOwnerMsg,
		},
		{
			desc: "Pause() as owner; ",
			action: func() error {
				_, err := op.Pause(sim.Acc(owner))
				return err
			},
		},
		{
			desc: "Unpause() as future owner; ",
			action: func() error {
				_, err := op.Pause(sim.Acc(newOwner))
				return err
			},
			errDiffAgainst: notOwnerMsg,
		},
		{
			desc: "TransferOwnership()",
			action: func() error {
				_, err := op.TransferOwnership(sim.Acc(owner), sim.Acc(newOwner).From)
				return err
			},
		},
		{
			desc: "Unpause() as new owner; ",
			action: func() error {
				_, err := op.Unpause(sim.Acc(newOwner))
				return err
			},
		},
	}

	for _, tt := range tests {
		if diff := errdiff.Check(tt.action(), tt.errDiffAgainst); diff != "" {
			t.Fatalf("%s %s", tt.desc, diff)
		}
	}
}

func TestDynamicBuffer(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, 3)
	_, _, dynBuf, err := DeployTestableDynamicBuffer(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployTestableDynamicBuffer() error %v", err)
	}

	const testStr = "This is a really long test string that we want to use."
	const testStrShort = "This is a short string"
	const testStr32 = "This test string is 32bytes long"
	const outOfBoundsMsg = "DynamicBuffer: Appending out of bounds."

	tests := []struct {
		desc           string
		capacity       int
		appendString   string
		repetitions    int
		errDiffAgainst interface{}
	}{
		{
			desc:         "Single append;",
			capacity:     len(testStr),
			repetitions:  1,
			appendString: testStr,
		},
		{
			desc:           "Double append out-of-bound;",
			capacity:       len(testStr),
			repetitions:    2,
			errDiffAgainst: outOfBoundsMsg,
			appendString:   testStr,
		},
		{
			desc:         "Mutliple append;",
			capacity:     420 * len(testStr),
			repetitions:  420,
			appendString: testStr,
		},
		{
			desc:           "Mutliple append out-of-bound;",
			capacity:       420 * len(testStr),
			repetitions:    421,
			errDiffAgainst: outOfBoundsMsg,
			appendString:   testStr,
		},
		{
			desc:         "Single append 32B;",
			capacity:     len(testStr32),
			repetitions:  1,
			appendString: testStr32,
		},
		{
			desc:           "Double append 32B out-of-bound;",
			capacity:       len(testStr32),
			repetitions:    2,
			errDiffAgainst: outOfBoundsMsg,
			appendString:   testStr32,
		},
		{
			desc:         "Mutliple append 32B;",
			capacity:     420 * len(testStr32),
			repetitions:  420,
			appendString: testStr32,
		},
		{
			desc:           "Mutliple append 32B out-of-bound;",
			capacity:       420 * len(testStr32),
			repetitions:    421,
			errDiffAgainst: outOfBoundsMsg,
			appendString:   testStr32,
		},
		{
			desc:         "Single append short;",
			capacity:     len(testStrShort),
			repetitions:  1,
			appendString: testStrShort,
		},
		{
			desc:           "Double append short out-of-bound;",
			capacity:       len(testStrShort),
			repetitions:    2,
			errDiffAgainst: outOfBoundsMsg,
			appendString:   testStrShort,
		},
		{
			desc:         "Mutliple append short;",
			capacity:     420 * len(testStrShort),
			repetitions:  420,
			appendString: testStrShort,
		},
		{
			desc:           "Mutliple append short out-of-bound;",
			capacity:       420 * len(testStrShort),
			repetitions:    421,
			errDiffAgainst: outOfBoundsMsg,
			appendString:   testStrShort,
		},
	}

	for _, tt := range tests {
		buffer, err := dynBuf.AllocateAndAppendRepeated(nil, big.NewInt(int64(tt.capacity)), tt.appendString, big.NewInt(int64(tt.repetitions)))

		if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
			t.Fatalf("%s %s", tt.desc, diff)
		}

		if tt.errDiffAgainst == "" {
			should := ""
			for rep := 0; rep < tt.repetitions; rep++ {
				should = should + tt.appendString
			}

			if buffer != should {
				t.Fatalf("%s: Expected %s instead of %s", tt.desc, should, buffer)
			}
		}
	}
}
