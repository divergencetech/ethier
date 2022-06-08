// Package openseatest provides test doubles for OpenSea's Wyvern protocol.
package openseatest

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/openseatest/openseatestabi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// proxyRegistry is the address at which the simulated ProxyRegistry is deployed
// by DeployProxyRegistry(). This is deterministic because mocked entities in
// ethtest.SimulatedBackend have deterministic keys.
var proxyRegistry = common.HexToAddress("E1a2bbc877b29ADBC56D2659DBcb0ae14ee62071")

// DeployProxyRegistry deploys a mocked Wyvern proxy registry to the
// SimulatedBackend.
//
// This function MUST only be called once for each SimulatedBackend; all future
// calls will deploy to a different address to the one expected by the
// OpenSeaGasFreeListing library.
func DeployProxyRegistry(sim *ethtest.SimulatedBackend) error {
	return sim.AsMockedEntity(ethtest.OpenSea, func(opts *bind.TransactOpts) error {

		addr, _, _, err := openseatestabi.DeploySimulatedProxyRegistry(opts, sim)
		if err != nil {
			return fmt.Errorf("openseatestabi.DeploySimulatedProxyRegistry() error %v", err)
		}

		if !bytes.Equal(addr.Bytes(), proxyRegistry.Bytes()) {
			return fmt.Errorf("unexpected deployment address %v; want %v", addr, proxyRegistry)
		}
		return nil

	})
}

// DeployProxyRegistryTB calls DeployProxyRegistry() and reports any errors with
// tb.Fatal.
func DeployProxyRegistryTB(tb testing.TB, sim *ethtest.SimulatedBackend) {
	tb.Helper()

	if err := DeployProxyRegistry(sim); err != nil {
		tb.Fatalf("openseatest.DeployProxyRegistry() error %v", err)
	}
}

// SetProxy sets the owner's proxy in the simulated Wyver proxy registry, which
// MUST already have been deployed with DeployProxyRegistry[TB}().
func SetProxy(sim *ethtest.SimulatedBackend, owner, proxy common.Address) error {
	reg, err := openseatestabi.NewSimulatedProxyRegistry(proxyRegistry, sim)
	if err != nil {
		return fmt.Errorf("openseatestabi.NewSimulatedProxyRegistry(): %v", err)
	}

	return sim.AsMockedEntity(ethtest.OpenSea, func(opts *bind.TransactOpts) error {
		_, err := reg.SetProxyFor(opts, owner, proxy)
		return err
	})
}

// SetProxyTB calls SetProxy() and reports any errors with tb.Fatal.
func SetProxyTB(tb testing.TB, sim *ethtest.SimulatedBackend, owner, proxy common.Address) {
	tb.Helper()

	if err := SetProxy(sim, owner, proxy); err != nil {
		tb.Fatalf("openseatestabi.SetProxy(%v, %v) error %v", owner, proxy, err)
	}
}
