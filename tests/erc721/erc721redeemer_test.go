package erc721

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/h-fam/errdiff"
)

// See erc721_test.go for const values associated with actors in the tests;
// their actual values are irrelevant, but their names are informative.

const (
	mintPerTransaction = 20
	// All tokens are allowed to redeem claims, but this one has an ERC721
	// approved alternative address.
	approvedRedeemableTokenID = 0
)

func deployRedeemer(t *testing.T, mintingTransactions int) (*ethtest.SimulatedBackend, *TestableERC721Redeemer, *ERC721RedeemerFilterer) {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, numAccounts)

	addr, _, redeem, err := DeployTestableERC721Redeemer(sim.Acc(deployer), sim)
	if err != nil {
		t.Fatalf("DeployTestableERC721Redeemer() error %v", err)
	}

	filter, err := NewERC721RedeemerFilterer(addr, sim)
	if err != nil {
		t.Fatalf("NewERC721RedeemerFilterer() error %v", err)
	}

	tokenAddr, err := redeem.Token(nil)
	if err != nil {
		t.Fatalf("%T.Token() error %v", redeem, err)
	}
	t.Logf("Address %v is ERC721 test token", tokenAddr)
	token, err := NewIERC721(tokenAddr, sim)
	if err != nil {
		t.Fatalf("NewIERC721(%T.Token()) error %v", redeem, tokenAddr)
	}

	for i := 0; i < mintingTransactions; i++ {
		sim.Must(t, "Mint(tokenOwner, %d)", mintPerTransaction)(
			redeem.Mint(
				sim.Acc(deployer),
				sim.Addr(tokenOwner),
				big.NewInt(mintPerTransaction),
			),
		)
	}
	sim.Must(t, "%T.Approve(%d)", token, approvedRedeemableTokenID)(
		token.Approve(
			sim.Acc(tokenOwner),
			sim.Addr(approved),
			big.NewInt(approvedRedeemableTokenID),
		),
	)

	totalSupply := int64(mintPerTransaction * mintingTransactions)
	wantNumWords := totalSupply / 256
	if totalSupply%256 != 0 {
		wantNumWords++
	}

	return sim, redeem, filter
}

type claim struct {
	asAccountID       int
	tokenIDs          []int64
	errDiffAgainst    interface{}
	wantRedeemedAfter map[int]int64
	wantEvents        []*ERC721RedeemerRedemption
	wantClaimedAfter  map[int64]int64
}

// maxOneRedemptionSteps returns test steps common to regular testing with a
// maxAllowance of 1, or for testing of SingleClaims redemption.
func maxOneRedemptionSteps() []claim {
	return []claim{
		{
			approved,
			[]int64{approvedRedeemableTokenID + 1},
			fmt.Sprintf("ERC721Redeemer: not approved nor owner of %d", approvedRedeemableTokenID+1),
			map[int]int64{},
			nil,
			map[int64]int64{
				approvedRedeemableTokenID:     0,
				approvedRedeemableTokenID + 1: 0,
			},
		},
		{
			approved,
			[]int64{approvedRedeemableTokenID},
			nil,
			map[int]int64{approved: 1},
			[]*ERC721RedeemerRedemption{
				{TokenId: big.NewInt(approvedRedeemableTokenID), N: big.NewInt(1)},
			},
			map[int64]int64{
				approvedRedeemableTokenID: 1,
			},
		},
		{
			tokenOwner,
			[]int64{approvedRedeemableTokenID},
			// Already used by approved address
			fmt.Sprintf("ERC721Redeemer: over allowance for %d", approvedRedeemableTokenID),
			map[int]int64{approved: 1},
			nil,
			map[int64]int64{
				approvedRedeemableTokenID: 1,
			},
		},
		{
			tokenOwner,
			[]int64{1, 2, 3},
			nil,
			map[int]int64{
				approved:   1,
				tokenOwner: 3,
			},
			[]*ERC721RedeemerRedemption{
				{TokenId: big.NewInt(1), N: big.NewInt(1)},
				{TokenId: big.NewInt(2), N: big.NewInt(1)},
				{TokenId: big.NewInt(3), N: big.NewInt(1)},
			},
			map[int64]int64{
				approvedRedeemableTokenID: 1,
				1:                         1,
				2:                         1,
				3:                         1,
			},
		},
		{
			tokenOwner,
			[]int64{2, 4},
			"ERC721Redeemer: over allowance for 2",
			map[int]int64{
				approved:   1,
				tokenOwner: 3,
			},
			nil,
			map[int64]int64{
				approvedRedeemableTokenID: 1,
				1:                         1,
				2:                         1,
				3:                         1,
				4:                         0, // reverted because of 2
			},
		},
	}
}

func TestRedeem(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		maxAllowance int64
		steps        []claim
	}{
		{
			name:         "single redemption allowed",
			maxAllowance: 1,
			steps:        maxOneRedemptionSteps(),
		},
		{
			name:         "5 redemptions allowed",
			maxAllowance: 5,
			steps: []claim{
				{
					asAccountID:       tokenOwner,
					tokenIDs:          []int64{1},
					wantRedeemedAfter: map[int]int64{tokenOwner: 1},
					wantEvents: []*ERC721RedeemerRedemption{
						{TokenId: big.NewInt(1), N: big.NewInt(1)},
					},
					wantClaimedAfter: map[int64]int64{
						1: 1,
						2: 0,
						3: 0,
					},
				},
				{
					asAccountID:       tokenOwner,
					tokenIDs:          []int64{1},
					wantRedeemedAfter: map[int]int64{tokenOwner: 2},
					wantEvents: []*ERC721RedeemerRedemption{
						{TokenId: big.NewInt(1), N: big.NewInt(1)},
					},
					wantClaimedAfter: map[int64]int64{
						1: 2,
						2: 0,
						3: 0,
					},
				},
				{
					asAccountID:       tokenOwner,
					tokenIDs:          []int64{1, 1, 1},
					wantRedeemedAfter: map[int]int64{tokenOwner: 5},
					wantEvents: []*ERC721RedeemerRedemption{
						{TokenId: big.NewInt(1), N: big.NewInt(3)},
					},
					wantClaimedAfter: map[int64]int64{
						1: 5,
						2: 0,
						3: 0,
					},
				},
				{
					asAccountID:       tokenOwner,
					tokenIDs:          []int64{1},
					errDiffAgainst:    "ERC721Redeemer: over allowance for 1",
					wantRedeemedAfter: map[int]int64{tokenOwner: 5},
					wantClaimedAfter: map[int64]int64{
						1: 5,
						2: 0,
						3: 0,
					},
				},
				{
					asAccountID:       tokenOwner,
					tokenIDs:          []int64{2, 3},
					wantRedeemedAfter: map[int]int64{tokenOwner: 7},
					wantEvents: []*ERC721RedeemerRedemption{
						{TokenId: big.NewInt(2), N: big.NewInt(1)},
						{TokenId: big.NewInt(3), N: big.NewInt(1)},
					},
					wantClaimedAfter: map[int64]int64{
						1: 5,
						2: 1,
						3: 1,
					},
				},
				{
					tokenOwner,
					// This is an inefficient way to call; placing identical IDs
					// next to each other will batch the checks.
					[]int64{2, 3, 2, 3},
					nil,
					map[int]int64{tokenOwner: 11},
					[]*ERC721RedeemerRedemption{
						{TokenId: big.NewInt(2), N: big.NewInt(1)},
						{TokenId: big.NewInt(3), N: big.NewInt(1)},
						// Note the repeated tokens with n==1
						{TokenId: big.NewInt(2), N: big.NewInt(1)},
						{TokenId: big.NewInt(3), N: big.NewInt(1)},
					},
					nil,
				},
				{
					tokenOwner,
					// The adjacent identical IDs will be batched for lower gas
					// consumption.
					[]int64{2, 2, 3, 3},
					nil,
					map[int]int64{tokenOwner: 15},
					[]*ERC721RedeemerRedemption{
						// Note the same tokens as the last step but n==2
						{TokenId: big.NewInt(2), N: big.NewInt(2)},
						{TokenId: big.NewInt(3), N: big.NewInt(2)},
					},
					nil,
				},
				{
					tokenOwner,
					[]int64{2},
					"ERC721Redeemer: over allowance for 2",
					map[int]int64{tokenOwner: 15},
					nil,
					nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim, redeem, filter := deployRedeemer(t, 1)
			bigMax := big.NewInt(tt.maxAllowance)

			for i, s := range tt.steps {
				t.Run(fmt.Sprintf("step %d as %s", i, accountName(s.asAccountID)), func(t *testing.T) {
					t.Logf("Redeeming %d", s.tokenIDs)

					tx, err := redeem.RedeemMaxN(sim.Acc(s.asAccountID), bigMax, toBigInts(s.tokenIDs))
					if diff := errdiff.Check(err, s.errDiffAgainst); diff != "" {
						t.Fatalf("RedeemMaxN(%q, maxAllowance=%d, %d) %s", accountName(s.asAccountID), bigMax, s.tokenIDs, diff)
					}

					t.Run("redeemed", func(t *testing.T) {
						gotRedeemed, err := redeem.AllRedeemed(nil)
						if err != nil {
							t.Fatalf("%T.AllRedeemed() error %v", redeem, err)
						}
						got := make(map[common.Address]*big.Int)
						for i, addr := range gotRedeemed.Redeemers {
							got[addr] = gotRedeemed.NumRedeemed[i]
						}

						want := make(map[common.Address]*big.Int)
						for acc, num := range s.wantRedeemedAfter {
							want[sim.Addr(acc)] = big.NewInt(num)
							t.Logf("Address %v is %s", sim.Addr(acc), accountName(acc))
						}

						if diff := cmp.Diff(want, got, ethtest.Comparers()...); diff != "" {
							t.Errorf("%s", diff)
						}
					})

					t.Run("total claimed after", func(t *testing.T) {
						for tokenID, want := range s.wantClaimedAfter {
							got, err := redeem.Claimed(nil, big.NewInt(tokenID))
							if err != nil || got.Cmp(big.NewInt(want)) != 0 {
								t.Errorf("%T.Claimed(token = %d) got %d, err = %v; want %d, nil err", redeem, tokenID, got, err, want)
							}
						}
					})

					if tx == nil {
						return
					}
					ethtest.LogGas(t, tx, fmt.Sprintf("Redeem: %d", s.tokenIDs))

					t.Run("events", func(t *testing.T) {
						testRedemptionLogs(ctx, t, sim, filter, tx, s.wantEvents)
					})
				})

				if t.Failed() {
					// Steps are deliberately not hermetic so stop the test and
					// move to the next one.
					return
				}
			}
		})
	}
}

func TestSingleRedeem(t *testing.T) {
	ctx := context.Background()
	sim, redeem, filter := deployRedeemer(t, 300/mintPerTransaction+1)

	steps := maxOneRedemptionSteps()

	// Include a token ID that is beyond the first word of the BitMap.
	steps = append(steps, claim{
		asAccountID: tokenOwner,
		tokenIDs:    []int64{260, 261},
		wantEvents: []*ERC721RedeemerRedemption{
			{TokenId: big.NewInt(260), N: big.NewInt(1)},
			{TokenId: big.NewInt(261), N: big.NewInt(1)},
		},
	})

	for _, s := range steps {
		tx, err := redeem.RedeemFromSingle(sim.Acc(s.asAccountID), toBigInts(s.tokenIDs))
		if diff := errdiff.Check(err, s.errDiffAgainst); diff != "" {
			t.Fatalf("RedeemFromSingle(%q, %d) %s", accountName(s.asAccountID), s.tokenIDs, diff)
		}

		t.Run("is claimed after", func(t *testing.T) {
			for tokenID, wantNum := range s.wantClaimedAfter {
				// We use the exact same steps for testing the general and the
				// max-single-claim redemption mechanisms, so we need to ensure
				// that this test is valid for conversion to a bool.
				if wantNum != 0 && wantNum != 1 {
					t.Fatalf("BAD TEST SETUP: want %d claims against single-claim token", wantNum)
				}
				want := wantNum == 1

				if got, err := redeem.SingleClaimed(nil, big.NewInt(tokenID)); err != nil || got != want {
					t.Errorf("%T.SingleClaimed(%d) got %t, err = %v; want %t, nil err", redeem, tokenID, got, err, want)
				}
			}
		})

		if tx == nil {
			continue
		}
		testRedemptionLogs(ctx, t, sim, filter, tx, s.wantEvents)
	}
}

func toBigInts(small []int64) []*big.Int {
	var bigs []*big.Int
	for _, s := range small {
		bigs = append(bigs, big.NewInt(s))
	}
	return bigs
}

// testRedemptionLogs compares `want` with the events returned by
// redemptionLogs(tx), reporting non-empty diffs as errors.
func testRedemptionLogs(ctx context.Context, t *testing.T, sim *ethtest.SimulatedBackend, filter *ERC721RedeemerFilterer, tx *types.Transaction, want []*ERC721RedeemerRedemption) {
	t.Helper()
	ignore := cmpopts.IgnoreFields(ERC721RedeemerRedemption{}, "Token", "Redeemer", "Raw")

	gotEvents := redemptionLogs(ctx, t, sim, filter, tx)
	if diff := cmp.Diff(want, gotEvents, ethtest.Comparers(ignore)...); diff != "" {
		t.Errorf("Redemption events diff (-want +got):\n%s", diff)
	}
}

// redemptionLogs waits for the transaction to be mined then returns all
// Redemption logs from the ERC721Redeemer library.
//
// TODO(aschlosberg): the filter iterator wasn't working as expected (perhaps
// due to the event coming from the library and not the contract although this
// seems unlikely) and it was easier to write a quick filter than to debug that.
func redemptionLogs(ctx context.Context, t *testing.T, sim *ethtest.SimulatedBackend, filter *ERC721RedeemerFilterer, tx *types.Transaction) []*ERC721RedeemerRedemption {
	t.Helper()

	topic := crypto.Keccak256([]byte("Redemption(address,address,uint256,uint256)"))

	rcpt, err := bind.WaitMined(ctx, sim, tx)
	if err != nil {
		t.Fatalf("bind.WaitMined() error %v", err)
	}

	var redemptions []*ERC721RedeemerRedemption
	for _, l := range rcpt.Logs {
		if len(l.Topics) == 0 || !bytes.Equal(l.Topics[0].Bytes(), topic) {
			continue
		}

		r, err := filter.ParseRedemption(*l)
		if err != nil {
			t.Errorf("%T.ParseRedemption(%+v) error %v", filter, l, err)
		}
		redemptions = append(redemptions, r)
	}
	return redemptions
}
