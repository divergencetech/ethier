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

/**
@notice An abstract contract providing the _purchase() function to:
 - Enforce per-wallet / per-transaction limits
 - Calculate required cost, forwarding to a beneficiary, and refunding extra
 */
abstract contract FixedSupplyRefund is FixedSupply, CappedRefund {
    function _capRequested(address, uint256 requested)
        internal
        virtual
        override
        returns (uint256)
    {
        uint256 remaining = totalInventory() - totalSold();
        return Math.min(requested, remaining);
    }

    function _beforePurchase(address to, uint256 num)
        internal
        virtual
        override(FixedSupply, CappedRefund)
        returns (address, uint256)
    {
        (to, num) = CappedRefund._beforePurchase(to, num);
        (to, num) = FixedSupply._beforePurchase(to, num);
        return (to, num);
    }
}
