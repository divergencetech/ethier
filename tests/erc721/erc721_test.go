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
	"github.com/h-fam/errdiff"
)

// Actors in the tests
const (
	deployer = iota
	admin
	steerer
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
	case admin:
		return "contract admin"
	case steerer:
		return "contract steerer"
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
	adminRoleName    = "DEFAULT_ADMIN_ROLE"
	steeringRoleName = "DEFAULT_STEERING_ROLE"
)

type deployed struct {
	sim                     *ethtest.SimulatedBackend
	nft                     *TestableERC721ACommon
	filter                  *ERC721Filterer
	adminRole, steeringRole [32]byte
}

func deploy(t *testing.T) deployed {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, numAccounts)
	openseatest.DeployProxyRegistryTB(t, sim)

	addr, _, nft, err := DeployTestableERC721ACommon(sim.Acc(deployer), sim, sim.Addr(admin), sim.Addr(steerer), sim.Addr(royaltyReceiver), big.NewInt(750))
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

	adminRole, err := nft.DEFAULTADMINROLE(nil)
	if err != nil {
		t.Fatalf("nft.DEFAULT_ADMIN_ROLE(): %v", err)
	}
	steeringRole, err := nft.DEFAULTSTEERINGROLE(nil)
	if err != nil {
		t.Fatalf("nft.DEFAULT_STEERING_ROLE(): %v", err)
	}

	return deployed{sim, nft, filter, adminRole, steeringRole}
}

func TestModifiers(t *testing.T) {
	d := deploy(t)
	nft := d.nft

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
	dep := deploy(t)
	sim := dep.sim
	nft := dep.nft

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

func TestRoyalties(t *testing.T) {
	dep := deploy(t)
	sim := dep.sim
	nft := dep.nft

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
			sim.Must(t, "nft.SetDefaultRoyalty(%s, %d)", tt.newReceiver, big.NewInt(tt.newBasisPoints))(nft.SetDefaultRoyalty(sim.Acc(steerer), tt.newReceiver, big.NewInt(tt.newBasisPoints)))
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

	if diff := revert.MissingRoleByName(sim.Addr(vandal), steeringRoleName).Diff(nft.SetDefaultRoyalty(sim.Acc(vandal), sim.Addr(vandal), big.NewInt(1000))); diff != "" {
		t.Errorf("SetDefaultRoyalty([as vandal]) %s", diff)
	}
}

func TestInterfaceSupport(t *testing.T) {
	dep := deploy(t)
	nft := dep.nft

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
	dep := deploy(t)
	sim := dep.sim
	nft := dep.nft

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
	if diff := revert.MissingRoleByName(sim.Addr(vandal), steeringRoleName).Diff(nft.SetBaseTokenURI(sim.Acc(vandal), "bad")); diff != "" {
		t.Errorf("SetBaseTokenURI([as vandal]) %s", diff)
	}
	wantURI(t, 1, "")
	wantURI(t, 42, "")

	const base = "good/"
	sim.Must(t, "SetBaseTokenURI(%q)", base)(nft.SetBaseTokenURI(sim.Acc(steerer), base))
	wantURI(t, 7, "good/7")
	wantURI(t, 42, "good/42")
}

func TestPause(t *testing.T) {
	dep := deploy(t)
	sim := dep.sim
	nft := dep.nft

	sim.Must(t, "Pause()")(nft.Pause(sim.Acc(steerer)))
	check := revert.Checker("ERC721ACommon: paused")
	if diff := check.Diff(nft.MintN(sim.Acc(deployer), big.NewInt(1))); diff != "" {
		t.Errorf("MintN() while paused; %s", diff)
	}
}

func TestRoleAssignment(t *testing.T) {
	dep := deploy(t)
	sim := dep.sim
	nft := dep.nft

	sim.Must(t, "grantRole(STEERING_ROLE, approved)")(nft.GrantRole(sim.Acc(admin), dep.steeringRole, sim.Addr(approved)))

	if got, err := nft.HasRole(nil, dep.steeringRole, sim.Addr(approved)); err != nil || !got {
		t.Errorf("HasRole(STEERING_ROLE, approved) got %t, err = %v; want %t, nil err", got, err, true)
	}

	check := revert.MissingRole(sim.Addr(approved), dep.adminRole)
	if diff := check.Diff(nft.GrantRole(sim.Acc(approved), dep.adminRole, sim.Addr(approved))); diff != "" {
		t.Errorf("grantRole(ADMIN_ROLE, approved, [as STEERING_ROLE]) %s", diff)
	}
}
