package sales

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-cmp/cmp"
	"github.com/h-fam/errdiff"
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

func deploy(t *testing.T, auctionConfig LinearDutchAuctionDutchAuctionConfig) (*ethtest.SimulatedBackend, common.Address, *TestableDutchAuction) {
	t.Helper()
	sim := ethtest.NewSimulatedBackendTB(t, 10)

	t.Logf("%T = %+v", auctionConfig, auctionConfig)
	sellerConfig := SellerSellerConfig{
		TotalInventory: big.NewInt(totalInventory),
		MaxPerAddress:  big.NewInt(maxPerAddress),
		MaxPerTx:       big.NewInt(maxPerTx),
	}
	t.Logf("%T = %+v", sellerConfig, sellerConfig)

	addr, _, auction, err := DeployTestableDutchAuction(
		sim.Acc(0), sim, auctionConfig, sellerConfig, beneficiary,
	)
	if err != nil {
		t.Fatalf("DeployTestableDutchAuction() error %v", err)
	}

	return sim, addr, auction
}

func deployConstantPrice(t *testing.T, price *big.Int) (*ethtest.SimulatedBackend, common.Address, *TestableDutchAuction) {
	return deploy(t, LinearDutchAuctionDutchAuctionConfig{
		StartPoint:       big.NewInt(1),
		NumDecreases:     big.NewInt(0),
		StartPrice:       price,
		DecreaseSize:     big.NewInt(0),
		DecreaseInterval: big.NewInt(1),
		Unit:             uint8(Block),
	})
}

// A unit represents the AuctionIntervalUnit enum.
type unit uint8

const (
	UnspecifiedAuctionUnit unit = iota
	Block
	Time
)

func (u unit) String() string {
	names := map[unit]string{
		Block: "Block",
		Time:  "Time",
	}
	n, ok := names[u]
	if !ok {
		return fmt.Sprintf("[UNSPECIFIED UNIT %d]", u)
	}
	return fmt.Sprintf("[UNIT: %s]", n)
}

// A config is equivalent to a LinearDutchAUctionDutchAuctionConfig but uses
// int64 values instead of *big.Int to make tests easier to write. Fields that
// refer to values remain as *big.Int to allow use of the eth.Ether*()
// functions.
type config struct {
	StartPoint, NumDecreases, DecreaseInterval int64
	StartPrice, DecreaseSize                   *big.Int
	Unit                                       unit
}

func (c config) convert() LinearDutchAuctionDutchAuctionConfig {
	return LinearDutchAuctionDutchAuctionConfig{
		StartPoint:       big.NewInt(c.StartPoint),
		NumDecreases:     big.NewInt(c.NumDecreases),
		StartPrice:       c.StartPrice,
		DecreaseSize:     c.DecreaseSize,
		DecreaseInterval: big.NewInt(c.DecreaseInterval),
		Unit:             uint8(c.Unit),
	}
}

func TestLinearPriceDecrease(t *testing.T) {
	const startBlock = 10

	type wantCost struct {
		point     int64 // block or time
		num       *big.Int
		totalCost *big.Int
	}

	one := big.NewInt(1)
	two := big.NewInt(2)

	tests := []struct {
		name string
		cfg  config
		want []wantCost
	}{
		{
			name: "constant",
			cfg: config{
				StartPoint:       startBlock,
				StartPrice:       eth.Ether(2),
				NumDecreases:     10,
				DecreaseInterval: 1,
				DecreaseSize:     big.NewInt(0),
				Unit:             Block,
			},
			want: []wantCost{
				{startBlock, one, eth.Ether(2)},
				{startBlock + 1, one, eth.Ether(2)},
				{startBlock + 9, two, eth.Ether(4)},
				{startBlock + 10, one, eth.Ether(2)},
				{startBlock + 20, one, eth.Ether(2)},
			},
		},
		{
			name: "immediate end",
			cfg: config{
				StartPoint:       startBlock,
				DecreaseInterval: 1,
				StartPrice:       eth.Ether(42),
				NumDecreases:     0,
				DecreaseSize:     eth.Ether(43),
				Unit:             Block,
			},
			want: []wantCost{
				{startBlock, one, eth.Ether(42)},
				{startBlock + 1, one, eth.Ether(42)},
			},
		},
		{
			name: "decreasing quickly",
			cfg: config{
				StartPoint:       startBlock,
				NumDecreases:     10,
				DecreaseInterval: 1,
				StartPrice:       eth.Ether(11),
				DecreaseSize:     eth.Ether(1),
				Unit:             Block,
			},
			want: []wantCost{
				{startBlock, one, eth.Ether(11)},
				{startBlock, two, eth.Ether(22)},
				{startBlock + 1, one, eth.Ether(10)},
				{startBlock + 2, one, eth.Ether(9)},
				{startBlock + 3, one, eth.Ether(8)},
				{startBlock + 4, two, eth.Ether(14)},
				{startBlock + 10, one, eth.Ether(1)},
				{startBlock + 11, one, eth.Ether(1)},
				{startBlock + 20, one, eth.Ether(1)},
			},
		},
		{
			name: "decreasing slowly",
			cfg: config{
				StartPoint:       startBlock,
				NumDecreases:     1000,
				DecreaseInterval: 1,
				StartPrice:       eth.Ether(101),
				DecreaseSize:     eth.EtherFraction(1, 10),
				Unit:             Block,
			},
			want: []wantCost{
				{startBlock + 1, one, eth.EtherFraction(1009, 10)},
				{startBlock + 2, one, eth.EtherFraction(1008, 10)},
			},
		},
		{
			name: "spread decrease with higher interval",
			cfg: config{
				StartPoint:       startBlock,
				NumDecreases:     5,
				DecreaseInterval: 7,
				StartPrice:       eth.Ether(10),
				DecreaseSize:     eth.Ether(1),
				Unit:             Block,
			},
			want: []wantCost{
				// Make sure to test boundaries before and after multiples of
				// decreaseInterval.
				{startBlock, one, eth.Ether(10)},
				//
				{startBlock + 6, one, eth.Ether(10)},
				{startBlock + 7, one, eth.Ether(9)},
				{startBlock + 8, one, eth.Ether(9)},
				//
				{startBlock + 13, one, eth.Ether(9)},
				{startBlock + 14, one, eth.Ether(8)},
				{startBlock + 15, one, eth.Ether(8)},
				//
				{startBlock + 34, one, eth.Ether(6)},
				{startBlock + 35, one, eth.Ether(5)},
				{startBlock + 36, one, eth.Ether(5)},
				// Respects numDecreases
				{startBlock + 43, one, eth.Ether(5)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg.convert()
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
				return tt.want[i].point < tt.want[j].point
			})

			for _, w := range tt.want {
				sim.FastForward(big.NewInt(w.point))
				got, err := auction.Cost(nil, w.num)
				if err != nil {
					t.Errorf("Cost(%d) at block %d; error %v", w.num, sim.BlockNumber(), err)
					continue
				}
				if want := w.totalCost; got.Cmp(want) != 0 {
					t.Errorf("Cost(%d) at point %d; got %d want %d", w.num, sim.BlockNumber(), got, want)
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

			if _, err := auction.Buy(acc, acc.From, big.NewInt(tt.buy)); err != nil {
				t.Fatalf("Buy(%d) error %v", tt.buy, err)
			}

			got, err := auction.Own(nil, acc.From)
			if err != nil {
				t.Fatalf("Own(%s) after Buy() error %v", acc.From, err)
			}
			if got.Cmp(big.NewInt(tt.wantPurchased)) != 0 {
				t.Errorf("Own(%s) after Buy(%d) with max per tx = %d; got %d; want %d", acc.From, tt.buy, maxPerTx, got, tt.wantPurchased)
			}
		})
	}
}

func TestAddressLimit(t *testing.T) {
	sim, auctionAddr, auction := deployConstantPrice(t, eth.Ether(0))

	// This test is primarily to demonstrate to readers of this code that config
	// values are as expected.
	t.Run("confirm config values", func(t *testing.T) {
		got, err := auction.SellerConfig(nil)
		if err != nil {
			t.Fatalf("SellerConfig() error %v", err)
		}

		// TODO: submit a PR so geth/accounts/abi/bind returns named structs.
		want := struct {
			TotalInventory, MaxPerAddress, MaxPerTx *big.Int
		}{
			big.NewInt(25), big.NewInt(10), big.NewInt(3),
		}

		if diff := cmp.Diff(want, got, ethtest.Comparers()...); diff != "" {
			t.Errorf("SellerConfig() diff (-want +got): \n%s", diff)
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
			Buy(*bind.TransactOpts, common.Address, *big.Int) (*types.Transaction, error)
		}
		payer, recipient   int
		buy, wantPurchased int64
		errDiffAgainst     interface{}
	}{
		{
			purchaseVia:    auction,
			payer:          0,
			recipient:      0,
			buy:            3,
			wantPurchased:  3,
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          0,
			recipient:      0,
			buy:            3,
			wantPurchased:  3,
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          0,
			recipient:      0,
			buy:            3,
			wantPurchased:  3,
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          0,
			recipient:      0,
			buy:            3,
			wantPurchased:  1, // capped by address limit
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          0,
			recipient:      0,
			buy:            1,
			wantPurchased:  0,
			errDiffAgainst: "Buyer limit",
		},
		{
			purchaseVia:    auction,
			payer:          0,
			recipient:      1, // can't buy for someone else either
			buy:            1,
			wantPurchased:  0,
			errDiffAgainst: "Sender limit",
		},
		{
			purchaseVia:    auction,
			payer:          1,
			recipient:      0, // can't be bought for by someone else
			buy:            1,
			wantPurchased:  0,
			errDiffAgainst: "Buyer limit",
		},
		{
			purchaseVia:    proxy,
			payer:          0,
			recipient:      1,
			buy:            1,
			wantPurchased:  0,
			errDiffAgainst: "Origin limit",
		},
		{
			purchaseVia:    auction,
			payer:          1,
			recipient:      1,
			buy:            3,
			wantPurchased:  3, // not capped because different address
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          1,
			recipient:      1,
			buy:            3,
			wantPurchased:  3,
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          1,
			recipient:      1,
			buy:            3,
			wantPurchased:  3,
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          1,
			recipient:      1,
			buy:            3,
			wantPurchased:  1, // capped again by wallet limit
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          2,
			recipient:      2,
			buy:            3,
			wantPurchased:  3,
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          2,
			recipient:      2,
			buy:            3,
			wantPurchased:  2, // capped by total inventory
			errDiffAgainst: nil,
		},
		{
			purchaseVia:    auction,
			payer:          2,
			recipient:      2,
			buy:            1,
			wantPurchased:  0,
			errDiffAgainst: "Sold out",
		},
	}

	for _, tt := range tests {
		payer := sim.Acc(tt.payer)
		recipient := sim.Acc(tt.recipient)

		before, err := auction.Own(nil, recipient.From)
		if err != nil {
			t.Fatalf("Own(<recipient>) before purchase; error %v", err)
		}

		_, err = tt.purchaseVia.Buy(payer, recipient.From, big.NewInt(tt.buy))
		if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
			t.Fatalf("Buy(account[%d], n=%d) as account %d; %s", tt.recipient, tt.buy, tt.payer, diff)
		}

		after, err := auction.Own(nil, recipient.From)
		if err != nil {
			t.Fatalf("Own(<recipient>) after purchase attempt; error %v", err)
		}
		if got := after.Sub(after, before); got.Cmp(big.NewInt(tt.wantPurchased)) != 0 {
			t.Errorf("Own(%s) got %d; want %d", recipient.From, got, tt.wantPurchased)
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

			tx, err := auction.Buy(
				sim.WithValueFrom(tt.account, tt.sendValue),
				sim.Acc(tt.account).From,
				big.NewInt(tt.num),
			)
			if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
				t.Fatalf("Buy() %s", diff)
			}

			t.Run("revenues", func(t *testing.T) {
				t.Parallel()
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
				case <-time.After(250 * time.Millisecond):
					// TODO: using a time delay isn't ideal; is there a way to
					// properly synchronise with the event filter?
					if tt.wantSpent != nil && tt.wantSpent.Cmp(big.NewInt(0)) != 0 {
						t.Errorf("No Revenue event; want amount %d", tt.wantSpent)
					}
				}
			})

			t.Run("refund", func(t *testing.T) {
				t.Parallel()
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
				case <-time.After(250 * time.Millisecond):
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
				after := sim.BalanceOf(ctx, t, sim.Acc(tt.account).From)
				gotSpent := new(big.Int).Sub(before, after)
				gotSpent.Sub(gotSpent, sim.GasSpent(ctx, t, tx))

				// TODO: there's a small (~0.0002ETH) discrepancy in the gas
				// calculation, which needs investigation. Although not idea,
				// it's ok to use a tolerance here because the value is in the
				// favour of the buyer.
				tolerance := eth.EtherFraction(1, 5000)
				diff := new(big.Int).Sub(tt.wantSpent, gotSpent)

				if diff.Cmp(tolerance) != -1 {
					t.Errorf("Buy(%d) at price %d; got balance reduction of %d (excluding gas); want %d within tolerance of %d", tt.num, price, gotSpent, tt.wantSpent, tolerance)
				}
			})
		})
	}
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
				got, err := seller.Cost(nil, big.NewInt(n))
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

func TestPausing(t *testing.T) {
	sim, _, auction := deployConstantPrice(t, eth.Ether(0))

	if _, err := auction.Buy(sim.Acc(0), sim.Acc(0).From, big.NewInt(1)); err != nil {
		t.Fatalf("When not paused, Buy() error %v", err)
	}
	if _, err := auction.Pause(sim.Acc(0)); err != nil {
		t.Fatalf("Pause() error %v", err)
	}

	_, err := auction.Buy(sim.Acc(0), sim.Acc(0).From, big.NewInt(1))
	if diff := ethtest.RevertDiff(err, "Pausable: paused"); diff != "" {
		t.Errorf("When paused, Buy() %s", diff)
	}
}

func TestReentrancyGuard(t *testing.T) {
	sim, auctionAddr, auction := deployConstantPrice(t, eth.Ether(1))

	// Contract that reenters upon receiving a refund.
	_, _, attacker, err := DeployReentrantProxyPurchaser(sim.Acc(0), sim, auctionAddr)
	if err != nil {
		t.Fatalf("DeployReentrantProxyPurchaser() error %v", err)
	}

	_, err = attacker.Buy(sim.WithValueFrom(0, eth.Ether(10)), sim.Acc(0).From, big.NewInt(1))
	if diff := errdiff.Check(err, "ReentrancyGuard: reentrant call"); diff != "" {
		t.Errorf("%T.Buy(); invoking Seller._purchase() through reentrant call; %s", attacker, diff)
	}

	got, err := auction.Own(nil, sim.Acc(0).From)
	if err != nil {
		t.Fatalf("Own() error %v", err)
	}
	if got.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Own(<reentrant attacker>) got %d; want 0", got)
	}
}
