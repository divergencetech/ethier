// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./Seller.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";

/// @notice A Seller with fixed per-item price.
abstract contract FixedPriceSeller is Seller {
    constructor(
        uint256 _price,
        Seller.SellerConfig memory sellerConfig,
        address payable _beneficiary
    ) Seller(sellerConfig, _beneficiary) {
        setPrice(_price);
    }

    /**
    @notice The fixed per-item price.
    @dev Fixed as in not changing with time nor number of items, but not a
    constant.
     */
    uint256 public price;

    /// @notice Sets the per-item price.
    function setPrice(uint256 _price) public onlyOwner {
        price = _price;
    }

    /// @notice Override of Seller.cost() with fixed price.
    function cost(uint256 n) public view override returns (uint256) {
        return n * price;
    }
}
