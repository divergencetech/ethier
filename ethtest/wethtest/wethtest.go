// Package openseatest provides a test double for wETH.
package wethtest

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/wethtest/wethtestabi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// wethAddr is the address at which the simulated wETH contract is deployed by
// DeployWETH(). This is deterministic because mocked entities in
// ethtest.SimulatedBackend have deterministic keys.
var wethAddr = common.HexToAddress("0x2336a902f2727C77867A5905dE392fEd3Ff3604b")

// Address returns the address at which DeployWETH deploys the simulated wETH
// contract.
func Address() common.Address {
	return wethAddr
}

// DeployWETH deploys a mocked wETH contract to the SimulatedBackend and returns
// an interface wrapping the address.
//
// This function MUST only be called once for each SimulatedBackend; all future
// calls will deploy to a different address to the one expected by the
// WETH library.
func DeployWETH(sim *ethtest.SimulatedBackend) (*wethtestabi.IwETH, error) {
	err := sim.AsMockedEntity(ethtest.WETH, func(opts *bind.TransactOpts) error {
		addr, _, _, err := wethtestabi.DeployWETH9(opts, sim)
		if err != nil {
			return fmt.Errorf("wethtestabi.DeployWETH9() error %v", err)
		}

		if !bytes.Equal(addr.Bytes(), wethAddr.Bytes()) {
			return fmt.Errorf("unexpected deployment address %v; want %v", addr, wethAddr)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return wethtestabi.NewIwETH(wethAddr, sim)
}

// DeployWETHTB calls DeployWETH() and reports any errors with tb.Fatal.
func DeployWETHTB(tb testing.TB, sim *ethtest.SimulatedBackend) *wethtestabi.IwETH {
	tb.Helper()

	weth, err := DeployWETH(sim)
	if err != nil {
		tb.Fatalf("wethtestabi.DeployWETH9() error %v", err)
	}
	return weth
}
