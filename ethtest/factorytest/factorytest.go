// The factorytest package provides test doubles for ethier factory contracts.
package factorytest

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/factorytest/factorytestabi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// splitterFactory is the expected deployment address for the
// PaymentSplitterFactory.
var splitterFactory = common.HexToAddress("a516d2c64ED7Fe2004A93Bc123854B229F3Bb738")

// DeployPaymentSplitterFactory deploys a test-double of the ethier
// PaymentSplitterFactory, compatible with the respective ethier library for
// chain-agnostic deployments.
func DeployPaymentSplitterFactory(sim *ethtest.SimulatedBackend) error {
	return sim.AsMockedEntity(ethtest.Ethier, func(opts *bind.TransactOpts) error {
		addr, _, _, err := factorytestabi.DeployPaymentSplitterFactory(opts, sim)

		if !bytes.Equal(addr.Bytes(), splitterFactory.Bytes()) {
			return fmt.Errorf("unexpected deployment address %v; want %v", addr, splitterFactory)
		}

		return err
	})
}

// DeployPaymentSplitterFactoryTB calls DeployPaymentSplitterFactory(),
// reporting any error on tb.Fatal.
func DeployPaymentSplitterFactoryTB(tb testing.TB, sim *ethtest.SimulatedBackend) {
	tb.Helper()
	if err := DeployPaymentSplitterFactory(sim); err != nil {
		tb.Fatalf("DeployPaymentSplitterFactory() error %v", err)
	}
}
