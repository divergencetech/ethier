// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

/// @notice Basic interface for a contract providing sellable content.
interface ISellable {
    // Todo should we have this return a success flag?
    /// @notice Handles the purchase of the sellable content.
    /// @dev This is usually only callable by sellers.
    function handlePurchase(address to, uint64 num) external payable;
}
