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
import "./Seller.sol";

/**
@notice An abstract contract providing the _purchase() function to:
 - Enforce per-wallet / per-transaction limits
 - Calculate required cost, forwarding to a beneficiary, and refunding extra
 */
abstract contract CapRefundSeller is Seller, Ownable {
    using Address for address payable;
    using Monotonic for Monotonic.Increaser;
    using Strings for uint256;

    uint256 internal totalInventory;

    /**
    @notice Tracks total number of items sold by this contract, including those
    purchased free of charge by the contract owner.
     */
    Monotonic.Increaser private _totalSold;

    constructor(uint256 totalInventory_) {
        totalInventory = totalInventory_;
    }

    function setTotalInventory(uint256 totalInventory_) external onlyOwner {
        totalInventory = totalInventory_;
    }

    /// @notice Returns the total number of items sold by this contract.
    function totalSold() public view returns (uint256) {
        return _totalSold.current();
    }

    /**
    @notice Enforces all purchase limits (counts and costs) before calling
    _handlePurchase(), after which the received funds are disbursed to the
    beneficiary, less any required refunds.
    @param to The final recipient of the item(s).
    @param requested The number of items requested for purchase, which MAY be
    reduced when passed to _handlePurchase().
     */
    function _purchase(address to, uint256 requested) internal override {
        requested = _capRequested(to, requested);
        _totalSold.add(requested);
        super._purchase(to, requested);
        _reimburseRest(requested);
    }

    function _capRequested(address to, uint256 requested)
        internal
        virtual
        returns (uint256)
    {
        uint256 remaining = totalInventory - _totalSold.current();
        if (remaining == 0) revert SoldOut();
        return Math.min(requested, remaining);
    }

    // -------------------------------------------------------------------------
    //
    //  Internals
    //
    // -------------------------------------------------------------------------

    function _reimburseRest(uint256 num) private {
        uint256 cost = _cost(num);

        // Ideally we'd be using a PullPayment here, but the user experience is
        // poor when there's a variable cost or the number of items purchased
        // has been capped. We've addressed reentrancy with both a nonReentrant
        // modifier and the checks, effects, interactions pattern.
        if (msg.value > cost) {
            address payable reimburse = payable(_msgSender());
            uint256 refund = msg.value - cost;

            // Using Address.sendValue() here would mask the revertMsg upon
            // reentrancy, but we want to expose it to allow for more precise
            // testing. This otherwise uses the exact same pattern as
            // Address.sendValue().
            (bool success, bytes memory returnData) = reimburse.call{
                value: refund
            }("");
            // Although `returnData` will have a spurious prefix, all we really
            // care about is that it contains the ReentrancyGuard reversion
            // message so we can check in the tests.
            require(success, string(returnData));

            emit Refund(reimburse, refund);
        }
    }

    // -------------------------------------------------------------------------
    //
    //  Events
    //
    // -------------------------------------------------------------------------

    /// @notice Emitted when a buyer is refunded.
    event Refund(address indexed buyer, uint256 amount);

    // -------------------------------------------------------------------------
    //
    //  Errors
    //
    // -------------------------------------------------------------------------

    error SoldOut();
}
