// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/sales/sellers/FixedPrice.sol";
import "../../../contracts/sales/sellers/SellableCallbacker.sol";

contract SellerMock is FixedPrice, SellableCallbacker {
    constructor(ISellable sellable)
        FixedPrice(1 ether)
        SellableCallbacker(sellable)
    {}

    uint64 public numPurchased;

    function purchase(address to, uint64 num) external payable {
        numPurchased += num;
        _handlePurchase(to, num);
    }

    function sellerType() external view virtual override returns (bytes32) {
        return hex"020f";
    }
}
