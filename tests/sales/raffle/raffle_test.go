package raffle

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
)

const (
	deployer = iota
	vandal
)

func deploy(t *testing.T, maxWinners int64, cost *big.Int) (*ethtest.SimulatedBackend, *TestableRaffleRunner, *Raffle, common.Address) {
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

	return sim, runner, raffle, addr
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
	ctx := context.Background()

	const (
		maxWinners  = 100
		cost        = 1 // Ether
		reserveFor  = 2
		reserveEach = 3
		reserveFree = 4
		entrants    = 50
	)
	sim, runner, raffle, raffleAddr := deploy(t, maxWinners, eth.Ether(cost))

	var totalEntries int64

	const (
		reservedFreeAcc = iota + vandal + 1
		reservedPaidAcc0
		reservedPaidAcc1
	)

	t.Run("free reservation by owner", func(t *testing.T) {
		_, err := raffle.ReserveFree(sim.Acc(vandal), sim.Acc(vandal).From, big.NewInt(1))
		if diff := ethtest.RevertDiff(err, "Ownable: caller is not the owner"); diff != "" {
			t.Errorf("ReserveFree() called by non-owner; %s", diff)
		}
		wantEntrantState(t, raffle, sim.Acc(vandal).From, 0, 0)

		if _, err := raffle.ReserveFree(sim.Acc(deployer), sim.Acc(reservedFreeAcc).From, big.NewInt(reserveFree)); err != nil {
			t.Fatalf("ReserveFree() called by owner; error %v", err)
		}
		wantEntrantState(t, raffle, sim.Acc(reservedFreeAcc).From, reserveFree, reserveFree)
		totalEntries += reserveFree
	})

	for i := 0; i < reserveFor; i++ {
		totalEntries += reserveEach
		acc := reservedPaidAcc0 + i
		if _, err := runner.Reserve(sim.WithValueFrom(acc, eth.Ether(cost*reserveEach)), big.NewInt(reserveEach)); err != nil {
			t.Fatalf("Reserve(%d) error %v", reserveEach, err)
		}
		wantEntrantState(t, raffle, sim.Acc(acc).From, reserveEach, reserveEach)
	}

	for i := 0; i < entrants; i++ {
		entries := int64(i%4 + 1)
		totalEntries += entries

		entrant := sim.WithValueFrom(i, eth.Ether(cost*entries))
		tx, err := raffle.Enter(entrant, entrant.From, big.NewInt(entries))
		if err != nil {
			t.Fatalf("Enter(%d) error %v", entries, err)
		}
		ethtest.GasPrice = 100
		ethtest.LogGas(t, tx, fmt.Sprintf("%d entries", entries))

		wantEntries := entries
		var wantWins int64
		switch i {
		case reservedFreeAcc:
			wantEntries += reserveFree
			wantWins = reserveFree
		case reservedPaidAcc0, reservedPaidAcc1:
			wantEntries += reserveEach
			wantWins = reserveEach
		}
		wantEntrantState(t, raffle, sim.Acc(i).From, wantEntries, wantWins)
	}

	// Entropy needs to be anything non-zero as this is the check for it being
	// set.
	var entropy [32]byte
	entropy[0] = 1

	noEntries := func(t *testing.T, when, wantRevertMsg string) {
		t.Helper()

		t.Run("no entries when "+when, func(t *testing.T) {
			entrant := sim.WithValueFrom(0, eth.Ether(cost))

			_, err := raffle.Enter(entrant, entrant.From, big.NewInt(1))
			if diff := ethtest.RevertDiff(err, wantRevertMsg); diff != "" {
				t.Errorf("Enter() %s", diff)
			}
		})
	}

	if _, err := raffle.Pause(sim.Acc(deployer)); err != nil {
		t.Errorf("Pause() error %v", err)
	}
	noEntries(t, "paused", "Pausable: paused")

	if _, err := raffle.SetEntropy(sim.Acc(deployer), entropy); err != nil {
		t.Fatalf("SetEntropy() error %v", err)
	}

	if _, err := raffle.Unpause(sim.Acc(deployer)); err != nil {
		t.Errorf("Unpause() error %v", err)
	}
	noEntries(t, "entropy set", "Raffle: entropy set")

	t.Run("no reservations when entropy set", func(t *testing.T) {
		const wantRevertMsg = "Raffle: entropy set"
		entrant := sim.WithValueFrom(0, eth.Ether(cost))

		_, err := runner.Reserve(entrant, big.NewInt(1))
		if diff := ethtest.RevertDiff(err, wantRevertMsg); diff != "" {
			t.Errorf("Reserve() %s", diff)
		}

		_, err = raffle.ReserveFree(sim.Acc(deployer), entrant.From, big.NewInt(1))
		if diff := ethtest.RevertDiff(err, wantRevertMsg); diff != "" {
			t.Errorf("ReserveFree() %s", diff)
		}
	})

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
		entrant, err := raffle.Entrants(nil, sim.Acc(i).From)
		if err != nil {
			entrant.Entries = big.NewInt(-1)
			entrant.Wins = big.NewInt(-1)
		}

		tx, err := raffle.Redeem(sim.Acc(i), sim.Acc(i).From)
		if err != nil {
			t.Fatalf("Redeem() error %v", err)
		}
		ethtest.LogGas(t, tx, fmt.Sprintf("Redeem %d entries with %d wins", entrant.Entries, entrant.Wins))
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

	t.Run("no excess", func(t *testing.T) {
		// Thorough testing of excess purchasing is performed elsewhere.
		_, err := raffle.PurchaseExcess(sim.WithValueFrom(0, eth.Ether(cost)), sim.Acc(0).From, big.NewInt(1))
		if diff := ethtest.RevertDiff(err, "Raffle: insufficient excess"); diff != "" {
			t.Errorf("PurchaseExcess(1) when none available; %s", diff)
		}
	})

	t.Run("withdrawal", func(t *testing.T) {
		wantRevenues := eth.Ether(cost * (maxWinners - reserveFree))

		if got := sim.BalanceOf(ctx, t, raffleAddr); got.Cmp(wantRevenues) != 0 {
			t.Errorf("Balance of raffle contract after all redemptions; got %d; want %d", got, wantRevenues)
		}

		beneficiary := sim.MonotonicAddress(t)

		_, err := raffle.WithdrawAll(sim.Acc(deployer+1), beneficiary)
		if diff := ethtest.RevertDiff(err, "Ownable: caller is not the owner"); diff != "" {
			t.Errorf("WithdrawAll() as non contract owner: %s", diff)
		}
		if got := sim.BalanceOf(ctx, t, beneficiary); got.Cmp(big.NewInt(0)) != 0 {
			t.Fatalf("After WithdrawAll() by non contract owner; balance of beneficiary is %d; want 0", got)
		}

		if _, err := raffle.WithdrawAll(sim.Acc(deployer), beneficiary); err != nil {
			t.Errorf("WithdrawAll() as contract owner; error %v", err)
		}
		if got := sim.BalanceOf(ctx, t, beneficiary); got.Cmp(wantRevenues) != 0 {
			t.Errorf("After WithdrawAll() by contract owner; balance of beneficiary is %d; want %d", got, wantRevenues)
		}

		_, err = raffle.Withdraw(sim.Acc(deployer), beneficiary, big.NewInt(1))
		if diff := ethtest.RevertDiff(err, "Raffle: all withdrawn"); diff != "" {
			t.Errorf("Withdraw(1 Wei) after WithdrawAll() %s", diff)
		}
	})
}

func TestReservationLimits(t *testing.T) {
	const (
		maxWinners = 10
		cost       = 1
	)
	sim, runner, raffle, _ := deploy(t, maxWinners, eth.Ether(cost))

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
