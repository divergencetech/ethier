// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./Seller.sol";

/**
@dev The Seller base contract has a convenience function _purchase(to,n) that
calls the standard function as _purchase(to,n,0). This would result in a free
purchase, to the convenience variant is overriden and always reverts with this
error.
 */
error ImplicitFreePurchase();

/**
@notice A Seller with an arbitrary price passed in externally.
 */
abstract contract ArbitraryPriceSeller is Seller {
    constructor(
        Seller.SellerConfig memory sellerConfig,
        address payable _beneficiary
    ) Seller(sellerConfig, _beneficiary) {} // solhint-disable-line no-empty-blocks

    /**
    @notice Block accidental usage of the convenience function that would
    default to a free sale.
     */
    function _purchase(address, uint256) internal pure override {
        revert ImplicitFreePurchase();
    }

    /**
    @notice Override of Seller.cost() with price passed via metadata.
    @return n*costEach;
     */
    function cost(uint256 n, uint256 costEach)
        public
        pure
        override
        returns (uint256)
    {
        return n * costEach;
    }
}
