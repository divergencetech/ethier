// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

// import "../utils/OwnerPausable.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./CappedRefund.sol";
import "./FixedSupply.sol";
import "./TxLimit.sol";

/**
@notice An abstract contract providing the _purchase() function to:
 - Enforce per-wallet / per-transaction limits
 - Calculate required cost, forwarding to a beneficiary, and refunding extra
 */
abstract contract FixedSupplyRefund is FixedSupply, TxLimit, CappedRefund {
    constructor(
        uint64 totalInventory_,
        uint64 maxPerTx_,
        uint64 maxPerAddress_
    ) FixedSupply(totalInventory_) TxLimit(maxPerTx_, maxPerAddress_) {} // solhint-disable-line no-empty-blocks

    function _beforePurchase(
        address to,
        uint256 num,
        uint256 cost
    )
        internal
        virtual
        override(FixedSupply, CappedRefund)
        returns (
            address,
            uint256,
            uint256
        )
    {
        // Capping based on `_capRequested`
        (to, num, cost) = CappedRefund._beforePurchase(to, num, cost);

        // Updating internal states
        (to, num, cost) = FixedSupply._beforePurchase(to, num, cost);
        return (to, num, cost);
    }

    function _capRequested(address to, uint256 requested)
        internal
        virtual
        override(CappedRefund, TxLimit)
        returns (uint256)
    {
        requested = TxLimit._capRequested(to, requested);
        requested = FixedSupply._capRequested(requested);
        return requested;
    }

    function _afterPurchase(
        address to,
        uint256 num,
        uint256 cost
    ) internal virtual override(CappedRefund, Seller) {
        CappedRefund._afterPurchase(to, num, cost);
    }
}
