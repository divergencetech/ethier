package payments_test

import (
	"math/big"
	"sort"
	"testing"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-cmp/cmp"
	"github.com/h-fam/errdiff"
)

// Default account IDs for various roles.
const (
	owner       = 0
	beneficiary = 9
)

// Default PurchaseManager config.
const (
	totalInventory = 25
	maxPerAddress  = 10
	maxPerTx       = 3
)

func deploy(t *testing.T, auctionConfig LinearDutchAuctionDutchAuctionConfig) (*ethtest.SimulatedBackend, common.Address, *TestableDutchAuction) {
	t.Helper()
	sim := ethtest.NewSimulatedBackendTB(t, 10)

	t.Logf("%T = %+v", auctionConfig, auctionConfig)
	purchaseConfig := PurchaseManagerPurchaseConfig{
		TotalInventory:  big.NewInt(totalInventory),
		MaxPerAddress:   big.NewInt(maxPerAddress),
		MaxPerTx:        big.NewInt(maxPerTx),
		AlsoLimitOrigin: true,
	}
	t.Logf("%T = %+v", purchaseConfig, purchaseConfig)

	addr, _, auction, err := DeployTestableDutchAuction(
		sim.Acc(owner), sim, auctionConfig, purchaseConfig, sim.Acc(beneficiary).From,
	)
	if err != nil {
		t.Fatalf("DeployTestableDutchAuction() error %v", err)
	}

	return sim, addr, auction
}

func deployConstantPrice(t *testing.T, price *big.Int) (*ethtest.SimulatedBackend, common.Address, *TestableDutchAuction) {
	return deploy(t, auctionConfig(1, 1, price, big.NewInt(0)))
}

func auctionConfig(startBlock, endBlock int64, startPrice, perBlockDecrease *big.Int) LinearDutchAuctionDutchAuctionConfig {
	return LinearDutchAuctionDutchAuctionConfig{
		StartBlock:       big.NewInt(startBlock),
		EndBlock:         big.NewInt(endBlock),
		StartPrice:       startPrice,
		PerBlockDecrease: perBlockDecrease,
	}
}

func TestLinearPriceDecrease(t *testing.T) {
	const startBlock = 10

	type wantCost struct {
		block     int64
		num       *big.Int
		totalCost *big.Int
	}

	one := big.NewInt(1)
	two := big.NewInt(2)

	tests := []struct {
		name                         string
		endBlock                     int64
		startPrice, perBlockDecrease *big.Int
		want                         []wantCost
	}{
		{
			name:             "constant",
			endBlock:         startBlock + 10,
			startPrice:       eth.Ether(2),
			perBlockDecrease: big.NewInt(0),
			want: []wantCost{
				{startBlock, one, eth.Ether(2)},
				{startBlock + 1, one, eth.Ether(2)},
				{startBlock + 9, two, eth.Ether(4)},
				{startBlock + 10, one, eth.Ether(2)},
				{startBlock + 100, one, eth.Ether(2)},
			},
		},
		{
			name:             "immediate end",
			endBlock:         startBlock,
			startPrice:       eth.Ether(42),
			perBlockDecrease: eth.Ether(43),
			want: []wantCost{
				{startBlock, one, eth.Ether(42)},
				{startBlock + 1, one, eth.Ether(42)},
			},
		},
		{
			name:             "decreasing quickly",
			endBlock:         startBlock + 10,
			startPrice:       eth.Ether(11),
			perBlockDecrease: eth.Ether(1),
			want: []wantCost{
				{startBlock, one, eth.Ether(11)},
				{startBlock, two, eth.Ether(22)},
				{startBlock + 1, one, eth.Ether(10)},
				{startBlock + 2, one, eth.Ether(9)},
				{startBlock + 3, one, eth.Ether(8)},
				{startBlock + 4, two, eth.Ether(14)},
				{startBlock + 10, one, eth.Ether(1)},
				{startBlock + 11, one, eth.Ether(1)},
				{startBlock + 100, one, eth.Ether(1)},
			},
		},
		{
			name:             "decreasing slowly",
			endBlock:         startBlock + 1000,
			startPrice:       eth.Ether(101),
			perBlockDecrease: eth.EtherFraction(1, 10),
			want: []wantCost{
				{startBlock + 1, one, eth.EtherFraction(1009, 10)},
				{startBlock + 2, one, eth.EtherFraction(1008, 10)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := auctionConfig(startBlock, tt.endBlock, tt.startPrice, tt.perBlockDecrease)
			t.Logf("%T = %+v", cfg, cfg)
			sim, _, auction := deploy(t, cfg)

			sim.FastForward(big.NewInt(startBlock - 1))
			t.Run("before start", func(t *testing.T) {
				_, err := auction.Cost(nil, big.NewInt(1))
				if diff := ethtest.RevertDiff(err, "LinearDutchAuction: Not started"); diff != "" {
					t.Errorf("Cost() before auction start; %s", diff)
				}
			})

			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].block < tt.want[j].block
			})

			for _, w := range tt.want {
				sim.FastForward(big.NewInt(w.block))
				got, err := auction.Cost(nil, w.num)
				if err != nil {
					t.Errorf("Cost(%d) at block %d; error %v", w.num, sim.BlockNumber(), err)
					continue
				}
				if want := w.totalCost; got.Cmp(want) != 0 {
					t.Errorf("Cost(%d) at block %d; got %d want %d", w.num, sim.BlockNumber(), got, want)
				}
			}
		})
	}
}

func TestTxLimit(t *testing.T) {
	tests := []struct {
		name               string
		buy, wantPurchased int64
	}{
		{
			name:          "exact max per tx",
			buy:           maxPerTx,
			wantPurchased: maxPerTx,
		},
		{
			name:          "one more than tx limit",
			buy:           maxPerTx + 1,
			wantPurchased: maxPerTx,
		},
		{
			name:          "10x tx limit",
			buy:           10 * maxPerTx,
			wantPurchased: maxPerTx,
		},
		{
			name:          "below tx limit",
			buy:           maxPerTx - 1,
			wantPurchased: maxPerTx - 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim, _, auction := deployConstantPrice(t, eth.Ether(0))
			acc := sim.Acc(0)

			if _, err := auction.Buy(acc, big.NewInt(tt.buy)); err != nil {
				t.Fatalf("Buy(%d) error %v", tt.buy, err)
			}

			got, err := auction.Purchased(nil, acc.From)
			if err != nil {
				t.Fatalf("Purchased(%s) after Buy() error %v", acc.From, err)
			}
			if got.Cmp(big.NewInt(tt.wantPurchased)) != 0 {
				t.Errorf("Purchased(%s) after Buy(%d) with max per tx = %d; got %d; want %d", acc.From, tt.buy, maxPerTx, got, tt.wantPurchased)
			}
		})
	}
}

func TestAddressLimit(t *testing.T) {
	sim, auctionAddr, auction := deployConstantPrice(t, eth.Ether(0))

	// This test is primarily to demonstrate to readers of this code that config
	// values are as expected.
	t.Run("confirm config values", func(t *testing.T) {
		got, err := auction.PurchaseConfig(nil)
		if err != nil {
			t.Fatalf("PurchaseConfig() error %v", err)
		}

		// TODO: submit a PR so geth/accounts/abi/bind returns named structs.
		want := struct {
			TotalInventory, MaxPerAddress, MaxPerTx *big.Int
			AlsoLimitOrigin                         bool
		}{
			big.NewInt(25), big.NewInt(10), big.NewInt(3), true,
		}

		if diff := cmp.Diff(want, got, ethtest.Comparers()...); diff != "" {
			t.Errorf("PurchaseConfig() diff (-want +got): \n%s", diff)
		}
	})

	// Allow testing for circumvention by buying through a contract.
	_, _, proxy, err := DeployProxyPurchaser(sim.Acc(0), sim, auctionAddr)
	if err != nil {
		t.Fatalf("DeployProxyPurchaser() error %v", err)
	}

	// Tests are deliberately not hermetic to demonstrate that there is no
	// spillover between addresses. Therefore, all errors MUST use t.Fatal.
	tests := []struct {
		purchaseVia interface { // auction or proxy
			Buy(*bind.TransactOpts, *big.Int) (*types.Transaction, error)
		}
		account            int
		buy, wantPurchased int64
		errDiffAgainst     interface{}
	}{
		{auction, 0, 3, 3, nil},
		{auction, 0, 3, 6, nil},
		{auction, 0, 3, 9, nil},
		{auction, 0, 3, 10, nil}, // capped
		{auction, 0, 1, 10, "Sender limit"},
		{proxy, 0, 1, 10, "Origin limit"},
		{auction, 1, 3, 3, nil}, // not capped because different address
		{auction, 1, 3, 6, nil},
		{auction, 1, 3, 9, nil},
		{auction, 1, 3, 10, nil}, // capped again
		{auction, 2, 3, 3, nil},
		{auction, 2, 3, 5, nil}, // capped by total inventory
		{auction, 2, 1, 5, "Sold out"},
	}

	for _, tt := range tests {
		acc := sim.Acc(tt.account)

		_, err := tt.purchaseVia.Buy(acc, big.NewInt(tt.buy))
		if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
			t.Fatalf("Buy(%d) as account %d; %s", tt.buy, tt.account, diff)
		}

		got, err := auction.Purchased(nil, acc.From)
		if err != nil {
			t.Fatalf("Purchased(%s) error %v", acc.From, err)
		}
		if got.Cmp(big.NewInt(tt.wantPurchased)) != 0 {
			t.Errorf("Purchased(%s) got %d; want %d", acc.From, got, tt.wantPurchased)
		}
	}

	t.Run("total purchased", func(t *testing.T) {
		got, err := auction.TotalSupply(nil)
		if err != nil {
			t.Fatalf("TotalSupply() error %v", err)
		}
		if got.Int64() != totalInventory {
			t.Errorf("TotalSupply() got %d; want %d", got, totalInventory)
		}
	})
}
