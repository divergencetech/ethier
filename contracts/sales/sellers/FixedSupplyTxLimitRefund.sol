// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "./CappedRefund.sol";
import "./FixedSupply.sol";
import "./TxLimit.sol";

/// @notice This composes the functionality of three modules providing a
/// capping of the requested number of items base on a total supply and limits
/// per transaction and purchaser.
abstract contract FixedSupplyTxLimitRefund is
    FixedSupply,
    TxLimit,
    CappedRefund
{
    constructor(
        uint64 totalInventory_,
        uint64 maxPerTx_,
        uint64 maxPerAddress_
    ) FixedSupply(totalInventory_) TxLimit(maxPerTx_, maxPerAddress_) {} // solhint-disable-line no-empty-blocks

    function _beforePurchase(
        address to,
        uint64 num,
        uint256 cost
    )
        internal
        virtual
        override(FixedSupply, TxLimit, CappedRefund)
        returns (
            address,
            uint64,
            uint256
        )
    {
        // Capping based on `_capRequested`
        (to, num, cost) = CappedRefund._beforePurchase(to, num, cost);

        // These calls are redundant because the checks performed in the subroutines
        // will always be true due to the capping. We perform them notheless to avoid
        // bugs if their implementations change at some point. Also the gas overhead
        // is negligible.
        (to, num, cost) = FixedSupply._beforePurchase(to, num, cost);
        (to, num, cost) = TxLimit._beforePurchase(to, num, cost);

        return (to, num, cost);
    }

    /// @dev Update internal states and reimburse the surplus.
    function _afterPurchase(
        address to,
        uint64 num,
        uint256 cost
    ) internal virtual override(FixedSupply, TxLimit, CappedRefund) {
        // Updating internal states
        TxLimit._afterPurchase(to, num, cost);
        FixedSupply._afterPurchase(to, num, cost);

        // Do the refunds
        CappedRefund._afterPurchase(to, num, cost);
    }

    /// @notice Compute the cap based on the tx/address and total supply limit.
    function _capRequested(address to, uint64 requested)
        internal
        view
        virtual
        override
        returns (uint64)
    {
        requested = TxLimit._capOnTxLimit(to, requested);
        requested = FixedSupply._capOnTotalSupply(requested);
        return requested;
    }
}
