package raffle

import (
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
)

const (
	deployer = iota
)

func deploy(t *testing.T, maxWinners int64, cost *big.Int) (*ethtest.SimulatedBackend, *TestableRaffleRunner, *Raffle) {
	t.Helper()
	sim := ethtest.NewSimulatedBackendTB(t, 100)

	t.Logf("%d max winners; entry cost %d", maxWinners, cost)

	_, _, runner, err := DeployTestableRaffleRunner(sim.Acc(deployer), sim, big.NewInt(maxWinners), cost)
	if err != nil {
		t.Fatalf("DeployTestableRaffleRunner() error %v", err)
	}

	addr, err := runner.Raffle(nil)
	if err != nil {
		t.Fatalf("Runner.Raffle() error %v", err)
	}
	raffle, err := NewRaffle(addr, sim)
	if err != nil {
		t.Fatalf("NewRaffle(<address from runner>) error %v", err)
	}

	return sim, runner, raffle
}

func wantEntrantState(t *testing.T, raffle *Raffle, addr common.Address, entries, wins int64) {
	t.Helper()

	got, err := raffle.Entrants(nil, addr)
	if err != nil {
		t.Fatalf("Entrants(%s) error %v", addr, err)
	}

	want := struct{ Entries, Wins *big.Int }{
		big.NewInt(entries), big.NewInt(wins),
	}
	if diff := cmp.Diff(want, got, ethtest.Comparers()...); diff != "" {
		t.Errorf("Entrants(%s) diff (-want +got):\n%s", addr, diff)
	}
}

func TestEndToEnd(t *testing.T) {
	const (
		maxWinners  = 100
		cost        = 1 // Ether
		reserveFor  = 2
		reserveEach = 3
		entrants    = 50
	)
	sim, runner, raffle := deploy(t, maxWinners, eth.Ether(cost))

	var totalEntries int64

	for i := 0; i < reserveFor; i++ {
		totalEntries += reserveEach
		if _, err := runner.Reserve(sim.WithValueFrom(i, eth.Ether(cost*reserveEach)), big.NewInt(reserveEach)); err != nil {
			t.Fatalf("Reserve(%d) error %v", reserveEach, err)
		}
		wantEntrantState(t, raffle, sim.Acc(i).From, reserveEach, reserveEach)
	}

	for i := 0; i < entrants; i++ {
		entries := int64(i%4 + 1)
		totalEntries += entries

		entrant := sim.WithValueFrom(i, eth.Ether(cost*entries))
		if _, err := raffle.Enter(entrant, entrant.From, big.NewInt(entries)); err != nil {
			t.Fatalf("Enter(%d) error %v", entries, err)
		}

		wantEntries := entries
		var wantWins int64
		if i < reserveFor {
			wantEntries += reserveEach
			wantWins = reserveEach
		}
		wantEntrantState(t, raffle, sim.Acc(i).From, wantEntries, wantWins)
	}

	// Entropy needs to be anything non-zero as this is the check for it being
	// set.
	var entropy [32]byte
	entropy[0] = 1
	// The canEnter modifier requires that the contract isn't paused and that
	// entropy hasn't been set. If we allow entropy to be set without first
	// pausing then we risk a race condition that might revert some entries and
	// unfairly cost them gas.
	_, err := raffle.SetEntropy(sim.Acc(deployer), entropy)
	if diff := ethtest.RevertDiff(err, "Pausable: not paused"); diff != "" {
		t.Errorf("SetEntropy() when not paused; %s", diff)
	}
	if _, err := raffle.Pause(sim.Acc(deployer)); err != nil {
		t.Errorf("Pause() error %v", err)
	}

	noEntries := func(t *testing.T, when, wantRevertMsg string) {
		t.Helper()

		t.Run("no entries nor reservations when "+when, func(t *testing.T) {
			entrant := sim.WithValueFrom(0, eth.Ether(cost))

			_, err = raffle.Enter(entrant, entrant.From, big.NewInt(1))
			if diff := ethtest.RevertDiff(err, wantRevertMsg); diff != "" {
				t.Errorf("Enter() when paused; %s", diff)
			}

			_, err = runner.Reserve(entrant, big.NewInt(1))
			if diff := ethtest.RevertDiff(err, wantRevertMsg); diff != "" {
				t.Errorf("Reserve() when paused; %s", diff)
			}
		})
	}
	noEntries(t, "paused", "Raffle: closed")

	if _, err := raffle.SetEntropy(sim.Acc(deployer), entropy); err != nil {
		t.Fatalf("SetEntropy() error %v", err)
	}

	if _, err := raffle.Unpause(sim.Acc(deployer)); err != nil {
		t.Errorf("Unpause() error %v", err)
	}
	noEntries(t, "entropy set", "Raffle: entropy set")

	t.Run("no redeem before shuffle", func(t *testing.T) {
		_, err := raffle.Redeem(sim.Acc(0), sim.Acc(0).From)
		if diff := ethtest.RevertDiff(err, "Raffle: winners not assigned"); diff != "" {
			t.Errorf("Redeem() before call to Shuffle(); %s", diff)
		}
	})

	tx, err := raffle.Shuffle(sim.Acc(0))
	if err != nil {
		t.Fatalf("Shuffle() error %v", err)
	}
	ethtest.LogGas(t, tx, "Raffle shuffle")

	for i := 0; i < entrants; i++ {
		if _, err := raffle.Redeem(sim.Acc(i), sim.Acc(i).From); err != nil {
			t.Fatalf("Redeem() error %v", err)
		}
	}

	t.Run("refunds", func(t *testing.T) {
		gotTotal := big.NewInt(0)

		refunds, err := raffle.RaffleFilterer.FilterRefund(nil, nil)
		if err != nil {
			t.Fatalf("FilterRefund() for logs; error %v", err)
		}

		for refunds.Next() {
			ev := refunds.Event
			gotTotal.Add(gotTotal, ev.Amount)
		}
		if err := refunds.Error(); err != nil {
			t.Errorf("%T.Error() got %v; want nil", refunds, err)
		}

		wantTotal := eth.Ether((totalEntries - maxWinners) * cost)
		if gotTotal.Cmp(wantTotal) != 0 {
			t.Errorf("Total refunds, determined via logs; got %d; want %d", gotTotal, wantTotal)
		}
	})

	t.Run("wins", func(t *testing.T) {
		gotTotal := big.NewInt(0)

		// TODO: when generics are released in Go 1.18 (~Feb 2022), abstract
		// an event-mapping function.
		redemptions, err := raffle.RaffleFilterer.FilterRedemption(nil, nil)
		if err != nil {
			t.Fatalf("FilterRedemption() for logs; error %v", err)
		}

		for redemptions.Next() {
			ev := redemptions.Event
			gotTotal.Add(gotTotal, ev.Wins)

			got, err := runner.Wins(nil, ev.Entrant)
			if err != nil {
				t.Errorf("runner.Wins() error %v", err)
				continue
			}
			if want := ev.Wins; got.Cmp(want) != 0 {
				t.Errorf("Wins recorded by RaffleRunner; got %d; want %d as reported by event logs", got, want)
			}
		}
		if err := redemptions.Error(); err != nil {
			t.Errorf("%T.Error() got %v; want nil", redemptions, err)
		}

		if gotTotal.Cmp(big.NewInt(maxWinners)) != 0 {
			t.Errorf("Total wins, determined via logs; got %d; want %d", gotTotal, maxWinners)
		}
	})
}

func TestReservationLimits(t *testing.T) {
	const (
		maxWinners = 10
		cost       = 1
	)
	sim, runner, raffle := deploy(t, maxWinners, eth.Ether(cost))

	t.Run("no direct reservation", func(t *testing.T) {
		const acc = 1
		_, err := raffle.Reserve(sim.WithValueFrom(acc, eth.Ether(cost)), sim.Acc(acc).From, big.NewInt(1))
		if diff := ethtest.RevertDiff(err, "Raffle: only runner can reserve"); diff != "" {
			t.Errorf("Direct call toRaffle.Reserve() %s", diff)
		}
		wantEntrantState(t, raffle, sim.Acc(acc).From, 0, 0)
	})

	t.Run("reservation via runner", func(t *testing.T) {
		const acc = 0
		for i := 0; i < 2; i++ {
			from := sim.WithValueFrom(acc, eth.Ether(cost*maxWinners/2))
			if _, err := runner.Reserve(from, big.NewInt(maxWinners/2)); err != nil {
				t.Fatalf("RaffleRunner.Reserve() error %v", err)
			}
			wantEntrantState(t, raffle, sim.Acc(acc).From, maxWinners/2*int64(i+1), maxWinners/2*int64(i+1))
		}

		t.Run("no over subscription", func(t *testing.T) {
			_, err := runner.Reserve(sim.WithValueFrom(acc, eth.Ether(cost)), big.NewInt(1))
			if diff := ethtest.RevertDiff(err, "Raffle: too many reserved"); diff != "" {
				t.Errorf("RaffleRunner.Reserve() when all reserved; %s", diff)
			}
			wantEntrantState(t, raffle, sim.Acc(acc).From, maxWinners, maxWinners)
		})
	})
}
