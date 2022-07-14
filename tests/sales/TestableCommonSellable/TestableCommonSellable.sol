// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/sales/presets/CommonSellable.sol";
import "./SellerMock.sol";

contract TestableCommonSellable is CommonSellable {
    uint64 public totalSupply;
    mapping(address => uint64) public bought;

    constructor() {
        _addSeller(address(new SellerMock(ISellable(this))));
        _addSeller(address(new SellerMock(ISellable(this))));
    }

    function _handlePurchase(address to, uint64 num) internal override {
        totalSupply += num;
        bought[to] += num;
    }
}
