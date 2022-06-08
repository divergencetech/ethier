// Package chainlinktest provides test doubles for Chainlinks's VRF.
package chainlinktest

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
	"github.com/divergencetech/ethier/ethtest/chainlinktest/chainlinktestabi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/google/go-cmp/cmp"
)

// Contracts carries addresses for Chainlink contracts.
type Contracts struct {
	LinkToken, VRFCoordinator common.Address
}

// Addresses returns the addresses to which DeployAll deploys Chainlink
// test-double contracts.
func Addresses() Contracts {
	return addresses
}

var addresses = Contracts{
	LinkToken:      common.HexToAddress("55B04d60213bcfdC383a6411CEff3f759aB366d6"),
	VRFCoordinator: common.HexToAddress("5FfD760b2B48575f3869722cd816d8b3f94DDb48"),
}

// DeployAll deploys all Chainlink test doubles to the backend and returns a
// VRFCoordinator responding to requests. To release resources, Close() must be
// called when the coordinator is no longer needed.
func DeployAll(sim *ethtest.SimulatedBackend) (*VRFCoordinator, error) {
	err := sim.AsMockedEntity(ethtest.Chainlink, func(opts *bind.TransactOpts) error {

		_, _, cl, err := chainlinktestabi.DeploySimulatedChainlink(opts, sim)
		if err != nil {
			return fmt.Errorf("chainlinktestabi.DeploySimulatedChainlink() error %v", err)
		}

		link, err := cl.LinkToken(nil)
		if err != nil {
			return fmt.Errorf("%T.LinkToken(): %v", cl, err)
		}
		vrf, err := cl.VrfCoordinator(nil)
		if err != nil {
			return fmt.Errorf("%T.VrfCoordinator(): %v", cl, err)
		}

		deployed := Contracts{
			LinkToken:      link,
			VRFCoordinator: vrf,
		}
		if want := addresses; !cmp.Equal(deployed, want) {
			return fmt.Errorf("unexpected deployment addresses %+v; expecting %+v", deployed, want)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return newVRFCoordinator(sim)
}

// DeployAllTB calls DeployAll(), reporting any errors on tb.Fatal. There is no
// need to call Close() on the returned coordinator as this is done by a
// tb.Cleanup() function.
func DeployAllTB(tb testing.TB, sim *ethtest.SimulatedBackend) *VRFCoordinator {
	tb.Helper()
	c, err := DeployAll(sim)
	if err != nil {
		tb.Fatalf("chainlinktest.DeployAll() error %v", err)
	}

	tb.Cleanup(func() {
		if err := c.Close(); err != nil {
			tb.Errorf("%T.Close() error %v", c, err)
		}
	})
	return c
}

// A VRFCoordinator fulfils requests for randomness, listening for events
// emitted by the mocked contract. Fulmilment occurs in a new transaction, and
// the WaitFulfilled*() methods can be used to block until this occurs.
type VRFCoordinator struct {
	vrf       *chainlinktestabi.SimulatedVRFCoordinator
	fulfilled chan error
	err       error

	reqs       chan *chainlinktestabi.SimulatedLinkTokenRandomnessRequested
	quit, done chan struct{}
	sub        event.Subscription
}

func newVRFCoordinator(sim *ethtest.SimulatedBackend) (*VRFCoordinator, error) {
	vrf, err := chainlinktestabi.NewSimulatedVRFCoordinator(addresses.VRFCoordinator, sim)
	if err != nil {
		return nil, fmt.Errorf("NewSimulatedVRFCoordinator(%v): %v", addresses.VRFCoordinator, err)
	}

	filter, err := chainlinktestabi.NewSimulatedLinkTokenFilterer(addresses.LinkToken, sim)
	if err != nil {
		return nil, fmt.Errorf("NewSimulatedLinkTokenFilterer(%v): %v", addresses.LinkToken, err)
	}

	reqs := make(chan *chainlinktestabi.SimulatedLinkTokenRandomnessRequested)
	sub, err := filter.WatchRandomnessRequested(nil, reqs, nil)
	if err != nil {
		close(reqs)
		return nil, fmt.Errorf("%T.WatchRandomnessRequested(): %v", filter, err)
	}

	coord := &VRFCoordinator{
		vrf:       vrf,
		fulfilled: make(chan error),
		reqs:      reqs,
		quit:      make(chan struct{}),
		done:      make(chan struct{}),
		sub:       sub,
	}
	go coord.fulfill(sim)

	return coord, nil
}

func (c *VRFCoordinator) fulfill(sim *ethtest.SimulatedBackend) {
	defer func() {
		c.sub.Unsubscribe()
		close(c.reqs)
		close(c.fulfilled)
		close(c.done)
	}()

	for {
		select {
		case <-c.quit:
			return
		case err := <-c.sub.Err():
			if c.err == nil {
				c.err = fmt.Errorf("%T.Err(): %w", c.sub, err)
			}
		case req := <-c.reqs:
			c.fulfilled <- sim.AsMockedEntity(ethtest.Chainlink, func(opts *bind.TransactOpts) error {
				_, err := c.vrf.Fulfill(opts, req.By)
				return err
			})
		}
	}
}

// Close releases all coordinator resources. After the call returns, further
// requests for randomness will not be fulfilled. The first error that occurred
// while the coordinator was listening for requests, if one occurred, will be
// returned.
func (c *VRFCoordinator) Close() error {
	close(c.quit)
	<-c.done
	return c.err
}

// WaitFulfilled blocks until the last request for randomness is fulfilled.
func (c *VRFCoordinator) WaitFulfilled() error {
	return <-c.fulfilled
}

// WaitFulfilledTB calls WaitFulfilled() and reports any errors on tb.Fatal.
func (c *VRFCoordinator) WaitFulfilledTB(tb testing.TB) {
	tb.Helper()
	if err := c.WaitFulfilled(); err != nil {
		tb.Fatalf("%T.WaitFulfilled() error %v", c, err)
	}
}

// Faucet increases the balances of the specified addresses by the respective
// amounts.
func Faucet(sim *ethtest.SimulatedBackend, extraBalances map[common.Address]*big.Int) error {
	link, err := chainlinktestabi.NewSimulatedLinkToken(addresses.LinkToken, sim)
	if err != nil {
		return fmt.Errorf("NewSimulatedLinkToken(%v): %v", addresses.LinkToken, err)
	}

	return sim.AsMockedEntity(ethtest.Chainlink, func(opts *bind.TransactOpts) error {
		for to, amt := range extraBalances {
			if _, err := link.Faucet(opts, to, amt); err != nil {
				return fmt.Errorf("%T.Faucet(%v, %d): %v", link, to, amt, err)
			}
		}
		return nil
	})
}

// FaucetTB calls Faucet() and reports any errors on tb.Fatal().
func FaucetTB(tb testing.TB, sim *ethtest.SimulatedBackend, extraBalances map[common.Address]*big.Int) {
	tb.Helper()
	if err := Faucet(sim, extraBalances); err != nil {
		tb.Fatalf("chainlinktest.Faucet(%v) error %v", extraBalances, err)
	}
}

// Fee returns the required fee, in LINK, for the specified number of requests
// for randomness.
func Fee(requests int64) *big.Int {
	return eth.Ether(2 * requests)
}
