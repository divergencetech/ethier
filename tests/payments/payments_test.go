package payments_test

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
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

// Default PurchaseManager config.
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
		sim.Acc(0), sim, auctionConfig, purchaseConfig, beneficiary,
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
		if got.Cmp(big.NewInt(totalInventory)) != 0 {
			t.Errorf("TotalSupply() got %d; want %d", got, totalInventory)
		}
	})
}

func TestFundsManagement(t *testing.T) {
	ctx := context.Background()
	price := eth.Ether(1)
	sim, auctionAddr, auction := deployConstantPrice(t, price)

	if got := sim.BalanceOf(ctx, t, beneficiary); got.Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("Bad test setup; beneficiary has non-zero balance; %T.BalanceOf(%s) got %d; want 0", sim, beneficiary, got)
	}

	refunds := make(chan *TestableDutchAuctionRefund)
	auction.TestableDutchAuctionFilterer.WatchRefund(&bind.WatchOpts{}, refunds, nil)
	defer close(refunds)

	revenues := make(chan *TestableDutchAuctionRevenue)
	auction.TestableDutchAuctionFilterer.WatchRevenue(&bind.WatchOpts{}, revenues, nil)
	defer close(revenues)

	tests := []struct {
		account                                             int
		num                                                 int64
		sendValue, wantSpent, wantRefund, wantTotalRevenues *big.Int
		errDiffAgainst                                      interface{}
	}{
		{
			account:           0,
			num:               3,
			sendValue:         eth.Ether(5),
			wantSpent:         eth.Ether(3),
			wantRefund:        eth.Ether(2),
			wantTotalRevenues: eth.Ether(3),
		},
		{
			account:           0,
			num:               1,
			sendValue:         new(big.Int).Sub(price, big.NewInt(1)),
			errDiffAgainst:    "Costs 1000000000 GWei",
			wantTotalRevenues: eth.Ether(3),
		},
		{
			account:           1,
			num:               2,
			sendValue:         eth.Ether(2),
			wantSpent:         eth.Ether(2),
			wantTotalRevenues: eth.Ether(5),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("purchase[%d]", i), func(t *testing.T) {
			before := sim.BalanceOf(ctx, t, sim.Acc(tt.account).From)

			tx, err := auction.Buy(sim.WithValueFrom(tt.account, tt.sendValue), big.NewInt(tt.num))
			if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
				t.Fatalf("Buy() %s", diff)
			}

			t.Run("revenues", func(t *testing.T) {
				if got := sim.BalanceOf(ctx, t, beneficiary); got.Cmp(tt.wantTotalRevenues) != 0 {
					t.Errorf("Calculating total revenues; %T.BalanceOf(<beneficiary>) got %d; want %d", sim, got, tt.wantTotalRevenues)
				}

				select {
				case got := <-revenues:
					got.Raw = types.Log{}
					want := &TestableDutchAuctionRevenue{
						Beneficiary:  beneficiary,
						NumPurchased: big.NewInt(tt.num),
						Amount:       tt.wantSpent,
					}
					if diff := cmp.Diff(want, got, ethtest.Comparers()...); diff != "" {
						t.Errorf("Revenue event diff (-want +got):\n%s", diff)
					}
				default:
					if tt.wantSpent != nil && tt.wantSpent.Cmp(big.NewInt(0)) != 0 {
						t.Errorf("No Revenue event; want amount %d", tt.wantSpent)
					}
				}
			})

			t.Run("refund", func(t *testing.T) {
				select {
				case got := <-refunds:
					got.Raw = types.Log{}
					want := &TestableDutchAuctionRefund{
						Buyer:  sim.Acc(tt.account).From,
						Amount: tt.wantRefund,
					}
					if diff := cmp.Diff(want, got, ethtest.Comparers()...); diff != "" {
						t.Errorf("Refund event diff (-want +got):\n%s", diff)
					}
				default:
					if tt.wantRefund != nil && tt.wantRefund.Cmp(big.NewInt(0)) != 0 {
						t.Errorf("No Refund event logged; want refund of %d", tt.wantRefund)
					}
				}
			})

			t.Run("complete disbursement of funds", func(t *testing.T) {
				if got := sim.BalanceOf(ctx, t, auctionAddr); got.Cmp(big.NewInt(0)) != 0 {
					t.Errorf("Contract kept funds; %T.BalanceOf(<PaymentManger>) got %d; want 0", sim, got)
				}
			})

			if err != nil {
				return
			}

			t.Run("buyer balance decrease", func(t *testing.T) {
				// TODO: these tests fail because of a small (~0.0001ETH)
				// discrepancy in the gas calculation, which needs investigation.
				t.Skip("TODO: Investigate gas cost discrepancy")

				after := sim.BalanceOf(ctx, t, sim.Acc(tt.account).From)
				spent := new(big.Int).Sub(before, after)
				spent.Sub(spent, sim.GasSpent(ctx, t, tx))
				if got := spent; got.Cmp(tt.wantSpent) != 0 {
					t.Errorf("Buy(%d) at price %d; got balance reduction of %d (excluding gas); want %d", tt.num, price, got, tt.wantSpent)
				}
			})
		})
	}
}
