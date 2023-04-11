// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./IPaymentSplitterFactory.sol";

/**
@notice Convenience library for using ethier's PaymentSplitterFactory for cheap
deployment of OpenZeppelin PaymentSplitters via minimal proxy contracts. A
single factory contract is deployed on supported chains, the respective
addresses of which are determined via the chainid() and returned by this
library's instance() function.
 */
library PaymentSplitterDeployer {
    /***
    @notice Returns the ethier PaymentSplitterFactory instance for the current
    chain.
     */
    function instance() internal view returns (IPaymentSplitterFactory) {
        address factory;

        assembly {
            switch chainid()
            case 1 {
                // mainnet
                factory := 0x10AF0996e2e289dC8f582e6510b0d251E061BFB5
            }
            case 5 {
                // Görli
                factory := 0x1353463A30F7A9Eb50A645F4D3950d5d5E41F1F4
            }
            case 1337 {
                // The geth SimulatedBackend iff used with the ethier
                // factorytest package.
                factory := 0xa516d2c64ED7Fe2004A93Bc123854B229F3Bb738
            }
        }

        require(
            factory != address(0),
            "PaymentSplitterFactory: not deployed on current chain"
        );
        return IPaymentSplitterFactory(factory);
    }
}
