// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

// import "../utils/Monotonic.sol";
// import "../utils/OwnerPausable.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./Seller.sol";

/**
@notice An abstract contract providing the _purchase() function to:
 - Enforce per-wallet / per-transaction limits
 - Calculate required cost, forwarding to a beneficiary, and refunding extra
 */
abstract contract FixedSupply is Seller {
    uint64 private _totalInventory;
    uint64 private _totalSold;

    constructor(uint64 totalInventory_) {
        _setTotalInventory(totalInventory_);
    }

    function _setTotalInventory(uint64 totalInventory_) internal {
        _totalInventory = totalInventory_;
    }

    /// @notice Returns the total number of items sold by this contract.
    function totalSold() public view returns (uint64) {
        return _totalSold;
    }

    /// @notice Returns the total number of items sold by this contract.
    function totalInventory() public view returns (uint64) {
        return _totalInventory;
    }

    // -------------------------------------------------------------------------
    //
    //  Internals
    //
    // -------------------------------------------------------------------------

    function _beforePurchase(
        address to,
        uint256 num,
        uint256 cost
    )
        internal
        virtual
        override(Seller)
        returns (
            address,
            uint256,
            uint256
        )
    {
        (to, num, cost) = Seller._beforePurchase(to, num, cost);
        require(
            _totalSold + num <= _totalInventory,
            "FixedSupply: To many requested"
        );
        _totalSold += uint64(num);
        return (to, num, cost);
    }

    function _capRequested(uint256 requested)
        internal
        virtual
        returns (uint256)
    {
        uint256 remaining = totalInventory() - totalSold();
        require(remaining > 0, "FixedSupply: Sold out");
        return Math.min(requested, remaining);
    }
}
