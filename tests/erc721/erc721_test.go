package erc721

import (
	"math/big"
	"reflect"
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

func deploy(t *testing.T) (*ethtest.SimulatedBackend, *TestableERC721CommonEnumerable, *ERC721Filterer) {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, numAccounts)
	openseatest.DeployProxyRegistryTB(t, sim)

	addr, _, nft, err := DeployTestableERC721CommonEnumerable(sim.Acc(deployer), sim)
	if err != nil {
		t.Fatalf("DeployTestableERC721CommonEnumerable() error %v", err)
	}

	if _, err := nft.Mint(sim.Acc(tokenOwner), big.NewInt(exists)); err != nil {
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
			errDiffAgainst: "ERC721Common: Token doesn't exist",
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
				sim.Must(t, "Mint(%d)", notExists)(nft.Mint(sim.Acc(tt.toAccountId), big.NewInt(notExists)))
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
				sim.Must(t, "Mint(%d)", notExists)(nft.Mint(sim.Acc(tokenReceiver), big.NewInt(notExists)))
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

func TestEnumerableInterface(t *testing.T) {
	enumerableMethods := []string{"TotalSupply", "TokenOfOwnerByIndex", "TokenByIndex"}
	// Common methods, expected on both, are used as a control because there is
	// significant embedding by abigen so we need toAccountId know which type actually
	// has the methods.
	commonMethods := []string{"BalanceOf", "OwnerOf", "IsApprovedForAll"}

	tests := []struct {
		contract interface{}
		methods  []string
		want     bool
	}{
		{
			contract: &ERC721CommonCaller{},
			methods:  commonMethods,
			want:     true,
		},
		{
			contract: &ERC721CommonEnumerableCaller{},
			methods:  commonMethods,
			want:     true,
		},
		{
			contract: &ERC721CommonCaller{},
			methods:  enumerableMethods,
			want:     false,
		},
		{
			contract: &ERC721CommonEnumerableCaller{},
			methods:  enumerableMethods,
			want:     true,
		},
	}

	for _, tt := range tests {
		typ := reflect.TypeOf(tt.contract)
		for _, method := range tt.methods {
			if _, got := typ.MethodByName(method); got != tt.want {
				t.Errorf("%T has method %q? got %t; want %t", tt.contract, method, got, tt.want)
			}
		}
	}

	// Effectively the same as above but this won't compile if we haven't
	// inherited properly.
	var enum ERC721CommonEnumerable
	_ = []interface{}{
		enum.TotalSupply,
		enum.TokenOfOwnerByIndex,
		enum.TokenByIndex,
	}
}

func TestBaseTokenURI(t *testing.T) {
	sim, nft, _ := deploy(t)

	for _, id := range []int64{1, 42, 101010} {
		sim.Must(t, "Mint(%d)", id)(nft.Mint(sim.Acc(deployer), big.NewInt(id)))
	}

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
	wantURI(t, 42, "good/42")
	wantURI(t, 101010, "good/101010")
}

func TestAutoIncrement(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, numAccounts)
	openseatest.DeployProxyRegistryTB(t, sim)

	_, _, nft, err := DeployTestableERC721AutoIncrement(sim.Acc(deployer), sim)
	if err != nil {
		t.Fatalf("DeployTestableERC721AutoIncrement() error %v", err)
	}

	const (
		o1 = tokenOwner
		o2 = tokenOwner2
	)

	tests := []struct {
		toAccountIndex   int
		n                int64
		wantOwnerIndices []int
		wantTotalSupply  int64
	}{
		{
			toAccountIndex:   o1,
			n:                2,
			wantOwnerIndices: []int{o1, o1},
			wantTotalSupply:  2,
		},
		{
			toAccountIndex:   o2,
			n:                0,
			wantOwnerIndices: []int{o1, o1},
			wantTotalSupply:  2,
		},
		{
			toAccountIndex:   o2,
			n:                1,
			wantOwnerIndices: []int{o1, o1, o2},
			wantTotalSupply:  3,
		},
		{
			toAccountIndex:   o2,
			n:                3,
			wantOwnerIndices: []int{o1, o1, o2, o2, o2, o2},
			wantTotalSupply:  6,
		},
		{
			toAccountIndex:   o1,
			n:                1,
			wantOwnerIndices: []int{o1, o1, o2, o2, o2, o2, o1},
			wantTotalSupply:  7,
		},
	}

	for _, tt := range tests {
		to := sim.Addr(tt.toAccountIndex)
		sim.Must(t, "safeMintN(%v, %d)", to, tt.n)(nft.SafeMintN(sim.Acc(deployer), to, big.NewInt(tt.n)))

		t.Run("owners", func(t *testing.T) {
			gotOwners, err := nft.AllOwners(nil)
			if err != nil {
				t.Fatalf("AllOwners() error %v", err)
			}

			var wantOwners []common.Address
			for _, i := range tt.wantOwnerIndices {
				wantOwners = append(wantOwners, sim.Addr(i))
			}

			if diff := cmp.Diff(wantOwners, gotOwners); diff != "" {
				t.Errorf("AllOwners() diff (-want +got):\n%s", diff)
			}
		})

		if got, err := nft.TotalSupply(nil); err != nil || got.Cmp(big.NewInt(tt.wantTotalSupply)) != 0 {
			t.Errorf("TotalSupply() got %d, err = %v; want %d, nil err", got, err, tt.wantTotalSupply)
		}
	}
}
