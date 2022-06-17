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

    function _afterPurchase(
        address to,
        uint64 num,
        uint256 cost
    ) internal virtual override(Seller) {
        Seller._afterPurchase(to, num, cost);
        _totalSold += num;
    }

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
