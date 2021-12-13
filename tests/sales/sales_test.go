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
	"github.com/divergencetech/ethier/ethtest/revert"
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

func deploy(t *testing.T, cfg config) (*ethtest.SimulatedBackend, common.Address, *TestableDutchAuction) {
	t.Helper()
	sim := ethtest.NewSimulatedBackendTB(t, 10)

	auctionConfig := cfg.convert()

	t.Logf("%T = %+v", auctionConfig, auctionConfig)
	sellerConfig := SellerSellerConfig{
		TotalInventory: big.NewInt(totalInventory),
		MaxPerAddress:  big.NewInt(maxPerAddress),
		MaxPerTx:       big.NewInt(maxPerTx),
		FreeQuota:      big.NewInt(0),
	}
	t.Logf("%T = %+v", sellerConfig, sellerConfig)

	addr, _, auction, err := DeployTestableDutchAuction(
		sim.Acc(0), sim, auctionConfig, cfg.ExpectedReserve, sellerConfig, beneficiary,
	)
	if err != nil {
		t.Fatalf("DeployTestableDutchAuction() error %v", err)
	}

	return sim, addr, auction
}

func deployConstantPrice(t *testing.T, price *big.Int) (*ethtest.SimulatedBackend, common.Address, *TestableDutchAuction) {
	t.Helper()
	return deploy(t, config{
		StartPoint:       1,
		NumDecreases:     0,
		StartPrice:       price,
		DecreaseSize:     big.NewInt(0),
		DecreaseInterval: 1,
		Unit:             Block,
		ExpectedReserve:  price,
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
	// Although not part of the config, the setter function requires an expected
	// reserve as a safety check.
	ExpectedReserve *big.Int
}

func (c config) convert() LinearDutchAuctionDutchAuctionConfig {
	if c.StartPrice == nil {
		c.StartPrice = big.NewInt(0)
	}
	if c.DecreaseSize == nil {
		c.DecreaseSize = big.NewInt(0)
	}
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
				ExpectedReserve:  eth.Ether(2),
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
				ExpectedReserve:  eth.Ether(42),
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
				ExpectedReserve:  eth.Ether(1),
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
				ExpectedReserve:  eth.Ether(1),
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
				ExpectedReserve:  eth.Ether(5),
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
			sim, _, auction := deploy(t, tt.cfg)

			sim.FastForward(big.NewInt(startBlock - 1))
			t.Run("before start", func(t *testing.T) {
				if diff := revert.NotStarted.Diff(auction.Cost(nil, big.NewInt(1))); diff != "" {
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

func TestTimeBasedDecrease(t *testing.T) {
	const (
		decreaseInterval = 60
		numDecreases     = 3
	)
	cfg := config{
		StartPoint:       0, // Will be reset later
		StartPrice:       eth.Ether(1),
		DecreaseSize:     eth.EtherFraction(1, 10),
		NumDecreases:     numDecreases,
		Unit:             Time,
		DecreaseInterval: decreaseInterval,
		ExpectedReserve:  eth.EtherFraction(7, 10),
	}

	sim, _, auction := deploy(t, cfg)

	timestamp := func(t *testing.T) int64 {
		t.Helper()
		tm, err := auction.Timestamp(nil)
		if err != nil {
			t.Errorf("Timestamp() error %v", err)
			return 0
		}
		if !tm.IsInt64() {
			t.Errorf("Timestamp().IsInt64() = false; want true")
		}
		return tm.Int64()
	}

	start := timestamp(t) + 60
	sim.Must(t, "SetAuctionStartPoint(now+60")(auction.SetAuctionStartPoint(sim.Acc(0), big.NewInt(start)))

	// Note: SimulatedBackend has undocumented functionality that adjusts its
	// clock by 10 seconds for every Commit(). The specifics of the step aren't
	// important, but we use Commit() as a means of ticking the clock.
	//
	// Because of this, we lock in the expected behaviour with a test to ensure
	// test coverage across the life of the auction. Each of these will be
	// deleted when encountered.
	unseenCosts := map[uint64]bool{
		eth.Ether(1).Uint64():             true,
		eth.EtherFraction(9, 10).Uint64(): true,
		eth.EtherFraction(8, 10).Uint64(): true,
		eth.EtherFraction(7, 10).Uint64(): true,
	}

	for i := 0; i < 50; i++ {
		sim.Commit() // i.e. clock tick

		// This type of logic-based test setup isn't ideal, but we're forced by
		// the weird time behaviour of the simulator. It attempts to mimic a
		// table-driven test pattern as closely as possible.
		var (
			want           *big.Int
			errDiffAgainst interface{}
		)

		dt := timestamp(t) - start
		t.Run(fmt.Sprintf("time %ds from start", dt), func(t *testing.T) {
			switch {
			case dt < 0:
				errDiffAgainst = "LinearDutchAuction: Not started"
			case dt < decreaseInterval:
				want = cfg.StartPrice
			case dt < 2*decreaseInterval:
				want = new(big.Int).Sub(cfg.StartPrice, cfg.DecreaseSize)
			case dt < 3*decreaseInterval:
				want = new(big.Int).Sub(cfg.StartPrice, new(big.Int).Mul(cfg.DecreaseSize, big.NewInt(2)))
			default:
				want = new(big.Int).Sub(cfg.StartPrice, new(big.Int).Mul(cfg.DecreaseSize, big.NewInt(numDecreases)))
			}

			got, err := auction.Cost(nil, big.NewInt(1))
			if diff := errdiff.Check(err, errDiffAgainst); diff != "" {
				t.Fatalf("Cost() %s", diff)
			}
			if errDiffAgainst != nil {
				return
			}

			if got.Cmp(want) != 0 {
				t.Errorf("Cost(1) got %d; want %d", got, want)
			} else {
				delete(unseenCosts, got.Uint64())
			}
		})
	}

	if len(unseenCosts) > 0 {
		t.Errorf("Incomplete test coverage; unseen costs: %v", unseenCosts)
	}
}

func TestReserveCheck(t *testing.T) {
	sim, _, auction := deployConstantPrice(t, eth.Ether(0))

	tests := []struct {
		name            string
		config          config
		expectedReserve *big.Int
		wantErr         bool
	}{
		{
			name: "zero reserve",
			config: config{
				StartPrice:   eth.Ether(1),
				DecreaseSize: eth.EtherFraction(1, 10),
				NumDecreases: 10,
				Unit:         Block,
			},
			expectedReserve: big.NewInt(0),
			wantErr:         false,
		},
		{
			name: "non-zero reserve",
			config: config{
				StartPrice:   eth.Ether(1),
				DecreaseSize: eth.EtherFraction(1, 10),
				NumDecreases: 2,
				Unit:         Block,
			},
			expectedReserve: eth.EtherFraction(8, 10),
			wantErr:         false,
		},
		{
			name: "non-zero reserve",
			config: config{
				StartPrice:   eth.Ether(200),
				DecreaseSize: eth.EtherFraction(1, 2),
				NumDecreases: 100,
				Unit:         Block,
			},
			expectedReserve: eth.Ether(150),
			wantErr:         false,
		},
		{
			name: "incorrect reserve",
			config: config{
				StartPrice:   eth.Ether(15),
				DecreaseSize: big.NewInt(1),
				NumDecreases: 10,
				Unit:         Block,
			},
			expectedReserve: eth.Ether(6),
			wantErr:         true,
		},
		{
			name: "underflow to 'negative' reserve",
			config: config{
				StartPrice:   eth.Ether(1),
				DecreaseSize: big.NewInt(1),
				NumDecreases: 10,
				Unit:         Block,
			},
			expectedReserve: eth.Ether(-9), // Impossible but underflow will trigger error
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.DecreaseInterval = 1

			var errDiffAgainst interface{}
			if tt.wantErr {
				errDiffAgainst = "LinearDutchAuction: incorrect reserve"
			}

			_, err := auction.SetAuctionConfig(sim.Acc(0), tt.config.convert(), tt.expectedReserve)
			if diff := errdiff.Check(err, errDiffAgainst); diff != "" {
				t.Errorf("SetAuctionConfig(%+v, %d) %s", tt.config, tt.expectedReserve, diff)
			}
		})
	}
}

func TestZeroStartToDisable(t *testing.T) {
	cfg := config{
		StartPoint:       0,
		StartPrice:       eth.Ether(1),
		DecreaseInterval: 1,
		Unit:             Block,
		ExpectedReserve:  eth.Ether(1),
	}
	sim, _, auction := deploy(t, cfg)

	const startBlock = 10
	sim.FastForward(big.NewInt(startBlock))

	if diff := revert.NotStarted.Diff(auction.Cost(nil, big.NewInt(1))); diff != "" {
		t.Errorf("Cost() when StartPoint==0; %s", diff)
	}

	if _, err := auction.SetAuctionStartPoint(sim.Acc(0), big.NewInt(startBlock)); err != nil {
		t.Fatalf("SetAuctionStartPoint() error %v", err)
	}

	got, err := auction.Cost(nil, big.NewInt(1))
	if err != nil {
		t.Fatalf("Cost() when StartPoint!=0; error %v", err)
	}
	if want := cfg.StartPrice; got.Cmp(want) != 0 {
		t.Errorf("Cost() when StartPoint!=0; got %d; want %d", got, want)
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
			TotalInventory, MaxPerAddress, MaxPerTx, FreeQuota  *big.Int
			ReserveFreeQuota, LockFreeQuota, LockTotalInventory bool
		}{
			big.NewInt(25), big.NewInt(10), big.NewInt(3), big.NewInt(0),
			false, false, false,
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

		if got, err := auction.ReceivedFree(nil, recipient.From); got.Cmp(big.NewInt(0)) != 0 || err != nil {
			t.Errorf("ReceivedFree(%s) got %d, err %v; want 0, nil err", recipient.From, got, err)
		}
	}

	t.Run("total sold", func(t *testing.T) {
		got, err := auction.TotalSold(nil)
		if err != nil {
			t.Fatalf("TotalSold() error %v", err)
		}
		if got.Cmp(big.NewInt(totalInventory)) != 0 {
			t.Errorf("TotalSold() got %d; want %d", got, totalInventory)
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

	sim.Must(t, "When not paused, Buy()")(auction.Buy(sim.Acc(0), sim.Acc(0).From, big.NewInt(1)))
	sim.Must(t, "Pause()")(auction.Pause(sim.Acc(0)))

	if diff := revert.Paused.Diff(auction.Buy(sim.Acc(0), sim.Acc(0).From, big.NewInt(1))); diff != "" {
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

	if diff := revert.Reentrant.Diff(
		attacker.Buy(sim.WithValueFrom(0, eth.Ether(10)), sim.Acc(0).From, big.NewInt(1)),
	); diff != "" {
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

func TestUnlimited(t *testing.T) {
	sim, _, auction := deployConstantPrice(t, big.NewInt(1))

	const total = 1e6
	cfg := SellerSellerConfig{
		TotalInventory: big.NewInt(total),
		MaxPerAddress:  big.NewInt(0),
		MaxPerTx:       big.NewInt(0),
		FreeQuota:      big.NewInt(0),
	}
	sim.Must(t, "SetSellerConfig(%+v)", cfg)(auction.SetSellerConfig(sim.Acc(0), cfg))

	buyer := sim.WithValueFrom(0, big.NewInt(total))
	sim.Must(t, "Buy(%d) with unlimited transaction / address limits", int(total))(
		auction.Buy(buyer, buyer.From, big.NewInt(total)),
	)

	if diff := revert.SoldOut.Diff(auction.Buy(buyer, buyer.From, big.NewInt(1))); diff != "" {
		t.Errorf("Buy(1) with no more inventory; %s", diff)
	}
}

// wantOwned is a helper to confirm total number of items owned by an address,
// and total received free of charge.
func wantOwned(t *testing.T, auction *TestableDutchAuction, addr common.Address, wantTotal, wantFree int64) {
	t.Helper()

	if got, err := auction.Own(nil, addr); got.Cmp(big.NewInt(wantTotal)) != 0 || err != nil {
		t.Errorf("Own(%q) got %d, err %v; want %d, nil err", addr, got, err, wantTotal)
	}
	if got, err := auction.ReceivedFree(nil, addr); got.Cmp(big.NewInt(wantFree)) != 0 || err != nil {
		t.Errorf("ReceivedFree(%q) got %d, err %v; want %d, nil err", addr, got, err, wantFree)
	}
}

func TestReservedFreePurchasing(t *testing.T) {
	sim, _, auction := deployConstantPrice(t, big.NewInt(1))

	const (
		deployer = iota
		rcvFree
		buyer
	)

	const (
		totalInventory = 10
		freeQuota      = 4
	)

	cfg := SellerSellerConfig{
		TotalInventory:   big.NewInt(totalInventory),
		FreeQuota:        big.NewInt(freeQuota),
		ReserveFreeQuota: true,
		// Required but irrelevant to tests:
		MaxPerAddress: big.NewInt(0),
		MaxPerTx:      big.NewInt(0),
	}
	sim.Must(t, "SetSellerConfig(%+v", cfg)(auction.SetSellerConfig(sim.Acc(deployer), cfg))

	t.Run("only owner purchases free", func(t *testing.T) {
		if diff := revert.OnlyOwner.Diff(auction.PurchaseFreeOfCharge(sim.Acc(rcvFree), sim.Acc(rcvFree).From, big.NewInt(1))); diff != "" {
			t.Errorf("PurchaseFreeOfCharge() as non-owner; %s", diff)
		}
		wantOwned(t, auction, sim.Acc(rcvFree).From, 0, 0)
	})

	t.Run("reserved free quota honoured", func(t *testing.T) {
		n := big.NewInt(totalInventory - freeQuota)
		sim.Must(t, "Buy(totalInventory - freeQuota)")(auction.Buy(sim.WithValueFrom(buyer, n), sim.Acc(buyer).From, n))

		if diff := revert.SoldOut.Diff(
			auction.Buy(sim.WithValueFrom(buyer, big.NewInt(1)), sim.Acc(buyer).From, big.NewInt(1)),
		); diff != "" {
			t.Errorf("After Buy(totalInventory - freeQuota); Buy(1) %s", diff)
		}

		wantOwned(t, auction, sim.Acc(buyer).From, totalInventory-freeQuota, 0)
	})

	t.Run("free quota enforced", func(t *testing.T) {
		sim.Must(t, "PurchaseFreeOfCharge(freeQuota)")(auction.PurchaseFreeOfCharge(sim.Acc(deployer), sim.Acc(rcvFree).From, big.NewInt(freeQuota)))

		c := revert.Checker("Seller: Free quota exceeded")
		if diff := c.Diff(auction.PurchaseFreeOfCharge(sim.Acc(deployer), sim.Acc(rcvFree).From, big.NewInt(1))); diff != "" {
			t.Errorf("PurchaseFreeOfCharge(1) after exhausting quota; %s", diff)
		}

		wantOwned(t, auction, sim.Acc(rcvFree).From, freeQuota, freeQuota)
	})
}

func TestIssue7Regression(t *testing.T) {
	sim, _, auction := deployConstantPrice(t, big.NewInt(1))

	const (
		deployer = iota
		rcvFree
		buyer
	)

	const (
		totalInventory = 20
		freeQuota      = 3
	)

	cfg := SellerSellerConfig{
		TotalInventory:   big.NewInt(totalInventory),
		FreeQuota:        big.NewInt(freeQuota),
		ReserveFreeQuota: true,
		// Required but irrelevant to tests:
		MaxPerAddress: big.NewInt(0),
		MaxPerTx:      big.NewInt(0),
	}
	sim.Must(t, "SetSellerConfig(%+v", cfg)(auction.SetSellerConfig(sim.Acc(deployer), cfg))
	sim.Must(t, "PurchaseFreeOfCharge(freeQuota)")(auction.PurchaseFreeOfCharge(sim.Acc(deployer), sim.Acc(rcvFree).From, big.NewInt(freeQuota)))
	sim.Must(t, "Buy(totalInventory-2*freeQuota)")(auction.Buy(sim.WithValueFrom(buyer, big.NewInt(totalInventory-2*freeQuota)), sim.Acc(buyer).From, big.NewInt(totalInventory-2*freeQuota)))

	// Without the fix, this call would have reverted due to https://github.com/divergencetech/ethier/issues/7
	sim.Must(t, "Buy(freeQuota)")(auction.Buy(sim.WithValueFrom(buyer, big.NewInt(freeQuota)), sim.Acc(buyer).From, big.NewInt(freeQuota)))

	if diff := revert.SoldOut.Diff(
		auction.Buy(sim.WithValueFrom(buyer, big.NewInt(1)), sim.Acc(buyer).From, big.NewInt(1)),
	); diff != "" {
		t.Errorf("After Buy(totalInventory - freeQuota); Buy(1) %s", diff)
	}
}

func TestUnreservedFreePurchasing(t *testing.T) {
	sim, _, auction := deployConstantPrice(t, big.NewInt(1))

	const (
		totalInventory = 10
		freeQuota      = 4
	)

	cfg := SellerSellerConfig{
		TotalInventory:   big.NewInt(totalInventory),
		FreeQuota:        big.NewInt(freeQuota),
		ReserveFreeQuota: false,
		// Required but irrelevant to tests:
		MaxPerAddress: big.NewInt(0),
		MaxPerTx:      big.NewInt(0),
	}
	sim.Must(t, "SetSellerConfig(%+v)", cfg)(auction.SetSellerConfig(sim.Acc(0), cfg))

	t.Run("unreserved free quota ignored", func(t *testing.T) {
		n := big.NewInt(totalInventory - freeQuota)
		sim.Must(t, "Buy(totalInventory - freeQuota)")(auction.Buy(sim.WithValueFrom(0, n), sim.Acc(0).From, n))
		sim.Must(t, "After Buy(totalInventory - freeQuota); Buy(1)")(
			auction.Buy(sim.WithValueFrom(0, big.NewInt(1)), sim.Acc(0).From, big.NewInt(1)),
		)

		wantOwned(t, auction, sim.Acc(0).From, totalInventory-freeQuota+1, 0)
	})

	t.Run("total inventory honoured even if free", func(t *testing.T) {
		sim.Must(t, "PurchsaeFreeOfCharge(freeQuota-1)")(
			auction.PurchaseFreeOfCharge(sim.Acc(0), sim.Acc(0).From, big.NewInt(freeQuota-1)),
		)

		if diff := revert.SoldOut.Diff(
			auction.PurchaseFreeOfCharge(sim.Acc(0), sim.Acc(0).From, big.NewInt(1)),
		); diff != "" {
			t.Errorf("PurchaseFreeOfCharge(1) when last unreserved quota already sold; %s", diff)
		}

		if diff := revert.SoldOut.Diff(
			auction.Buy(sim.Acc(0), sim.Acc(0).From, big.NewInt(1)),
		); diff != "" {
			t.Errorf("Buy(1) when last unreserved quota already sold; %s", diff)
		}

		wantOwned(t, auction, sim.Acc(0).From, totalInventory, freeQuota-1)
	})
}
