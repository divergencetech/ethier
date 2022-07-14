// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/sales/presets/ArbitraryPriceRefundSeller.sol";
import "../SellableMock.sol";

/// @notice A concrete FixedPriceSeller for testing the cost() function.
contract TestableArbitraryPriceSeller is ArbitraryPriceRefundSeller {
    constructor(uint64 inventory)
        ArbitraryPriceRefundSeller(
            Config({totalInventory: inventory, maxPerTx: 0, maxPerAddress: 0}),
            new SellableMock()
        )
    {} // solhint-disable-line no-empty-blocks

    /**
    @notice Allow purchasing at any price; exposed only for testing.
    @dev DO NOT USE IN PRODUCTION; the caller MUST NOT be able to control the
    cost of an item.
     */
    function purchase(uint64 num, uint256 costEach)
        external
        payable
        whenNotPaused
    {
        _purchase(msg.sender, num, costEach * num);
    }

    function sellerType() external view virtual override returns (bytes32) {}
}
