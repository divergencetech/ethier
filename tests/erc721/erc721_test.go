package erc721

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/openseatest"
	"github.com/divergencetech/ethier/ethtest/revert"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/h-fam/errdiff"
)

// Actors in the tests
const (
	deployer = iota
	tokenOwner
	tokenOwner2
	tokenReceiver
	approved
	vandal
	proxy
	royaltyReceiver

	numAccounts // last account + 1 ;)
)

func accountName(id int) string {
	switch id {
	case deployer:
		return "contract deployer"
	case tokenOwner:
		return "token owner"
	case approved:
		return "approved for token"
	case vandal:
		return "evil villain"
	case proxy:
		return "OpenSea proxy"
	default:
		return "unknown account"
	}
}

// Token IDs
const (
	exists = iota
	notExists
)

func deploy(t *testing.T) (*ethtest.SimulatedBackend, *TestableERC721ACommon, *ERC721Filterer) {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, numAccounts)
	openseatest.DeployProxyRegistryTB(t, sim)

	addr, _, nft, err := DeployTestableERC721ACommon(sim.Acc(deployer), sim, sim.Addr(royaltyReceiver), big.NewInt(750))
	if err != nil {
		t.Fatalf("TestableERC721ACommon() error %v", err)
	}

	if _, err := nft.Mint(sim.Acc(tokenOwner)); err != nil {
		t.Fatalf("Mint(%d) error %v", exists, err)
	}
	if _, err := nft.Approve(sim.Acc(tokenOwner), sim.Addr(approved), big.NewInt(exists)); err != nil {
		t.Fatalf("Approve(<approved account>, %d) error %v", exists, err)
	}

	filter, err := NewERC721Filterer(addr, sim)
	if err != nil {
		t.Fatalf("NewERC721RedeemerFilterer() error %v", err)
	}

	return sim, nft, filter
}

func TestModifiers(t *testing.T) {
	_, nft, _ := deploy(t)

	tests := []struct {
		name           string
		tokenID        int64
		errDiffAgainst interface{}
	}{
		{
			name:           "existing token",
			tokenID:        exists,
			errDiffAgainst: nil,
		},
		{
			name:           "non-existent token",
			tokenID:        notExists,
			errDiffAgainst: "ERC721ACommon: Token doesn't exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := nft.MustExist(nil, big.NewInt(tt.tokenID))
			if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
				t.Errorf("MustExist([%s]), modified with tokenExists(); %s", tt.name, diff)
			}
		})
	}
}

func TestOnlyApprovedOrOwner(t *testing.T) {
	sim, nft, _ := deploy(t)

	tests := []struct {
		name           string
		account        *bind.TransactOpts
		errDiffAgainst interface{}
	}{
		{
			name:           "token owner",
			account:        sim.Acc(tokenOwner),
			errDiffAgainst: nil,
		},
		{
			name:           "approved",
			account:        sim.Acc(approved),
			errDiffAgainst: nil,
		},
		{
			name:           "vandal",
			account:        sim.Acc(vandal),
			errDiffAgainst: string(revert.ERC721ApproveOrOwner),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := nft.MustBeApprovedOrOwner(tt.account, big.NewInt(exists))
			if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
				t.Errorf("MustBeApprovedOrOwner([%s]), modified with onlyApprovedOrOwner(); %s", tt.name, diff)
			}
		})
	}
}

func TestOpenSeaProxyApproval(t *testing.T) {
	sim, nft, _ := deploy(t)

	tests := []struct {
		name string
		want bool
	}{
		{
			name: "before setting proxy",
			want: false,
		},
		// Note that the proxy is set between the tests
		{
			name: "after setting proxy",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nft.IsApprovedForAll(nil, sim.Addr(deployer), sim.Addr(proxy))
			if err != nil || got != tt.want {
				t.Errorf("IsApprovedForAll([owner], [proxy]) got %t, err = %v; want %t, nil err", got, err, tt.want)
			}
		})

		openseatest.SetProxyTB(t, sim, sim.Addr(deployer), sim.Addr(proxy))
	}
}

func diffApprovalForAllEvents(t *testing.T, filter *ERC721Filterer, start uint64, want []*ERC721ApprovalForAll) string {
	ignore := ethtest.Comparers(cmpopts.IgnoreFields(ERC721ApprovalForAll{}, "Raw"))
	filterOpts := bind.FilterOpts{
		Start:   start + 1,
		End:     nil,
		Context: nil,
	}

	iter, err := filter.FilterApprovalForAll(&filterOpts, nil, nil)
	if err != nil {
		t.Fatalf("FilterApprovalForAll(%v, nil, nil) error %v", filterOpts, err)
	}
	defer func() {
		err := iter.Close()
		if err != nil {
			t.Fatalf("Closing iterator FilterApprovalForAll(%v, nil, nil) error %v", filterOpts, err)
		}
	}()

	var gotEvents []*ERC721ApprovalForAll
	for iter.Next() {
		gotEvents = append(gotEvents, iter.Event)
	}
	if iter.Error() != nil {
		t.Fatalf("Iterating over FilterApprovalForAll(%v, nil, nil) error %v", filterOpts, err)
	}

	return cmp.Diff(want, gotEvents, ignore...)
}

func TestOpenSeaProxyPreApproval(t *testing.T) {
	sim, _, _ := deploy(t)

	const zeroAddressID = -1

	tests := []struct {
		name            string
		fromAccountId   int
		toAccountId     int
		hasProxy        bool
		wantPreApproved bool
		wantEvents      []*ERC721ApprovalForAll
	}{
		{
			name:            "mint without proxy",
			fromAccountId:   zeroAddressID,
			toAccountId:     tokenReceiver,
			hasProxy:        false,
			wantPreApproved: false,
			wantEvents:      nil,
		},
		{
			name:            "mint with proxy",
			fromAccountId:   zeroAddressID,
			toAccountId:     tokenReceiver,
			hasProxy:        true,
			wantPreApproved: true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:            "transfer without proxy",
			fromAccountId:   tokenOwner,
			toAccountId:     tokenReceiver,
			hasProxy:        false,
			wantPreApproved: false,
			wantEvents:      nil,
		},
		{
			name:            "transfer with proxy",
			fromAccountId:   tokenOwner,
			toAccountId:     tokenReceiver,
			hasProxy:        true,
			wantPreApproved: true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:            "burn without proxy",
			fromAccountId:   tokenOwner,
			toAccountId:     zeroAddressID,
			hasProxy:        false,
			wantPreApproved: false,
			wantEvents:      nil,
		},
		{
			name:            "burn with proxy",
			fromAccountId:   tokenOwner,
			toAccountId:     zeroAddressID,
			hasProxy:        false,
			wantPreApproved: false,
			wantEvents:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim, nft, filter := deploy(t)
			start := sim.BlockNumber().Uint64()

			var to common.Address
			if tt.toAccountId != zeroAddressID {
				to = sim.Addr(tt.toAccountId)
			}

			if tt.hasProxy {
				openseatest.SetProxyTB(t, sim, to, sim.Addr(proxy))
			}

			switch {
			case tt.fromAccountId == zeroAddressID:
				// mint
				sim.Must(t, "Mint()")(nft.Mint(sim.Acc(tt.toAccountId)))
			case tt.toAccountId == zeroAddressID:
				// burn
				sim.Must(t, "Burn(%d)", exists)(nft.Burn(sim.Acc(tt.fromAccountId), big.NewInt(exists)))
			default:
				// transfer
				sim.Must(t, "TransferFrom(%d)", exists)(nft.TransferFrom(sim.Acc(tt.fromAccountId), sim.Addr(tt.fromAccountId), to, big.NewInt(exists)))
			}

			if diff := diffApprovalForAllEvents(t, filter, start, tt.wantEvents); diff != "" {
				t.Errorf("ApprovalForAll events diff (-want +got):\n%s", diff)
			}

			if got, err := nft.IsApprovedForAll(nil, to, sim.Addr(proxy)); err != nil || got != tt.wantPreApproved {
				t.Errorf("IsApprovedForAll([owner], [proxy]) got %t, err = %v; want %t, nil err", got, err, tt.wantPreApproved)
			}
		})
	}
}

func TestOpenSeaProxyPreApprovalOptInOut(t *testing.T) {
	sim, _, _ := deploy(t)

	tests := []struct {
		name                 string
		hasProxyBeforeMint   bool
		mints                bool
		createProxyAfterMint bool
		hasProxyAfterMint    bool
		doSetApproved        bool
		setApprovedTo        bool
		wantApproved         bool
		wantEvents           []*ERC721ApprovalForAll
	}{
		{
			name:                 "opt in with proxy after mint",
			hasProxyBeforeMint:   true,
			mints:                true,
			createProxyAfterMint: false,
			doSetApproved:        true,
			setApprovedTo:        true,
			wantApproved:         true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:                 "opt out with proxy after mint",
			hasProxyBeforeMint:   true,
			mints:                true,
			createProxyAfterMint: false,
			doSetApproved:        true,
			setApprovedTo:        false,
			wantApproved:         false,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: false,
				},
			},
		},
		{
			name:                 "opt in with proxy",
			hasProxyBeforeMint:   true,
			mints:                false,
			createProxyAfterMint: false,
			doSetApproved:        true,
			setApprovedTo:        true,
			wantApproved:         true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:                 "opt out with proxy",
			hasProxyBeforeMint:   true,
			mints:                false,
			createProxyAfterMint: false,
			doSetApproved:        true,
			setApprovedTo:        false,
			wantApproved:         false,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: false,
				},
			},
		},
		{
			name:                 "opt in after mint",
			hasProxyBeforeMint:   false,
			mints:                true,
			createProxyAfterMint: false,
			doSetApproved:        true,
			setApprovedTo:        true,
			wantApproved:         true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:                 "opt out after mint",
			hasProxyBeforeMint:   false,
			mints:                true,
			createProxyAfterMint: false,
			doSetApproved:        true,
			setApprovedTo:        false,
			wantApproved:         false,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: false,
				},
			},
		},
		{
			name:                 "opt in",
			hasProxyBeforeMint:   false,
			mints:                false,
			createProxyAfterMint: false,
			doSetApproved:        true,
			setApprovedTo:        true,
			wantApproved:         true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:                 "opt out",
			hasProxyBeforeMint:   false,
			mints:                false,
			createProxyAfterMint: false,
			doSetApproved:        true,
			setApprovedTo:        false,
			wantApproved:         false,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: false,
				},
			},
		},
		{
			name:                 "opt in with proxy after minting without proxy",
			hasProxyBeforeMint:   false,
			mints:                true,
			createProxyAfterMint: true,
			doSetApproved:        true,
			setApprovedTo:        true,
			wantApproved:         true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:                 "opt out with proxy after minting without proxy",
			hasProxyBeforeMint:   false,
			mints:                true,
			createProxyAfterMint: true,
			doSetApproved:        true,
			setApprovedTo:        false,
			wantApproved:         false,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: false,
				},
			},
		},
		{
			name:                 "creating proxy after minting without proxy",
			hasProxyBeforeMint:   false,
			mints:                true,
			createProxyAfterMint: true,
			doSetApproved:        false,
			wantApproved:         false,
			wantEvents:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim, nft, filter := deploy(t)
			start := sim.BlockNumber().Uint64()

			if tt.hasProxyBeforeMint {
				openseatest.SetProxyTB(t, sim, sim.Addr(tokenReceiver), sim.Addr(proxy))
			}

			if tt.mints {
				sim.Must(t, "Mint()")(nft.Mint(sim.Acc(tokenReceiver)))
			}

			if tt.createProxyAfterMint {
				openseatest.SetProxyTB(t, sim, sim.Addr(tokenReceiver), sim.Addr(proxy))
			}

			if tt.doSetApproved {
				sim.Must(t, "%T.setApprovalForAll(proxy,%v)", nft, tt.setApprovedTo)(
					nft.SetApprovalForAll(sim.Acc(tokenReceiver), sim.Addr(proxy), tt.setApprovedTo),
				)
			}

			got, err := nft.IsApprovedForAll(nil, sim.Addr(tokenReceiver), sim.Addr(proxy))
			if err != nil || got != tt.wantApproved {
				t.Errorf("IsApprovedForAll([owner], [proxy]) got %t, err = %v; want %t, nil err", got, err, tt.wantApproved)
			}

			if diff := diffApprovalForAllEvents(t, filter, start, tt.wantEvents); diff != "" {
				t.Errorf("ApprovalForAll events diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRoyalties(t *testing.T) {
	sim, nft, _ := deploy(t)

	tests := []struct {
		setConfig         bool
		newReceiver       common.Address
		newBasisPoints    int64
		salesPrice        int64
		wantRoyaltyAmount int64
	}{
		{
			salesPrice:        200_000,
			wantRoyaltyAmount: 15_000,
		},
		{
			setConfig:         true,
			newReceiver:       sim.Addr(deployer),
			newBasisPoints:    300,
			salesPrice:        700_000,
			wantRoyaltyAmount: 21_000,
		},
	}

	for _, tt := range tests {
		if tt.setConfig {
			sim.Must(t, "nft.SetDefaultRoyalty(%s, %d)", tt.newReceiver, big.NewInt(tt.newBasisPoints))(nft.SetDefaultRoyalty(sim.Acc(deployer), tt.newReceiver, big.NewInt(tt.newBasisPoints)))
		}

		wantReceiver := sim.Addr(royaltyReceiver)

		if tt.setConfig {
			wantReceiver = tt.newReceiver
		}

		receiver, amount, err := nft.RoyaltyInfo(nil, big.NewInt(0), big.NewInt(tt.salesPrice))

		if err != nil || big.NewInt(tt.wantRoyaltyAmount).Cmp(amount) != 0 || wantReceiver != receiver {
			t.Errorf("RoyaltyInfo(0, %d) got (%s, %d), err = %v; want (%s, %d), nil err", tt.salesPrice, receiver, amount, err, wantReceiver, tt.wantRoyaltyAmount)
		}
	}

	if diff := revert.OnlyOwner.Diff(nft.SetDefaultRoyalty(sim.Acc(vandal), sim.Addr(vandal), big.NewInt(1000))); diff != "" {
		t.Errorf("SetDefaultRoyalty([as vandal]) %s", diff)
	}
}

func TestInterfaceSupport(t *testing.T) {
	_, nft, _ := deploy(t)

	tests := []struct {
		interfaceID string
		wantSupport bool
	}{
		{
			interfaceID: "80ac58cd", // ERC721
			wantSupport: true,
		},
		{
			interfaceID: "2a55205a", // ERC2981
			wantSupport: true,
		},
	}

	for _, tt := range tests {
		b, err := hex.DecodeString(tt.interfaceID)
		if err != nil || len(b) != 4 {
			t.Errorf("hex.Decode(%q): err = %v, got len = (%d), want len = 4;", tt.interfaceID, err, len(b))
		}

		var id [4]byte
		copy(id[:], b)

		got, err := nft.SupportsInterface(nil, id)

		if err != nil || got != tt.wantSupport {
			t.Errorf("SupportsInterface(%s) got = %t, err = %v; want = %t", id, got, err, tt.wantSupport)
		}

	}
}

func TestBaseTokenURI(t *testing.T) {
	sim, nft, _ := deploy(t)

	const quantity = 50
	sim.Must(t, "MintN(%d)", quantity)(nft.MintN(sim.Acc(deployer), big.NewInt(quantity)))

	wantURI := func(t *testing.T, id int64, want string) {
		t.Helper()
		got, err := nft.TokenURI(nil, big.NewInt(id))
		if err != nil || got != want {
			t.Errorf("tokenURI(%d) got %q, err = %v; want %q, nil err", id, got, err, want)
		}
	}

	// OpenZeppelin's ERC721 returns an empty string if no base is set.
	wantURI(t, 1, "")
	if diff := revert.OnlyOwner.Diff(nft.SetBaseTokenURI(sim.Acc(vandal), "bad")); diff != "" {
		t.Errorf("SetBaseTokenURI([as vandal]) %s", diff)
	}
	wantURI(t, 1, "")
	wantURI(t, 42, "")

	const base = "good/"
	sim.Must(t, "SetBaseTokenURI(%q)", base)(nft.SetBaseTokenURI(sim.Acc(deployer), base))
	wantURI(t, 7, "good/7")
	wantURI(t, 42, "good/42")
}

func TestPause(t *testing.T) {
	sim, nft, _ := deploy(t)

	sim.Must(t, "Pause()")(nft.Pause(sim.Acc(deployer)))
	check := revert.Checker("ERC721ACommon: paused")
	if diff := check.Diff(nft.MintN(sim.Acc(deployer), big.NewInt(1))); diff != "" {
		t.Errorf("MintN() while paused; %s", diff)
	}
}
