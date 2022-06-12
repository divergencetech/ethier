// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./Seller.sol";

/// @notice A Seller with fixed per-item price.
abstract contract FixedPrice is InternalCostSeller {
    /**
    @notice The fixed per-item price.
    @dev Fixed as in not changing with time nor number of items, but not a
    constant.
     */
    uint256 private _price;

    constructor(uint256 price_) {
        _price = price_;
    }

    /// @notice Sets the per-item price.
    function _setPrice(uint256 price_) internal {
        _price = price_;
    }

    function price() external view returns (uint256) {
        return _price;
    }

    /// @notice Override of Seller.cost() with fixed price.
    function _cost(uint256 num)
        internal
        view
        virtual
        override
        returns (uint256)
    {
        return num * _price;
    }
}
