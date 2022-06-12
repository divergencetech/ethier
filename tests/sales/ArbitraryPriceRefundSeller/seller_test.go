package ArbitraryPriceRefundSeller

import (
	"context"
	"log"
	"math/big"
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

func deploy(t *testing.T, inventory int64) (*ethtest.SimulatedBackend, *TestableArbitraryPriceSeller) {
	t.Helper()
	sim := ethtest.NewSimulatedBackendTB(t, 1)

	_, _, seller, err := DeployTestableArbitraryPriceSeller(sim.Acc(0), sim, big.NewInt(inventory))
	if err != nil {
		t.Fatalf("DeployTestableArbitraryPriceSeller(%d) error %v", inventory, err)
	}
	return sim, seller
}

func TestArbitraryPriceSeller(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		numToPurchase      int64
		costEach, wantPaid *big.Int
	}{
		{
			numToPurchase: 1,
			costEach:      eth.Ether(2),
			wantPaid:      eth.Ether(2),
		},
		{
			numToPurchase: 2,
			costEach:      eth.Ether(1),
			wantPaid:      eth.Ether(2),
		},
		{
			numToPurchase: 3,
			costEach:      eth.Ether(2),
			wantPaid:      eth.Ether(6),
		},
		{
			numToPurchase: 5,
			costEach:      eth.EtherFraction(2, 5),
			wantPaid:      eth.Ether(2),
		},
		{
			numToPurchase: 4,
			costEach:      eth.EtherFraction(256, 1000),
			wantPaid:      eth.EtherFraction(1024, 1000),
		},
	}

	for _, tt := range tests {
		sim, seller := deploy(t, tt.numToPurchase)

		// An address out of our control so its balance can only increase in value.
		beneficiary, err := seller.Sellable(nil)
		if err != nil {
			t.Fatalf("Unable to get sellable address: %v", err)
		}

		if got := sim.BalanceOf(ctx, t, beneficiary); got.Cmp(big.NewInt(0)) != 0 {
			t.Fatalf("Bad test setup; before purchase %T.BalanceOf(beneficiary) got %d; want 0", seller, got)
		}

		acc := sim.WithValueFrom(0, tt.wantPaid)
		sim.Must(t, "Purchase(%d, %d)", tt.numToPurchase, tt.costEach)(seller.Purchase(acc, big.NewInt(tt.numToPurchase), tt.costEach))

		if got := sim.BalanceOf(ctx, t, beneficiary); got.Cmp(tt.wantPaid) != 0 {
			t.Errorf("After purchase of %d items for %d each; BalanceOf(beneficiary) got %d; want %d", tt.numToPurchase, tt.costEach, got, tt.wantPaid)
		}
	}

}

func TestPausing(t *testing.T) {
	sim, seller := deploy(t, 100)

	price := eth.Ether(1)
	acc := sim.WithValueFrom(0, price)

	sim.Must(t, "When not paused, Purchase()")(seller.Purchase(acc, big.NewInt(1), price))
	sim.Must(t, "Pause()")(seller.Pause(sim.Acc(0)))

	if diff := revert.Paused.Diff(seller.Purchase(acc, big.NewInt(1), price)); diff != "" {
		t.Errorf("When paused, Purchase() %s", diff)
	}
}
