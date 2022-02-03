// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/sales/ArbitraryPriceSeller.sol";

/// @notice A concrete ArbitraryPriceSeller for testing the cost() function.
contract TestableArbitraryPriceSeller is ArbitraryPriceSeller {
    constructor(uint256 totalInventory)
        ArbitraryPriceSeller(
            SellerConfig(totalInventory, 0, 0, 0, false, false, false),
            payable(0)
        )
    {}

    /**
    @notice Allow purchasing at any price; exposed only for testing.
    @dev DO NOT USE IN PRODUCTION; the caller MUST NOT be able to control the
    cost of an item.
     */
    function purchase(uint256 n, uint256 costEach) external payable {
        _purchase(msg.sender, n, costEach);
    }

    /**
    @notice Emitted by accidentalFreePurchase() to stop solc from suggesting
    that it be marked as view because we want to test that the transaction
    reverts.
     */
    event Event();

    /**
    @notice The convenience _purchase() function is deliberately disabled to
    avoid accidentally giving something away for free. This will always revert.
     */
    function accidentalFreePurchase() external {
        emit Event();
        _purchase(msg.sender, 1);
    }

    function _handlePurchase(
        address,
        uint256,
        bool
    ) internal override {}
}
