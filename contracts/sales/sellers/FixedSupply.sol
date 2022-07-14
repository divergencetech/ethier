// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/utils/math/Math.sol";
import "./Seller.sol";

/// @notice A seller module to limit the number of purchased tokens based on
/// a maximum total supply.
abstract contract FixedSupply is Seller {
    /// @notice The number of tokens that can be sold by the seller.
    uint64 private _totalInventory;

    /// @notice The number of tokens that have already been sold by the seller.
    uint64 private _totalSold;

    constructor(uint64 totalInventory_) {
        _setTotalInventory(totalInventory_);
    }

    /// @notice Changes the inventory of the seller.
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

    /// @notice Checks if the number of requested purchases is below the limit
    /// given by the inventory.
    /// @dev Reverts otherwise.
    function _beforePurchase(
        address to,
        uint64 num,
        uint256 cost
    )
        internal
        virtual
        override(Seller)
        returns (
            address,
            uint64,
            uint256
        )
    {
        (to, num, cost) = Seller._beforePurchase(to, num, cost);
        require(
            num <= _capOnTotalSupply(num),
            "FixedSupply: To many requested"
        );
        return (to, num, cost);
    }

    /// @notice Updating the total number of sold tokens.
    function _afterPurchase(
        address to,
        uint64 num,
        uint256 cost
    ) internal virtual override(Seller) {
        Seller._afterPurchase(to, num, cost);
        _totalSold += num;
    }

    /// @notice Computes the maximum number of purchases that can be performed in
    /// the current transaction based on the remaining inventory.
    /// @dev This function can be used to dynamically adapt the number of purchased
    /// tokens in case to many are requested.
    function _capOnTotalSupply(uint64 requested)
        internal
        view
        returns (uint64)
    {
        uint64 remaining = _totalInventory - _totalSold;
        require(remaining > 0, "FixedSupply: Sold out");
        return uint64(Math.min(requested, remaining));
    }
}
