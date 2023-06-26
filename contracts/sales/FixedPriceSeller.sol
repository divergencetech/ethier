// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./Seller.sol";

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
    function setPrice(uint256 _price) public onlyRole(DEFAULT_STEERING_ROLE) {
        price = _price;
    }

    /**
    @notice Override of Seller.cost() with fixed price.
    @dev The second parameter, metadata propagated from the call to _purchase(),
    is ignored.
     */
    function cost(uint256 n, uint256) public view override returns (uint256) {
        return n * price;
    }
}
