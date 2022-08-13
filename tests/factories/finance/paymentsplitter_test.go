package finance

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/factorytest"
	"github.com/divergencetech/ethier/ethtest/revert"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func deploy(t *testing.T, numAccounts, deploymentAcc int) (*ethtest.SimulatedBackend, *PaymentSplitterFactory, <-chan *PaymentSplitterFactoryPaymentSplitterDeployed) {
	sim := ethtest.NewSimulatedBackendTB(t, numAccounts)

	_, _, dep, err := DeployPaymentSplitterFactory(sim.Acc(deploymentAcc), sim)
	if err != nil {
		t.Fatalf("DeployPaymentSplitterFactory() error %v", err)
	}

	events := make(chan *PaymentSplitterFactoryPaymentSplitterDeployed)
	t.Cleanup(func() { close(events) })
	if _, err := dep.WatchPaymentSplitterDeployed(nil, events); err != nil {
		t.Fatalf("%T.WatchPaymentSplitterDeployed() error %v", dep, err)
	}

	return sim, dep, events
}

func TestPaymentSplitterProxy(t *testing.T) {
	ctx := context.Background()

	const (
		numPayees = 5
		// An extra account acts as the deployer and payer to avoid gas charges
		// confounding tests of balance changes.
		numAccounts = numPayees + 1
		deployer    = 5 // not a payee
	)
	sim, dep, events := deploy(t, numAccounts, deployer)

	tests := []struct {
		name   string
		deploy func(*bind.TransactOpts, []common.Address, []*big.Int) (*types.Transaction, error)
	}{
		{
			name: "raw deploy",
			deploy: func(opts *bind.TransactOpts, payees []common.Address, shares []*big.Int) (*types.Transaction, error) {
				return dep.Deploy(opts, payees, shares)
			},
		},
		{
			name: "deterministic deploy",
			deploy: func(opts *bind.TransactOpts, payees []common.Address, shares []*big.Int) (*types.Transaction, error) {
				var salt [32]byte
				return dep.DeployDeterministic(opts, salt, payees, shares)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				payees      []common.Address
				shares      []*big.Int
				totalShares int64
			)
			for i := 0; i < numPayees; i++ {
				payees = append(payees, sim.Addr(i))

				s := int64(i + 1)
				shares = append(shares, big.NewInt(s))
				totalShares += s
			}
			t.Logf("Payees: %v", payees)
			t.Logf("Shares: %d", shares)

			if _, err := tt.deploy(sim.Acc(deployer), payees, shares); err != nil {
				t.Fatalf("Deploy(%v, %d) error %v", payees, shares, err)
			}
			ev, ok := <-events
			if !ok {
				t.Fatal("Event channel closed unexpectedly")
			}
			cloned := ev.ClonedPaymentSplitter

			t.Run("no re-init", func(t *testing.T) {
				d, err := NewDelegatedPaymentSplitter(cloned, sim)
				if err != nil {
					t.Fatalf("NewDelegatedPaymentSplitter(address from Deploy() event) error %v", err)
				}

				if diff := revert.AlreadyInitialized.Diff(d.Initialize(sim.Acc(deployer), payees[:2], shares[:2])); diff != "" {
					t.Errorf("%T.Initialize() when attempting to overwrite payees; %s", d, diff)
				}
			})

			// Note that although we have deployed a proxy to a
			// DelegatedPaymentSplitter, it has an otherwise identical function
			// signature so we test with the standard contract.
			split, err := NewPaymentSplitter(cloned, sim)
			if err != nil {
				t.Fatalf("NewPaymentSplitter(address from Deploy() event) error %v", err)
			}

			t.Run("constants", func(t *testing.T) {
				if got, err := split.TotalShares(nil); err != nil || got.Cmp(big.NewInt(totalShares)) != 0 {
					t.Errorf("%T.TotalShares() got %d, err = %v; want %d, nil err", split, got, err, totalShares)
				}

				for i := 0; i < numPayees; i++ {
					if got, err := split.Payee(nil, big.NewInt(int64(i))); err != nil || !bytes.Equal(got.Bytes(), payees[i].Bytes()) {
						t.Errorf("%T.Payee(%d) got %v, err = %v; want %v, nil err", split, i, got, err, payees[i])
					}
					if got, err := split.Shares(nil, payees[i]); err != nil || got.Cmp(shares[i]) != 0 {
						t.Errorf("%T.Shares(%d) got %d, err = %v; want %d, nil err", split, i, got, err, shares[i])
					}
				}
			})

			t.Run("payment and splitting", func(t *testing.T) {
				raw := &PaymentSplitterRaw{split}
				sim.Must(t, "%T.Transfer(1 ETH)", raw)(raw.Transfer(sim.WithValueFrom(deployer, eth.Ether(1))))
				if got, want := sim.BalanceOf(ctx, t, cloned), eth.Ether(1); got.Cmp(want) != 0 {
					t.Fatalf("Bad setup; balance of cloned PaymentSplitter got %d; want %d", got, want)
				}

				for i := 0; i < numPayees; i++ {
					before := sim.BalanceOf(ctx, t, payees[i])
					sim.Must(t, "%T.Release(payee %d)", split, i)(split.Release(sim.Acc(deployer), payees[i]))
					after := sim.BalanceOf(ctx, t, payees[i])

					want := eth.EtherFraction(shares[i].Int64(), totalShares)
					if got := new(big.Int).Sub(after, before); got.Cmp(want) != 0 {
						t.Errorf("After %T.Release() payee %d balance diff of %d (%d to %d); want diff of %d", split, i, got, before, after, want)
					}
				}
			})

		})
	}
}

func TestDeterministicDeploymentAddress(t *testing.T) {
	sim, dep, events := deploy(t, 1, 0)

	payees := []common.Address{sim.Addr(0)}
	shares := []*big.Int{big.NewInt(1)}

	for i := byte(0); i < 10; i++ {
		var salt [32]byte
		salt[0] = i

		predicted, err := dep.PredictDeploymentAddress(nil, salt)
		if err != nil {
			t.Errorf("%T.PredictDeploymentAddress(%#x) error %v", dep, salt, err)
			continue
		}

		sim.Must(t, "%T.DeployDeterministic(%#x)", dep, salt)(dep.DeployDeterministic(sim.Acc(0), salt, payees, shares))
		ev, ok := <-events
		if !ok {
			t.Fatal("Event channel closed unexpectedly")
		}

		if got := ev.ClonedPaymentSplitter; !bytes.Equal(got.Bytes(), predicted.Bytes()) {
			t.Errorf("%T.DeployDeterministic(%#x) got deployed address %v; want %v as returned by PredictDeploymentAddress()", dep, salt, got, predicted)
		}
	}
}

func TestChainAgnosticLibrary(t *testing.T) {
	ctx := context.Background()
	sim := ethtest.NewSimulatedBackendTB(t, 3)

	// Deploys an instance of the PaymentSplitterFactory to a deterministic
	// address, similarly to how there will be known addresses of main- and
	// test-net deployments. All of these addresses are made available in the
	// PaymentSplitterDeployer library, which determines the correct address
	// based on the chain ID.
	factorytest.DeployPaymentSplitterFactoryTB(t, sim)

	// Deploys a contract that uses the aforementioned library.
	_, _, deployer, err := DeployTestablePaymentSplitterDeployer(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployTestablePaymentSplitterDeployer() error %v", err)
	}

	var salt [32]byte
	addr, err := deployer.PredictDeploymentAddress(nil, salt)
	if err != nil {
		t.Fatalf("%T.PredictDeploymentAddress(%#x) error %v", deployer, salt, err)
	}

	tx, err := deployer.DeployDeterministic(sim.Acc(0), salt, []common.Address{sim.Addr(1), sim.Addr(2)}, []*big.Int{big.NewInt(1), big.NewInt(2)})
	if err != nil {
		t.Fatalf("%T.DeployDeterministic(salt %#x) error %v", deployer, salt, err)
	}
	rcpt, err := bind.WaitMined(ctx, sim, tx)
	if err != nil {
		t.Fatalf("bind.WaitMined(transaction from deployment with library) error %v", err)
	}
	got := rcpt.Logs

	topics := func(event string) []common.Hash {
		h := common.BytesToHash(crypto.Keccak256([]byte(event)))
		t.Logf("Event %q -> topic hash %v", event, h)
		return []common.Hash{h}
	}

	want := []*types.Log{
		{
			Topics: topics("PayeeAdded(address,uint256)"),
			Data:   append(common.LeftPadBytes(sim.Addr(1).Bytes(), 32), common.LeftPadBytes([]byte{1}, 32)...),
		},
		{
			Topics: topics("PayeeAdded(address,uint256)"),
			Data:   append(common.LeftPadBytes(sim.Addr(2).Bytes(), 32), common.LeftPadBytes([]byte{2}, 32)...),
		},
		{
			Topics: topics("Initialized(uint8)"),
			Data:   common.LeftPadBytes([]byte{1}, 32),
		},
		{
			Topics: topics("PaymentSplitterDeployed(address)"),
			Data:   common.LeftPadBytes(addr.Bytes(), 32),
		},
	}

	ignore := cmpopts.IgnoreFields(types.Log{}, "Address", "BlockHash", "BlockNumber", "Index", "TxHash")

	if diff := cmp.Diff(want, got, ignore); diff != "" {
		t.Errorf("Events emitted by %T.DeployDeterministic(salt %#x) (-want +got) diff:\n%s", deployer, salt, diff)
	}
}
