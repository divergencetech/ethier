package TestableCommonSellable

import (
	"fmt"
	"testing"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/revert"
	"github.com/ethereum/go-ethereum/common"
)

const (
	deployer int = iota
	nobody
)

func deploy(t *testing.T) (*ethtest.SimulatedBackend, common.Address, *TestableCommonSellable) {
	t.Helper()
	sim := ethtest.NewSimulatedBackendTB(t, 10)

	addr, _, sellable, err := DeployTestableCommonSellable(
		sim.Acc(0), sim,
	)
	if err != nil {
		t.Fatalf("DeployTestableDutchAuction() error %v", err)
	}

	return sim, addr, sellable
}

func deploySeller(t *testing.T, sim *ethtest.SimulatedBackend, a common.Address) (common.Address, *SellerMock) {
	t.Helper()
	addr, _, seller, err := DeploySellerMock(
		sim.Acc(0), sim, a,
	)
	if err != nil {
		t.Fatalf("DeployTestableDutchAuction() error %v", err)
	}

	return addr, seller
}

func TestOnlySellerCanCall(t *testing.T) {
	sim, sAddr, s := deploy(t)
	badSeller, _ := deploySeller(t, sim, sAddr)
	sellers, err := s.GetSellers(nil)

	if err != nil {
		t.Errorf("Unable to get approved sellers: %v", err)
	}

	tests := []struct {
		buyer       int
		addr        common.Address
		numPurchase uint64
		wantSell    bool
		addSellers  []common.Address
		rmSellers   []common.Address
	}{
		{
			buyer:       0,
			addr:        sellers[0],
			numPurchase: 3,
			wantSell:    true,
		},
		{
			buyer:       1,
			addr:        sellers[1],
			numPurchase: 7,
			wantSell:    true,
		},
		{
			buyer:       3,
			addr:        badSeller,
			numPurchase: 12,
			wantSell:    false,
		},
		{
			buyer:       3,
			addr:        badSeller,
			wantSell:    true,
			numPurchase: 12,
			addSellers:  []common.Address{badSeller},
		},
		{
			buyer:       4,
			addr:        sellers[1],
			wantSell:    false,
			numPurchase: 13,
			rmSellers:   sellers,
		},
		{
			buyer:       5,
			addr:        sellers[0],
			wantSell:    false,
			numPurchase: 15,
		},
	}

	totalSupply := uint64(0)

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Purchasing via seller %d", i), func(t *testing.T) {
			sm, err := NewSellerMock(tt.addr, sim)
			if err != nil {
				t.Errorf("Unable connect to seller mock: %v", err)
			}

			if diff := revert.OnlyOwner.Diff(s.ChangeSellers(sim.Acc(nobody), tt.rmSellers, tt.addSellers)); diff != "" {
				t.Errorf("ChangeSigners(%v, %v) from non-owner account: %s", tt.rmSellers, tt.addSellers, diff)
			}

			_, err = s.ChangeSellers(sim.Acc(deployer), tt.rmSellers, tt.addSellers)
			if err != nil {
				t.Errorf("ChangeSigners(%v, %v): %v", tt.rmSellers, tt.addSellers, err)
			}

			if !tt.wantSell {
				if diff := revert.Checker("Unauthorized seller").Diff(
					sm.Purchase(
						sim.WithValueFrom(tt.buyer, eth.Ether(int64(tt.numPurchase))),
						sim.Addr(tt.buyer),
						tt.numPurchase,
					),
				); diff != "" {
					t.Errorf("Purchasing from unapproved seller: %s", diff)
				}

				tt.numPurchase = 0
			} else {
				sim.Must(t, "Purchase from approved seller")(sm.Purchase(
					sim.WithValueFrom(tt.buyer, eth.Ether(int64(tt.numPurchase))),
					sim.Addr(tt.buyer),
					tt.numPurchase,
				))
			}

			totalSupply += tt.numPurchase

			nBought, err := s.Bought(nil, sim.Addr(tt.buyer))
			if err != nil {
				t.Errorf("Unable to fetch bought: %v", err)
			}
			if nBought != tt.numPurchase {
				t.Errorf("Mismatching num bought: want %d, got %d", tt.numPurchase, nBought)
			}

			sup, err := s.TotalSupply(nil)
			if err != nil {
				t.Errorf("Unable to fetch totalSupply: %v", err)
			}
			if sup != totalSupply {
				t.Errorf("Mismatching totalSupply: want %d, got %d", totalSupply, sup)
			}

		})
	}
}
