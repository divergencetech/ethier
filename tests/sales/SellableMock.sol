// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/sales/ISellable.sol";

/// @notice A concrete FixedPriceSeller for testing the cost() function.
contract SellableMock is ISellable {
    /// @notice Emitted on all purchases of non-zero amount.
    event Revenue(
        address indexed beneficiary,
        uint256 numPurchased,
        uint256 amount
    );

    function handlePurchase(address, uint256 num) external payable {
        emit Revenue(address(this), num, msg.value);
    }
}
