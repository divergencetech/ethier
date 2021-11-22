// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/sales/FixedPriceSeller.sol";

/// @notice A concrete FixedPriceSeller for testing the cost() function.
contract TestableFixedPriceSeller is FixedPriceSeller {
    constructor(uint256 price)
        FixedPriceSeller(price, SellerConfig(0, 0, 0), payable(0))
    {}

    function _handlePurchase(address, uint256) internal override {}

    function totalSupply() public pure override returns (uint256) {
        return 0;
    }
}
