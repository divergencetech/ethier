// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../sellers/FixedSupplyRefund.sol";
import "../sellers/FixedPrice.sol";
import "../sellers/SellableCallbacker.sol";
import "../../utils/OwnerPausable.sol";

contract FixedPriceRefundSeller is
    FixedSupplyRefund,
    FixedPrice,
    SellableCallbacker,
    OwnerPausable
{
    struct Config {
        uint64 totalInventory;
        uint64 maxPerTx;
        uint64 maxPerAddress;
        uint256 price;
    }

    constructor(Config memory cfg, ISellable sellable)
        FixedSupplyRefund(cfg.totalInventory, cfg.maxPerTx, cfg.maxPerAddress)
        FixedPrice(cfg.price)
        SellableCallbacker(sellable)
    {} // solhint-disable-line no-empty-blocks

    function setSellerConfig(Config memory cfg) external onlyOwner {
        _setTotalInventory(cfg.totalInventory);
        _setPrice(cfg.price);
        _setTxLimits(cfg.maxPerTx, cfg.maxPerAddress);
    }

    function _beforePurchase(
        address to,
        uint256 num,
        uint256 cost
    )
        internal
        virtual
        override(FixedSupplyRefund, InternalCostSeller)
        returns (
            address,
            uint256,
            uint256
        )
    {
        (to, num, cost) = FixedSupplyRefund._beforePurchase(to, num, cost);
        return (to, num, cost);
    }

    function _afterPurchase(
        address to,
        uint256 num,
        uint256 cost
    ) internal virtual override(FixedSupplyRefund, Seller) {
        CappedRefund._afterPurchase(to, num, cost);
    }

    function purchase(address to, uint256 num) external payable whenNotPaused {
        _purchase(to, num, 0);
    }
}
