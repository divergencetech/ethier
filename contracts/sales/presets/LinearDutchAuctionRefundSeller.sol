// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../sellers/FixedSupplyRefund.sol";
import "../sellers/LinearDutchAuction.sol";
import "../sellers/SellableCallbacker.sol";
import "../../utils/OwnerPausable.sol";

contract LinearDutchAuctionRefundSeller is
    FixedSupplyRefund,
    LinearDutchAuction,
    SellableCallbacker,
    OwnerPausable
{
    struct Config {
        uint64 totalInventory;
        uint64 maxPerTx;
        uint64 maxPerAddress;
    }

    constructor(
        Config memory cfg,
        AuctionConfig memory config,
        uint256 expectedReserve,
        ISellable sellable
    )
        FixedSupplyRefund(cfg.totalInventory, cfg.maxPerTx, cfg.maxPerAddress)
        LinearDutchAuction(config, expectedReserve)
        SellableCallbacker(sellable)
    {}

    function setSellerConfig(Config memory cfg) external onlyOwner {
        _setTotalInventory(cfg.totalInventory);
        _setTxLimits(cfg.maxPerTx, cfg.maxPerAddress);
    }

    function setAuctionConfig(
        AuctionConfig memory config,
        uint256 expectedReserve
    ) external onlyOwner {
        _setAuctionConfig(config, expectedReserve);
    }

    function setAuctionStartPoint(uint64 startPoint) external onlyOwner {
        _setAuctionStartPoint(startPoint);
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
        (to, num, cost) = InternalCostSeller._beforePurchase(to, num, cost);
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
