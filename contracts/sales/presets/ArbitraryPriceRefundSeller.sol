// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../sellers/FixedSupplyRefund.sol";
import "../sellers/FixedPrice.sol";
import "../sellers/SellableCallbacker.sol";
import "../../utils/OwnerPausable.sol";

abstract contract ArbitraryPriceRefundSeller is
    FixedSupplyRefund,
    SellableCallbacker,
    OwnerPausable
{
    constructor(uint64 maxSupply, ISellable sellable)
        FixedSupply(maxSupply)
        SellableCallbacker(sellable)
    {}

    function purchase(address to, uint64 num) external whenNotPaused {
        _purchase(to, num);
    }

    function _beforePurchase(address to, uint256 num)
        internal
        virtual
        override(FixedSupplyRefund)
        returns (address, uint256)
    {
        (to, num) = FixedSupplyRefund._beforePurchase(to, num);
        return (to, num);
    }
}
