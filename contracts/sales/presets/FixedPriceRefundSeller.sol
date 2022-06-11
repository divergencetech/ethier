// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../sellers/FixedSupplyRefund.sol";
import "../sellers/FixedPrice.sol";
import "../sellers/SellableCallbacker.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract FixedPriceRefundSeller is
    FixedSupplyRefund,
    FixedPrice,
    SellableCallbacker,
    Ownable
{
    constructor(
        uint64 maxSupply,
        uint256 price,
        ISellable sellable
    ) FixedSupply(maxSupply) FixedPrice(price) SellableCallbacker(sellable) {}

    function _beforePurchase(address to, uint256 num)
        internal
        virtual
        override(FixedSupplyRefund, Seller)
        returns (address, uint256)
    {
        (to, num) = FixedSupplyRefund._beforePurchase(to, num);
        return (to, num);
    }

    function purchase(address to, uint64 num) external {
        _purchase(to, num);
    }
}

// import "../sellers/FixedPrice.sol";
// import "./ArbitraryPriceRefundSeller.sol";

// contract FixedPriceRefundSeller is ArbitraryPriceRefundSeller, FixedPrice {
//     constructor(
//         uint64 maxSupply,
//         uint256 price,
//         ISellable sellable
//     ) ArbitraryPriceRefundSeller(maxSupply, sellable) FixedPrice(price) {}

//     function _beforePurchase(address to, uint256 num)
//         internal
//         virtual
//         override(ArbitraryPriceRefundSeller, Seller)
//         returns (address, uint256)
//     {
//         (to, num) = ArbitraryPriceRefundSeller._beforePurchase(to, num);
//         return (to, num);
//     }
// }
