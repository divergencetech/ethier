// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../utils/Monotonic.sol";
import "../utils/OwnerPausable.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./CapRefundSeller.sol";

/**
@notice An abstract contract providing the _purchase() function to:
 - Enforce per-wallet / per-transaction limits
 - Calculate required cost, forwarding to a beneficiary, and refunding extra
 */
abstract contract TotalSupplyCapSeller is CapRefundSeller {
    using Address for address payable;
    using Monotonic for Monotonic.Increaser;
    using Strings for uint256;

    uint256 private _totalInventory;

    /**
    @notice Tracks total number of items sold by this contract, including those
    purchased free of charge by the contract owner.
     */
    Monotonic.Increaser private _totalSold;

    constructor(uint256 totalInventory_) {
        _totalInventory = totalInventory_;
    }

    function setTotalInventory(uint256 totalInventory_) internal {
        _totalInventory = totalInventory_;
    }

    /// @notice Returns the total number of items sold by this contract.
    function totalSold() public view returns (uint256) {
        return _totalSold.current();
    }

    function _capRequested(address to, uint256 requested)
        internal
        virtual
        override
        returns (uint256)
    {
        uint256 remaining = _totalInventory - _totalSold.current();
        if (remaining == 0) revert SoldOut();
        return Math.min(requested, remaining);
    }

    function _beforePurchase(
        address,
        uint256 num,
        uint256
    ) internal virtual override {
        _totalSold.add(num);
    }

    // -------------------------------------------------------------------------
    //
    //  Errors
    //
    // -------------------------------------------------------------------------

    error SoldOut();
}
