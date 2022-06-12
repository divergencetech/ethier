package FreeOwnerAirdropper

import (
	"log"
	"os"
	"testing"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/revert"
	"github.com/ethereum/go-ethereum/common"
)

// Default Seller config.
const (
	totalInventory = 25
	maxPerAddress  = 10
	maxPerTx       = 3
)

var beneficiary common.Address

func TestMain(m *testing.M) {
	// Create a random beneficiary address that we can't control because we
	// throw away the key. That way it can only increase in value.
	s, err := eth.NewSigner(128)
	if err != nil {
		log.Fatalf("eth.NewSigner(128) error %v", err)
	}
	beneficiary = s.Address()

	os.Exit(m.Run())
}

type SellerConfig struct {
	TotalInventory, MaxPerAddress, MaxPerTx uint64
}

func deploy(t *testing.T, totalInventory uint64) (*ethtest.SimulatedBackend, common.Address, *TestableFreeOwnerAirdropper, *SellableMock) {
	t.Helper()
	sim := ethtest.NewSimulatedBackendTB(t, 10)

	addr, _, seller, err := DeployTestableFreeOwnerAirdropper(
		sim.Acc(0), sim, totalInventory,
	)
	if err != nil {
		t.Fatalf("TestableFreeOwnerAirdropper() error %v", err)
	}

	s, err := seller.Sellable(nil)
	if err != nil {
		t.Fatalf("Unable to get sellable: %v", err)
	}
	sellable, err := NewSellableMock(s, sim)
	if err != nil {
		t.Fatalf("Unable to create sellable: %v", err)
	}

	return sim, addr, seller, sellable
}

// wantOwned is a helper to confirm total number of items owned by an address,
// and total received free of charge.
func wantOwned(t *testing.T, s *SellableMock, addr common.Address, wantTotal uint64) {
	t.Helper()

	if got, err := s.BalanceOf(nil, addr); got != wantTotal || err != nil {
		t.Errorf("Own(%q) got %d, err %v; want %d, nil err", addr, got, err, wantTotal)
	}
}

func TestAirdrop(t *testing.T) {
	const (
		deployer = iota
		nobody
		rcvFree1
		rcvFree2
	)

	const (
		totalInventory = 20
	)

	sim, _, seller, sellable := deploy(t, totalInventory)

	singleRcv := []FreeOwnerAirdropperReceiver{
		{
			To:  sim.Addr(rcvFree1),
			Num: 1,
		},
	}

	t.Run("only owner purchases free", func(t *testing.T) {
		if diff := revert.OnlyOwner.Diff(seller.Airdrop(sim.Acc(nobody), singleRcv)); diff != "" {
			t.Errorf("Airdrop() as non-owner; %s", diff)
		}
		wantOwned(t, sellable, sim.Addr(deployer), 0)
		wantOwned(t, sellable, sim.Addr(nobody), 0)
		wantOwned(t, sellable, sim.Addr(rcvFree1), 0)
		wantOwned(t, sellable, sim.Addr(rcvFree2), 0)
	})

	t.Run("airdrop to the correct addresses", func(t *testing.T) {
		rcvs := []FreeOwnerAirdropperReceiver{
			{
				To:  sim.Addr(rcvFree1),
				Num: 15,
			},
			{
				To:  sim.Addr(rcvFree2),
				Num: 5,
			},
		}

		sim.Must(t, "Airdrop(totalInventory)")(seller.Airdrop(sim.Acc(deployer), rcvs))

		wantOwned(t, sellable, sim.Addr(deployer), 0)
		wantOwned(t, sellable, sim.Addr(nobody), 0)
		wantOwned(t, sellable, sim.Addr(rcvFree1), 15)
		wantOwned(t, sellable, sim.Addr(rcvFree2), 5)
	})

	t.Run("limit enforced", func(t *testing.T) {
		if diff := revert.Checker("FixedSupply: To many requested").Diff(
			seller.Airdrop(sim.Acc(deployer), singleRcv),
		); diff != "" {
			t.Errorf("After Airdrop(totalInventory); Airdrop(1) %s", diff)
		}
	})
}
