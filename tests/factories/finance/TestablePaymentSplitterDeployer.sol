// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/factories/PaymentSplitterDeployer.sol";

/**
@notice Test-only contract to expose functionality on the instance of the
PaymentSplitterFactory deployed on the current chain for use in testing.
@dev Production code need not expose the functions if unnecessary and can merely
call PaymentSplitterDeployer.instance().<function>().
 */
contract TestablePaymentSplitterDeployer {
    function deployDeterministic(
        bytes32 salt,
        address[] memory payees,
        uint256[] memory shares
    ) external {
        PaymentSplitterDeployer.instance().deployDeterministic(
            salt,
            payees,
            shares
        );
    }

    function predictDeploymentAddress(bytes32 salt)
        external
        view
        returns (address)
    {
        return
            PaymentSplitterDeployer.instance().predictDeploymentAddress(salt);
    }
}
