// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./Shield.sol";
import "@openzeppelin/contracts/proxy/Clones.sol";

/**
@notice This contract deploys an asset Shield  EIP-1167 minimal proxy contract,
and exposes functions to deploy such proxies to random or deterministic
addresses.
@dev As only a single instance of this contract need exist, the addresses of
deployments will be published on github.com/divergencetech/ethier.

NOTE: there is likely no need to import this contract directly; instead see the
ethier documentation for the deployed factory addresses.
 */
contract ShieldFactory {
    using Clones for address;

    /**
    @notice The primary instance of the Shield, delegated to by all clones.
     */
    address public immutable implementation;

    constructor() {
        implementation = address(new Shield());
    }

    /// @notice Emitted when a new Shield is deployed.
    event ShieldDeployed(address indexed creator, Shield clonedShield);

    /// @notice Deploys a minimal contract proxy to a Shield.
    function deploy() external returns (address) {
        address clone = implementation.clone();
        _postDeploy(Shield(clone), msg.sender);
        return clone;
    }

    /**
    @notice Deploys a minimal contract proxy to a Shield, at a deterministic
    address.
    @dev Use predictDeploymentAddress() with the same salt to predit the address
    before calling deployDeterministic(). See OpenZeppelin's proxy/Clones.sol
    for details and caveats, primarily that this will revert if a salt is
    reused.
     */
    function deployDeterministic(bytes32 salt) external returns (address) {
        address clone = implementation.cloneDeterministic(salt);
        _postDeploy(Shield(clone), msg.sender);
        return clone;
    }

    /**
    @notice Calls initialize(controller) on the proxy contract and then emits an
    event to log the new address.
     */
    function _postDeploy(Shield clone, address creator) internal {
        clone.initialize(creator);
        emit ShieldDeployed(creator, clone);
    }

    /**
    @notice Returns the address at which a new Shield will be deployed if using
    the same salt as passed to this function.
     */
    function predictDeploymentAddress(bytes32 salt)
        external
        view
        returns (address)
    {
        return implementation.predictDeterministicAddress(salt);
    }
}
