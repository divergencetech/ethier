package weth

import (
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/revert"
	"github.com/divergencetech/ethier/ethtest/wethtest"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate ethier gen TestableWETH.sol

const (
	deployer = iota
	hodler0
	hodler1
	numAccounts
)

var hodlers = []int{hodler0, hodler1}

func TestSimulatedWETH(t *testing.T) {
	sim := ethtest.NewSimulatedBackendTB(t, numAccounts)
	weth := wethtest.DeployWETHTB(t, sim)

	t.Run("insufficient balance revert", func(t *testing.T) {
		if diff := revert.WETHWithdraw.Diff(weth.Withdraw(sim.Acc(hodler0), eth.Ether(1))); diff != "" {
			t.Errorf("%T.Withdraw() with 0 balance diff %s", weth, diff)
		}
		if diff := revert.WETHTransfer.Diff(weth.Transfer(sim.Acc(hodler0), sim.Addr(hodler1), eth.Ether(1))); diff != "" {
			t.Errorf("%T.Transfer() with 0 balance diff %s", weth, diff)
		}
	})

	const deposit = 10

	wantBalance := func(t *testing.T, accountID int, want *big.Int) {
		t.Helper()
		addr := sim.Addr(accountID)
		if got, err := weth.BalanceOf(nil, addr); err != nil || got.Cmp(want) != 0 {
			t.Errorf("weth.BalanceOf(account %d = %v) got %d, err = %v; want %d, nil err", accountID, addr, got, err, want)
		}
	}

	for i, h := range hodlers {
		wantBalance(t, h, eth.Ether(0))
		sim.Must(t, "Deposit(%dETH as hodler %d)", deposit, i)(weth.Deposit(sim.WithValueFrom(h, eth.Ether(deposit))))
		wantBalance(t, h, eth.Ether(deposit))
	}

	spenderAddr, _, spender, err := DeployTestableWETH(sim.Acc(deployer), sim)
	if err != nil {
		t.Fatalf("DeployTestableWETH() error %v", err)
	}

	t.Run("without allowance", func(t *testing.T) {
		if diff := revert.WETHAllowance.Diff(spender.TransferFrom(sim.Acc(deployer), sim.Addr(hodler0), sim.Addr(hodler1), eth.Ether(1))); diff != "" {
			t.Errorf("%T.TransferFrom() diff %s", spender, diff)
		}
		wantBalance(t, hodler0, eth.Ether(deposit))
		wantBalance(t, hodler1, eth.Ether(deposit))
	})

	t.Run("with allowance", func(t *testing.T) {
		wantAllowance := func(t *testing.T, owner, approved common.Address, want *big.Int) {
			t.Helper()
			if got, err := weth.Allowance(nil, owner, approved); err != nil || got.Cmp(want) != 0 {
				t.Errorf("weth.Allowance(%v, %v) got %d, err = %v; want %d, nil err", owner, approved, got, err, want)
			}
		}

		const approve = 3
		sim.Must(t, "Approve(spender contract, %dETH)", approve)(weth.Approve(sim.Acc(hodler0), spenderAddr, eth.Ether(approve)))
		wantAllowance(t, sim.Addr(hodler0), spenderAddr, eth.Ether(approve))

		const spend = 1
		sim.Must(t, "TransferFrom(%d) as approved", spend)(spender.TransferFrom(sim.Acc(deployer), sim.Addr(hodler0), sim.Addr(hodler1), eth.Ether(spend)))
		wantAllowance(t, sim.Addr(hodler0), spenderAddr, eth.Ether(approve-spend))

		wantBalance(t, hodler0, eth.Ether(deposit-spend))
		wantBalance(t, hodler1, eth.Ether(deposit+spend))
	})

}
