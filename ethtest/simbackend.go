// The ethtest package provides helpers for testing Ethereum smart contracts.
package ethtest

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/divergencetech/ethier/eth"
)

// A SimulatedBackend embeds a go-ethereum SimulatedBackend and extends its
// functionality to simplify standard testing.
type SimulatedBackend struct {
	*backends.SimulatedBackend

	AutoCommit bool
	accounts   []*bind.TransactOpts

	// See comment on MockedEntity.
	mockAccounts map[MockedEntity]*bind.TransactOpts
}

var _ bind.ContractBackend = (*SimulatedBackend)(nil)

// NewSimulatedBackend returns a new simulated ETH backend with the specified
// number of accounts. Transactions are automatically committed unless. Close()
// must be called to free resources after use.
//
// Accounts are deterministically generated so have identical addresses between
// backends, but balances are coupled to the specific instance of the backend.
func NewSimulatedBackend(numAccounts int) (*SimulatedBackend, error) {
	sb := &SimulatedBackend{
		AutoCommit:   true,
		mockAccounts: make(map[MockedEntity]*bind.TransactOpts),
	}
	alloc := make(core.GenesisAlloc)

	// Ensure that the pre-compiled contracts are available.
	// TODO: check if this is absolutely necessary.
	for addr := byte(1); addr <= 8; addr++ {
		alloc[common.BytesToAddress([]byte{addr})] = core.GenesisAccount{
			Balance: big.NewInt(1),
		}
	}

	createAccount := func(seed []byte) (*bind.TransactOpts, error) {
		entropy := bytes.NewReader(crypto.Keccak512(seed))
		key, err := ecdsa.GenerateKey(crypto.S256(), entropy)
		if err != nil {
			return nil, fmt.Errorf("ecdsa.GenerateKey(crypto.S256, [deterministic entropy; Keccak512(%q)]): %v", seed, err)
		}

		txOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		if err != nil {
			return nil, fmt.Errorf("NewKeyedTransactorWithChainID(<new key>, sim-backend-id=1337): %v", err)
		}
		alloc[txOpts.From] = core.GenesisAccount{
			Balance: eth.Ether(100),
		}
		return txOpts, nil
	}

	for i := 0; i < numAccounts; i++ {
		txOpts, err := createAccount([]byte(fmt.Sprintf("account:%d", i)))
		if err != nil {
			return nil, err
		}
		sb.accounts = append(sb.accounts, txOpts)
	}

	// These accounts need to be deterministic so that any contracts they deploy
	// have deterministic addresses.
	for _, mock := range []MockedEntity{OpenSea, Chainlink} {
		txOpts, err := createAccount([]byte(mock))
		if err != nil {
			return nil, err
		}
		sb.mockAccounts[mock] = txOpts
	}

	sb.SimulatedBackend = backends.NewSimulatedBackend(alloc, 3e7)

	sb.AdjustTime(365 * 24 * time.Hour)
	sb.Commit()

	return sb, nil
}

// NewSimulatedBackendTB calls NewSimulatedBackend(), reports any errors with
// tb.Fatal, and calls Close() with tb.Cleanup().
func NewSimulatedBackendTB(tb testing.TB, numAccounts int) *SimulatedBackend {
	tb.Helper()

	sim, err := NewSimulatedBackend(numAccounts)
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() {
		if err := sim.Close(); err != nil {
			tb.Errorf("%T.Close(): %v", sim.SimulatedBackend, err)
		}
	})

	return sim
}

// SendTransaction functions pipes its parameters to the embedded backend and
// also calls Commit() if sb.AutoCommit==true.
func (sb *SimulatedBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if err := sb.SimulatedBackend.SendTransaction(ctx, tx); err != nil {
		return err
	}
	if sb.AutoCommit {
		sb.SimulatedBackend.Commit()
	}
	return nil
}

// Acc returns a TransactOpts signing as the specified account number.
func (sb *SimulatedBackend) Acc(account int) *bind.TransactOpts {
	acc := sb.accounts[account]
	return &bind.TransactOpts{
		From:   acc.From,
		Signer: acc.Signer,
	}
}

// Addr returns the Address of the specified account number.
func (sb *SimulatedBackend) Addr(account int) common.Address {
	return sb.accounts[account].From
}

// WithValueFrom returns a TransactOpts that sends the specified value from the
// account. If value==0, sb.Acc(account) can be used directly.
func (sb *SimulatedBackend) WithValueFrom(account int, value *big.Int) *bind.TransactOpts {
	opts := sb.Acc(account)
	opts.Value = value
	return opts
}

// CallFrom returns a CallOpts from the specified account number.
func (sb *SimulatedBackend) CallFrom(account int) *bind.CallOpts {
	return &bind.CallOpts{
		From: sb.accounts[account].From,
	}
}

// A MockedEntity mocks a real-world entity such as Uniswap or Opensea with
// deterministically generated accounts (and therefore contract addresses).
type MockedEntity string

// Mocked entities.
const (
	OpenSea   = MockedEntity("OpenSea")
	Chainlink = MockedEntity("Chainlink")
)

// AsMockedEntity calls the provided function with the mocked entity's account
// if it is supported, propagating any returned error.
//
// MockedEntity accounts SHOULD NOT be used in general tests; prefer provided
// packages like openseatest to using AsMockedEntity() directly.
func (sb *SimulatedBackend) AsMockedEntity(mock MockedEntity, fn func(*bind.TransactOpts) error) error {
	a, ok := sb.mockAccounts[mock]
	if !ok {
		return fmt.Errorf("unsupported %T %q", mock, mock)
	}

	return fn(a)
}

// BalanceOf returns the current balance of the address, calling tb.Fatalf on
// error.
func (sb *SimulatedBackend) BalanceOf(ctx context.Context, tb testing.TB, addr common.Address) *big.Int {
	tb.Helper()
	bal, err := sb.BalanceAt(ctx, addr, nil)
	if err != nil {
		tb.Fatalf("%T.BalanceAt(ctx, %s, nil) error %v", sb.SimulatedBackend, addr, err)
	}
	return bal
}

// BlockNumber returns the current block number.
func (sb *SimulatedBackend) BlockNumber() *big.Int {
	return sb.Blockchain().CurrentBlock().Number()
}

// FastForward calls sb.Commit() until sb.BlockNumber() >= blockNumber. It
// returns whether fast-forwarding was required; i.e. false if the requested
// block number is current or in the past.
//
// NOTE: FastForward is O(curr - requested).
func (sb *SimulatedBackend) FastForward(blockNumber *big.Int) bool {
	done := false
	for ; blockNumber.Cmp(sb.BlockNumber()) == 1; done = true {
		// TODO: is there a more efficient way to do this?
		sb.Commit()
	}
	return done
}

// GasSpent returns the gas spent (i.e. used*cost) by the transaction.
func (sb *SimulatedBackend) GasSpent(ctx context.Context, tb testing.TB, tx *types.Transaction) *big.Int {
	rcpt, err := bind.WaitMined(ctx, sb, tx)
	if err != nil {
		tb.Fatalf("bind.WaitMined(<simulated backend>, %s) error %v", tx.Hash(), err)
	}
	return new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(rcpt.GasUsed))
}

// Must returns a function that ensures a successful transaction, reporting any
// error on tb.Fatal, or propagating the transaction.
//
// Intended usage:
//  sb.Must(t, "ContractFunc()")(foo.ContractFunc(sim.Acc(<acc>), …))
//
// The description format and associated args will be used as a prefix in any
// reported errors. The returned function MUST be used immediately, and can only
// be used once.
func (sb *SimulatedBackend) Must(tb testing.TB, descFormat string, descArgs ...interface{}) func(*types.Transaction, error) *types.Transaction {
	// This function is "naughty" and not strictly within idiomatic Go
	// style. Similarly to how contexts mustn't be held within structs, holding
	// a testing.T risks it becoming irrelevant with respect to the scope within
	// which it's used. To avoid this, we limit the returned function to single
	// use; it's not a perfect solution, but a user would have to deliberately
	// misuse the API.
	var used bool

	desc := fmt.Sprintf(descFormat, descArgs...)

	return func(tx *types.Transaction, err error) *types.Transaction {
		tb.Helper()
		if used {
			tb.Errorf("Function returned by %T.Must(%q) must only be used once", sb, desc)
		}
		used = true

		if err != nil {
			tb.Fatalf("%s; got err %v; want nil err", desc, err)
			return nil
		}
		return tx
	}
}
