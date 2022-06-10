// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./Seller.sol";
import "../utils/OwnerPausable.sol";

/// @notice A Seller with fixed per-item price.
abstract contract FixedPriceSeller is Seller, OwnerPausable {
    /**
    @notice The fixed per-item price.
    @dev Fixed as in not changing with time nor number of items, but not a
    constant.
     */
    uint256 public price;

    constructor(uint256 price_) {
        price = price_;
    }

    /// @notice Sets the per-item price.
    function setPrice(uint256 price_) external onlyOwner {
        price = price_;
    }

    /// @notice Override of Seller.cost() with fixed price.
    function _cost(uint256 num) internal view override returns (uint256) {
        return num * price;
    }
}
