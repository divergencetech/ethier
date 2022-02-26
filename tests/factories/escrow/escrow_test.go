package escrow

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate ethier gen ../../../factories/escrow/ColdEscrowFactory.sol ./MockERC721.sol

// Actors in the tests
const (
	deployer = iota
	fallbackOwner
	owner0
	owner1
	owner2
	numAccounts
)

var owners = []int{owner0, owner1, owner2}

func TestColdEscrow(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, numAccounts)

	_, _, factory, err := DeployColdEscrowFactory(sim.Acc(deployer), sim)
	if err != nil {
		t.Fatalf("DeployColdEscrowFactory() error %v", err)
	}

	type nft struct {
		addr common.Address
		m    *MockERC721
	}
	var nfts []nft

	for i := 0; i < 3; i++ {
		addr, _, m, err := DeployMockERC721(sim.Acc(deployer), sim)
		if err != nil {
			t.Fatalf("DeployMockERC721() error %v", err)
		}
		nfts = append(nfts, nft{addr, m})

		// Maintains a mapping from token ID to owner index, so owner0 owns
		// token 0, etc.
		for j, o := range owners {
			sim.Must(t, "%T.Mint() as owner %d", m, j)(m.Mint(sim.Acc(o), big.NewInt(int64(j))))
		}
	}

	sim.Must(t, "Deploy()")(factory.Deploy(sim.Acc(deployer), sim.Addr(fallbackOwner)))

	it, err := factory.FilterColdEscrowDeployed(nil, nil)
	if err != nil {
		t.Fatalf("FilterColdEscrowDeployed() error %v", err)
	}
	if !it.Next() {
		t.Fatalf("Got no %T events after %T.Deploy(); want 1", it.Event, factory)
	}

	escrowAddr := it.Event.ClonedColdEscrow
	escrow, err := NewColdEscrow(escrowAddr, sim)
	if err != nil {
		t.Fatalf("NewColdEscrow([address from factory event]) error %v", err)
	}
	if it.Next() {
		t.Errorf("Got second %T event after %T.Deploy(); want 1", it.Event, factory)
	}
	if err := it.Error(); err != nil {
		t.Errorf("%T.Error() got %v; want nil", it, err)
	}

	if got, err := escrow.Controller(nil); err != nil || !bytes.Equal(got.Bytes(), sim.Addr(fallbackOwner).Bytes()) {
		t.Errorf("Controller() got %v, err = %v; want address passed to %T.Deploy() = %v, nil error", got, err, factory, sim.Addr(fallbackOwner))
	}

	// TODO(aschlosberg): flesh out these tests and abstract all checks for
	// address equality; currently only for proof of concept.
	sim.Must(t, "")(nfts[0].m.SafeTransferFrom(sim.Acc(owner0), sim.Addr(owner0), escrowAddr, big.NewInt(0)))
	if got, err := escrow.Erc721Owners(nil, nfts[0].addr, big.NewInt(0)); err != nil || !bytes.Equal(got.Bytes(), sim.Addr(owner0).Bytes()) {
		t.Errorf("Erc721Owners() got %v, err = %v; want %v, nil err", got, err, sim.Addr(owner0))
	}
	if got, err := nfts[0].m.OwnerOf(nil, big.NewInt(0)); err != nil || !bytes.Equal(got.Bytes(), escrowAddr.Bytes()) {
		t.Errorf("OwnerOf() got %v, err = %v; want escrow contract = %v, nil err", got, err, escrowAddr)
	}
	sim.Must(t, "")(escrow.Reclaim(sim.Acc(owner0), nfts[0].addr, big.NewInt(0)))
	if got, err := nfts[0].m.OwnerOf(nil, big.NewInt(0)); err != nil || !bytes.Equal(got.Bytes(), sim.Addr(owner0).Bytes()) {
		t.Errorf("OwnerOf() got %v, err = %v; want EOA owner = %v, nil err", got, err, sim.Addr(owner0))
	}

	// TODO(aschlosberg): include a test for fallback ownership (controller
	// address) when tokens are transferred with non-safe methods.
}
