package chainlink

import (
	"context"
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/chainlinktest"
	"github.com/divergencetech/ethier/ethtest/chainlinktest/chainlinktestabi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/h-fam/errdiff"
)

func deploy(t *testing.T, requests int64) (*ethtest.SimulatedBackend, *chainlinktest.VRFCoordinator, *TestableVRFConsumerHelper) {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, 1)
	vrf := chainlinktest.DeployAllTB(t, sim)

	addr, _, rand, err := DeployTestableVRFConsumerHelper(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("DeployTestableVRFConsumerHelper() error %v", err)
	}
	chainlinktest.FaucetTB(t, sim, map[common.Address]*big.Int{
		addr: chainlinktest.Fee(requests),
	})

	return sim, vrf, rand
}

func TestSimulatedChainlink(t *testing.T) {
	const requests = 10
	sim, vrf, rand := deploy(t, requests)

	// The SimulatedVRFCoordinator always fulfills randomness with
	// keccak256(requestID) but has to compute the ID itself. This is confirmed
	// here as well as in TestableVRFConsumerHelper; this test also adds
	// coverage for chainlintest.VRFCoordinator to ensure that it fulfills
	// requests.
	for i := 0; i < requests; i++ {
		sim.Must(t, "%T.HelperRequestRandomNess() call [%d]", rand, i)(rand.HelperRequestRandomness(sim.Acc(0)))
		vrf.WaitFulfilledTB(t)

		id, err := rand.LastRequestId(nil)
		if err != nil {
			t.Errorf("%T.LastRequestId() error %v", rand, err)
			continue
		}

		got, err := rand.Randomness(nil, id)
		if err != nil {
			t.Errorf("%T.Randomness() error %v", rand, err)
			continue
		}
		if got.Cmp(big.NewInt(0)) == 0 {
			t.Errorf("Call %d to fulfillRandomness() did not happen", i)
			continue
		}

		want := new(big.Int).SetBytes(crypto.Keccak256(id[:]))
		if got.Cmp(want) != 0 {
			// There's no point in including the actual values as they're all
			// cryptographic hashes.
			t.Errorf("Call %d to fulfillRandomness() got unexpected randomness; SimulatedVRFCoordinator calculates requestId incorrectly", i)
		}
	}
}

func TestExcessGas(t *testing.T) {
	ctx := context.Background()
	// Although we only need the faucet to dispense 2x requests' worth of LINK,
	// this would mean that the last one would receive a gas refund for setting
	// the balance to 0. 100 is arbitrarily large.
	sim, vrf, rand := deploy(t, 100)

	// The very first call incurs additional gas not directly related to the
	// approach, but specific to the test contract.
	sim.Must(t, "%T.HelperRequestRandomness", rand)(rand.HelperRequestRandomness(sim.Acc(0)))
	vrf.WaitFulfilledTB(t)

	const (
		helper = "VRFConsumerHelper"
		base   = "VRFConsumerBase"
	)
	fns := map[string](func(*bind.TransactOpts) (*types.Transaction, error)){
		helper: rand.HelperRequestRandomness,
		base:   rand.StandardRequestRandomness,
	}
	gas := make(map[string]int64)

	for method, fn := range fns {
		tx := sim.Must(t, "%T call to %s.requestRandomness()", rand, method)(fn(sim.Acc(0)))
		vrf.WaitFulfilledTB(t)

		rcpt, err := bind.WaitMined(ctx, sim, tx)
		if err != nil {
			t.Fatalf("with %s; bind.WaitMined() error %v", method, err)
		}
		gas[method] = int64(rcpt.GasUsed)
	}

	t.Log(gas)
	got := gas[helper] - gas[base]
	if got < 0 {
		got = -got
	}
	if got > 1000 {
		// At a gas price of 100 gwei, this is a difference of <0.0001 ETH.
		t.Errorf("Absolute gas differential between VRFConsumerHelper and VRFConsumerBase = %d; want < 1000", got)
	}
}

func TestMultipleConsumers(t *testing.T) {
	sim, vrf, rand := deploy(t, 1)

	addr, _, rand2, err := DeployTestableVRFConsumerHelper(sim.Acc(0), sim)
	if err != nil {
		t.Fatalf("Second deployment: DeployTestableVRFConsumerHelper() error %v", err)
	}
	chainlinktest.FaucetTB(t, sim, map[common.Address]*big.Int{
		addr: chainlinktest.Fee(1),
	})

	// TestableVRFConsumerHelper validates the expected request ID. Running a
	// call to both contracts will exercise this and ensure that the simulated
	// VRFCoordinator keeps nonces in sync on a per-consumer basis.
	sim.Must(t, "On first %T; HelperRequestRandomness()", rand)(rand.HelperRequestRandomness(sim.Acc(0)))
	vrf.WaitFulfilledTB(t)
	sim.Must(t, "On second %T; HelperRequestRandomess()", rand)(rand2.HelperRequestRandomness(sim.Acc(0)))
	vrf.WaitFulfilledTB(t)
}

func TestTransfer(t *testing.T) {
	const faucet = 10
	sim, _, rand := deploy(t, faucet)

	addr := chainlinktest.Addresses().LinkToken
	link, err := chainlinktestabi.NewSimulatedLinkToken(addr, sim)
	if err != nil {
		t.Fatalf("NewSimulatedLink(%v) error %v", addr, err)
	}

	wantBalance := func(t *testing.T, want *big.Int) {
		t.Helper()
		if got, err := link.BalanceOf(nil, sim.Addr(0)); err != nil || got.Cmp(want) != 0 {
			t.Errorf("%T.BalanceOf(%v) got %d, err = %v; want %d, nil err", link, sim.Addr(0), got, err, want)
		}
	}

	tests := []struct {
		wantBefore, transfer, wantAfter *big.Int
		errDiffAgainst                  interface{}
	}{
		{
			wantBefore: big.NewInt(0),
			transfer:   chainlinktest.Fee(1),
			wantAfter:  chainlinktest.Fee(1),
		},
		{
			wantBefore: chainlinktest.Fee(1),
			transfer:   chainlinktest.Fee(2),
			wantAfter:  chainlinktest.Fee(3),
		},
		{
			wantBefore: chainlinktest.Fee(3),
			transfer:   chainlinktest.Fee(3),
			wantAfter:  chainlinktest.Fee(6),
		},
		{
			wantBefore: chainlinktest.Fee(6),
			transfer:   chainlinktest.Fee(4),
			wantAfter:  chainlinktest.Fee(10),
		},
		{
			wantBefore:     chainlinktest.Fee(10),
			transfer:       big.NewInt(1),
			wantAfter:      chainlinktest.Fee(10),
			errDiffAgainst: "ERC20: transfer amount exceeds balance",
		},
	}

	for _, tt := range tests {
		wantBalance(t, tt.wantBefore)

		_, err := rand.TransferLink(sim.Acc(0), sim.Addr(0), tt.transfer)
		if diff := errdiff.Check(err, tt.errDiffAgainst); diff != "" {
			t.Errorf("TransferLink(%d) %s", tt.transfer, diff)
		}

		wantBalance(t, tt.wantAfter)
	}
}
