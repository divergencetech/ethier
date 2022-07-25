package erc721

import (
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/openseatest"
	"github.com/divergencetech/ethier/ethtest/revert"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

func deploy(t *testing.T) (*ethtest.SimulatedBackend, *TestableERC721ACommon, *ERC721Filterer) {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, numAccounts)
	openseatest.DeployProxyRegistryTB(t, sim)

	addr, _, nft, err := DeployTestableERC721ACommon(sim.Acc(deployer), sim)
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
