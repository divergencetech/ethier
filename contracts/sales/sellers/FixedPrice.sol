// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./InternalCostSeller.sol";

/// @notice A Seller with a fixed per-item price.
abstract contract FixedPrice is InternalCostSeller {
    /// @notice The per-item price of the seller.
    uint256 private _price;

    constructor(uint256 price_) {
        _price = price_;
    }

    /// @notice Sets the per-item price.
    function _setPrice(uint256 price_) internal {
        _price = price_;
    }

    /// @notice Returns the price per item.
    /// @dev Intended for third-party integration.
    function price() external view returns (uint256) {
        return _price;
    }

    /// @notice Computes the total cost for `num` tokens.
    function _cost(uint64 num)
        internal
        view
        virtual
        override
        returns (uint256)
    {
        return num * _price;
    }
}
