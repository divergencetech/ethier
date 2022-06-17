// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/utils/Address.sol";
import "../ISellable.sol";
import "./Seller.sol";

abstract contract SellableCallbacker is PurchaseHandler {
    ISellable public immutable sellable;

    constructor(ISellable sellable_) {
        sellable = ISellable(sellable_);
    }

    function _handlePurchase(
        address to,
        uint64 num,
        uint256 cost
    ) internal virtual override {
        sellable.handlePurchase{value: cost}(to, num);
    }
}
