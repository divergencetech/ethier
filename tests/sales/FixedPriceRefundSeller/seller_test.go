package FixedPriceRefundSeller

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
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

func TestFixedPriceSeller(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, 1)

	tests := []struct {
		price     *big.Int
		wantCosts map[int64]*big.Int
	}{
		{
			price: eth.Ether(1),
			wantCosts: map[int64]*big.Int{
				0: big.NewInt(0),
				1: eth.Ether(1),
				2: eth.Ether(2),
				7: eth.Ether(7),
			},
		},
		{
			price: eth.EtherFraction(1, 2),
			wantCosts: map[int64]*big.Int{
				0:  big.NewInt(0),
				1:  eth.EtherFraction(1, 2),
				2:  eth.Ether(1),
				7:  eth.EtherFraction(7, 2),
				42: eth.Ether(21),
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d wei", tt.price), func(t *testing.T) {
			_, _, seller, err := DeployTestableFixedPriceSeller(sim.Acc(0), sim, tt.price)
			if err != nil {
				t.Fatalf("DeployTestableFixedPriceSeller() error %v", err)
			}

			for n, want := range tt.wantCosts {
				got, err := seller.Cost(nil, uint64(n))
				if err != nil {
					t.Errorf("Cost(%d) error %v", n, err)
					continue
				}
				if got.Cmp(want) != 0 {
					t.Errorf("Cost(%d) got %d; want %d", n, got, want)
				}
			}
		})
	}
}
