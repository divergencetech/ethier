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
    struct Config {
        uint64 totalInventory;
        uint64 maxPerTx;
        uint64 maxPerAddress;
    }

    constructor(Config memory cfg, ISellable sellable)
        FixedSupplyRefund(cfg.totalInventory, cfg.maxPerTx, cfg.maxPerAddress)
        SellableCallbacker(sellable)
    {} // solhint-disable-line no-empty-blocks

    function setSellerConfig(Config memory cfg) external onlyOwner {
        _setTotalInventory(cfg.totalInventory);
        _setTxLimits(cfg.maxPerTx, cfg.maxPerAddress);
    }

    function _purchase(
        address to,
        uint256 num,
        uint256 cost
    ) internal virtual override whenNotPaused {
        Seller._purchase(to, num, cost);
    }

    function _beforePurchase(
        address to,
        uint256 num,
        uint256 cost
    )
        internal
        virtual
        override(FixedSupplyRefund)
        returns (
            address,
            uint256,
            uint256
        )
    {
        (to, num, cost) = FixedSupplyRefund._beforePurchase(to, num, cost);
        return (to, num, cost);
    }
}
