// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/sales/FixedPriceSeller.sol";

/// @notice A concrete FixedPriceSeller for testing the cost() function.
contract TestableFixedPriceSeller is FixedPriceSeller {
    constructor(uint256 price)
        FixedPriceSeller(
            price,
            SellerConfig(0, 0, 0, 0, false, false, false),
            payable(0)
        )
    {} // solhint-disable-line no-empty-blocks

    function _handlePurchase(
        address,
        uint256,
        bool
    ) internal override {} // solhint-disable-line no-empty-blocks
}
