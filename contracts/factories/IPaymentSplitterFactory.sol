// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

interface IPaymentSplitterFactory {
    /// @notice Deploys a minimal contract proxy to a PaymentSplitter.
    function deploy(address[] memory payees, uint256[] memory shares)
        external
        returns (address);

    /**
    @notice Deploys a minimal contract proxy to a PaymentSplitter, at a
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
    ) external returns (address);

    /**
    @notice Returns the address at which a new PaymentSplitter will be deployed
    if using the same salt as passed to this function.
     */
    function predictDeploymentAddress(bytes32 salt)
        external
        view
        returns (address);
}
