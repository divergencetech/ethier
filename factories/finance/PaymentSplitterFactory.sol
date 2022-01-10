// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./DelegatedPaymentSplitter.sol";
import "@openzeppelin/contracts/proxy/Clones.sol";

/**
@notice This contract deploys a modified version of OpenZeppelin's
PaymentSplitter that can be used with an EIP-1167 minimal proxy contract, and
exposes functions to deploy such proxies to random or deterministic addresses.
@dev As only a single instance of this contract need exist, the addresses of
deployments will be published on github.com/divergencetech/ethier.

NOTE: there is likely no need to import this contract directly; instead see the
ethier documentation for the deployed factory addresses.
 */
contract PaymentSplitterFactory {
    using Clones for address;

    /**
    @notice The primary instance of the modified PaymentSplitter, delegated to
    by all clones.
     */
    address public immutable implementation;

    constructor() {
        implementation = address(new DelegatedPaymentSplitter());
    }

    /// @notice Emitted when a new PaymentSplitter is deployed.
    event PaymentSplitterDeployed(address clonedPaymentSplitter);

    /// @notice Deploy a minimal contract proxy to a PaymentSplitter.
    function deploy(address[] memory payees, uint256[] memory shares)
        external
        returns (address)
    {
        address clone = implementation.clone();
        _postDeploy(clone, payees, shares);
        return clone;
    }

    /**
    @notice Deploy a minimal contract proxy to a PaymentSplitter, at a
    deterministic address.
    @dev Use predictDeploymentAddress() with the same salt to predit the address
    before calling deployDeterministic(). See OpenZeppelin's proxy/Clones.sol
    for details and caveats, primarily that this will revert if a salt is
    reused.
     */
    function deployDeterministic(
        bytes32 salt,
        address[] memory payees,
        uint256[] memory shares
    ) external returns (address) {
        address clone = implementation.cloneDeterministic(salt);
        _postDeploy(clone, payees, shares);
        return clone;
    }

    /**
    @notice Calls initialize(payees, shares) on the proxy contract and then
    emits an event to log the new address.
     */
    function _postDeploy(
        address clone,
        address[] memory payees,
        uint256[] memory shares
    ) internal {
        DelegatedPaymentSplitter(payable(clone)).initialize(payees, shares);
        emit PaymentSplitterDeployed(clone);
    }

    /**
    @notice Returns the address at which a new PaymentSplitter will be deployed
    if using the same salt as passed to this function.
     */
    function predictDeploymentAddress(bytes32 salt)
        external
        view
        returns (address)
    {
        return implementation.predictDeterministicAddress(salt);
    }
}
