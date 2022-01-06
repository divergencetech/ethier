// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./ProxiedPaymentSplitter.sol";
import "@openzeppelin/contracts/proxy/Clones.sol";

/**
@notice This contract deploys a modified version of OpenZeppelin's
PaymentSplitter that can be used with a minimal proxy contract, and exposes
functions to deploy such proxies to random or deterministic addresses.
@dev As only a single instance of this contract need exist, the addresses of
deployments will be published on github.com/divergencetech/ethier.
 */
contract PaymentSplitterDeployer {
    using Clones for address;

    /**
    @notice The primary instance of the modified PaymentSplitter, delegated to
    by all clones.
     */
    address public immutable primary;

    constructor() {
        primary = address(new ProxiedPaymentSplitter());

        // The primary instance also needs to be initialised even though it will
        // likely never be used.
        address[] memory payees = new address[](1);
        payees[0] = msg.sender;
        uint256[] memory shares = new uint256[](1);
        shares[0] = 1;

        _postDeploy(primary, payees, shares);
    }

    /// @notice Emitted when a new PaymentSplitter is deployed.
    event PaymentSplitterDeployed(address clonedPaymentSplitter);

    /// @notice Deploy a minimal contract proxy to a PaymentSplitter.
    function deploy(address[] memory payees, uint256[] memory shares)
        external
        returns (address)
    {
        address clone = primary.clone();
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
        address clone = primary.cloneDeterministic(salt);
        _postDeploy(clone, payees, shares);
        return clone;
    }

    /**
    @notice Calls init(payees, shares) on the proxy contract and then emits an
    event to log the new address.
     */
    function _postDeploy(
        address clone,
        address[] memory payees,
        uint256[] memory shares
    ) internal {
        ProxiedPaymentSplitter(payable(clone)).init(payees, shares);
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
        return primary.predictDeterministicAddress(salt);
    }
}
