package utils

import (
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
	if want := sim.Addr(owner); !cmp.Equal(got, want) {
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
				_, err := op.TransferOwnership(sim.Acc(owner), sim.Addr(newOwner))
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
