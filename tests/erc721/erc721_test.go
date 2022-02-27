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

func testApprovalForAllEvents(t *testing.T, filter *ERC721Filterer, start uint64, want []*ERC721ApprovalForAll) {
	ignore := ethtest.Comparers(cmpopts.IgnoreFields(ERC721ApprovalForAll{}, "Raw"))

	iter, err := filter.FilterApprovalForAll(&bind.FilterOpts{
		Start:   start + 1,
		End:     nil,
		Context: nil,
	}, nil, nil)
	if err != nil {
		t.Fatalf("%T.FilterApprovalForAll(nil, nil, nil) error %v", filter, err)
	}
	defer iter.Close()

	var gotEvents []*ERC721ApprovalForAll
	for iter.Next() {
		gotEvents = append(gotEvents, iter.Event)
	}

	if diff := cmp.Diff(want, gotEvents, ignore...); diff != "" {
		t.Errorf("After %T deployment and single ownership transfer; Transfer events diff (-want +got):\n%s", filter, diff)
	}
}

func TestOpenSeaProxyPreApproval(t *testing.T) {
	sim, _, _ := deploy(t)

	tests := []struct {
		name            string
		from            int
		to              int
		hasProxy        bool
		wantPreApproved bool
		wantEvents      []*ERC721ApprovalForAll
	}{
		{
			name:            "mint without proxy",
			from:            -1,
			to:              tokenReceiver,
			hasProxy:        false,
			wantPreApproved: false,
			wantEvents:      nil,
		},
		{
			name:            "mint with proxy",
			from:            -1,
			to:              tokenReceiver,
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
			from:            tokenOwner,
			to:              tokenReceiver,
			hasProxy:        false,
			wantPreApproved: false,
			wantEvents:      nil,
		},
		{
			name:            "transfer with proxy",
			from:            tokenOwner,
			to:              tokenReceiver,
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
			from:            tokenOwner,
			to:              -1,
			hasProxy:        false,
			wantPreApproved: false,
			wantEvents:      nil,
		},
		{
			name:            "burn with proxy",
			from:            tokenOwner,
			to:              -1,
			hasProxy:        false,
			wantPreApproved: false,
			wantEvents:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim, nft, filter := deploy(t)
			start := sim.BlockNumber().Uint64()

			toAddr := common.HexToAddress("0x0")
			if tt.to >= 0 {
				toAddr = sim.Addr(tt.to)
			}

			if tt.hasProxy {
				openseatest.SetProxyTB(t, sim, toAddr, sim.Addr(proxy))
			}

			if tt.from == -1 {
				sim.Must(t, "Mint(%d)", notExists)(nft.Mint(sim.Acc(tt.to), big.NewInt(notExists)))
			} else if tt.to == -1 {
				sim.Must(t, "Burn(%d)", exists)(nft.Burn(sim.Acc(tt.from), big.NewInt(exists)))
			} else {
				sim.Must(t, "TransferFrom(%d)", exists)(nft.TransferFrom(sim.Acc(tt.from), sim.Addr(tt.from), toAddr, big.NewInt(exists)))
			}

			testApprovalForAllEvents(t, filter, start, tt.wantEvents)

			{
				got, err := nft.IsApprovedForAll(nil, toAddr, sim.Addr(proxy))
				if err != nil || got != tt.wantPreApproved {
					t.Errorf("IsApprovedForAll([owner], [proxy]) got %t, err = %v; want %t, nil err", got, err, tt.wantPreApproved)
				}
			}
		})
	}
}

func TestOpenSeaProxyPreApprovalOptInOut(t *testing.T) {
	sim, _, _ := deploy(t)

	tests := []struct {
		name         string
		mints        bool
		hasProxy     bool
		wantApproved bool
		wantEvents   []*ERC721ApprovalForAll
	}{
		{
			name:         "opt in with proxy after mint",
			mints:        true,
			hasProxy:     true,
			wantApproved: true,
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
			name:         "opt out with proxy after mint",
			mints:        true,
			hasProxy:     true,
			wantApproved: false,
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
			name:         "opt in with proxy",
			mints:        false,
			hasProxy:     true,
			wantApproved: true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:         "opt out with proxy",
			mints:        false,
			hasProxy:     true,
			wantApproved: false,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: false,
				},
			},
		},
		{
			name:         "opt in after mint",
			mints:        true,
			hasProxy:     false,
			wantApproved: true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:         "opt out after mint",
			mints:        true,
			hasProxy:     false,
			wantApproved: false,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: false,
				},
			},
		},
		{
			name:         "opt in",
			mints:        false,
			hasProxy:     false,
			wantApproved: true,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: true,
				},
			},
		},
		{
			name:         "opt out",
			mints:        false,
			hasProxy:     false,
			wantApproved: false,
			wantEvents: []*ERC721ApprovalForAll{
				{
					Owner:    sim.Addr(tokenReceiver),
					Operator: sim.Addr(proxy),
					Approved: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim, nft, filter := deploy(t)
			start := sim.BlockNumber().Uint64()

			if tt.hasProxy {
				openseatest.SetProxyTB(t, sim, sim.Addr(tokenReceiver), sim.Addr(proxy))
			}

			if tt.mints {
				sim.Must(t, "Mint(%d)", notExists)(nft.Mint(sim.Acc(tokenReceiver), big.NewInt(notExists)))
			}

			sim.Must(t, "%T.setApprovalForAll(proxy,%v)", nft, tt.wantApproved)(
				nft.SetApprovalForAll(sim.Acc(tokenReceiver), sim.Addr(proxy), tt.wantApproved),
			)

			got, err := nft.IsApprovedForAll(nil, sim.Addr(tokenReceiver), sim.Addr(proxy))
			if err != nil || got != tt.wantApproved {
				t.Errorf("IsApprovedForAll([owner], [proxy]) got %t, err = %v; want %t, nil err", got, err, tt.wantApproved)
			}

			testApprovalForAllEvents(t, filter, start, tt.wantEvents)
		})
	}
}

func TestEnumerableInterface(t *testing.T) {
	enumerableMethods := []string{"TotalSupply", "TokenOfOwnerByIndex", "TokenByIndex"}
	// Common methods, expected on both, are used as a control because there is
	// significant embedding by abigen so we need to know which type actually
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
