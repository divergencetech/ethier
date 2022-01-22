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
	numRedeemableTokens = 10
	// All tokens are allowed to redeem claims, but this one has an ERC721
	// approved alternative address.
	approvedRedeemableTokenID = 0
)

func deployRedeemer(t *testing.T) (*ethtest.SimulatedBackend, *TestableERC721Redeemer, *ERC721RedeemerFilterer) {
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

	sim.Must(t, "Mint(tokenOwner, %d)", numRedeemableTokens)(
		redeem.Mint(
			sim.Acc(deployer),
			sim.Addr(tokenOwner),
			big.NewInt(numRedeemableTokens),
		),
	)
	sim.Must(t, "%T.Approve(%d)", token, approvedRedeemableTokenID)(
		token.Approve(
			sim.Acc(tokenOwner),
			sim.Addr(approved),
			big.NewInt(approvedRedeemableTokenID),
		),
	)

	return sim, redeem, filter
}

func TestRedeem(t *testing.T) {
	ctx := context.Background()

	type claim struct {
		asAccountID        int
		tokenIDs           []int64
		errDiffAgainst     interface{}
		wantRedeemedAfter  map[int]int64
		wantEvents         []*ERC721RedeemerRedemption
		wantUnclaimedAfter map[int64]int64
	}

	tests := []struct {
		name         string
		maxAllowance int64
		steps        []claim
	}{
		{
			name:         "single redemption allowed",
			maxAllowance: 1,
			steps: []claim{
				{
					approved,
					[]int64{approvedRedeemableTokenID + 1},
					fmt.Sprintf("ERC721Redeemer: not approved nor owner of %d", approvedRedeemableTokenID+1),
					map[int]int64{},
					nil,
					nil,
				},
				{
					approved,
					[]int64{approvedRedeemableTokenID},
					nil,
					map[int]int64{approved: 1},
					[]*ERC721RedeemerRedemption{
						{TokenId: big.NewInt(approvedRedeemableTokenID), N: big.NewInt(1)},
					},
					nil,
				},
				{
					tokenOwner,
					[]int64{approvedRedeemableTokenID},
					// Already used by approved address
					fmt.Sprintf("ERC721Redeemer: over allowance for %d", approvedRedeemableTokenID),
					map[int]int64{approved: 1},
					nil,
					nil,
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
					nil,
				},
				{
					tokenOwner,
					[]int64{2},
					"ERC721Redeemer: over allowance for 2",
					map[int]int64{
						approved:   1,
						tokenOwner: 3,
					},
					nil,
					nil,
				},
			},
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
					wantUnclaimedAfter: map[int64]int64{
						1: 4,
						2: 5,
						3: 5,
					},
				},
				{
					asAccountID:       tokenOwner,
					tokenIDs:          []int64{1},
					wantRedeemedAfter: map[int]int64{tokenOwner: 2},
					wantEvents: []*ERC721RedeemerRedemption{
						{TokenId: big.NewInt(1), N: big.NewInt(1)},
					},
					wantUnclaimedAfter: map[int64]int64{
						1: 3,
						2: 5,
						3: 5,
					},
				},
				{
					asAccountID:       tokenOwner,
					tokenIDs:          []int64{1, 1, 1},
					wantRedeemedAfter: map[int]int64{tokenOwner: 5},
					wantEvents: []*ERC721RedeemerRedemption{
						{TokenId: big.NewInt(1), N: big.NewInt(3)},
					},
					wantUnclaimedAfter: map[int64]int64{
						1: 0,
						2: 5,
						3: 5,
					},
				},
				{
					asAccountID:       tokenOwner,
					tokenIDs:          []int64{1},
					errDiffAgainst:    "ERC721Redeemer: over allowance for 1",
					wantRedeemedAfter: map[int]int64{tokenOwner: 5},
					wantUnclaimedAfter: map[int64]int64{
						1: 0,
						2: 5,
						3: 5,
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
					wantUnclaimedAfter: map[int64]int64{
						1: 0,
						2: 4,
						3: 4,
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
			sim, redeem, filter := deployRedeemer(t)
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

					t.Run("unclaimed", func(t *testing.T) {
						for tokenID, want := range s.wantUnclaimedAfter {
							got, err := redeem.UnclaimedIfMaxN(nil, bigMax, big.NewInt(tokenID))
							if err != nil || got.Cmp(big.NewInt(want)) != 0 {
								t.Errorf("%T.UnclaimedIfMaxN(%d, token = %d) got %d, err = %v; want %d, nil err", redeem, bigMax, tokenID, got, err, want)
							}
						}
					})

					if tx == nil {
						return
					}
					ethtest.LogGas(t, tx, fmt.Sprintf("Redeem: %d", s.tokenIDs))

					t.Run("events", func(t *testing.T) {
						ignore := cmpopts.IgnoreFields(ERC721RedeemerRedemption{}, "Token", "Redeemer", "Raw")

						gotEvents := redemptionLogs(ctx, t, sim, filter, tx)
						if diff := cmp.Diff(s.wantEvents, gotEvents, ethtest.Comparers(ignore)...); diff != "" {
							t.Errorf("Redemption events diff (-want +got):\n%s", diff)
						}
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

func toBigInts(small []int64) []*big.Int {
	var bigs []*big.Int
	for _, s := range small {
		bigs = append(bigs, big.NewInt(s))
	}
	return bigs
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
