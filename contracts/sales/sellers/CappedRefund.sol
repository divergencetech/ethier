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
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "./Seller.sol";

/**
@notice An abstract contract providing the _purchase() function to:
 - Enforce per-wallet / per-transaction limits
 - Calculate required cost, forwarding to a beneficiary, and refunding extra
 */
abstract contract CappedRefund is Seller, Context {
    function _capRequested(address to, uint256 requested)
        internal
        virtual
        returns (uint256);

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
        num = _capRequested(to, num);
        return (to, num, cost);
    }

    function _afterPurchase(
        address,
        uint256,
        uint256 cost
    ) internal virtual override(Seller) {
        _reimburseRest(cost);
    }

    // -------------------------------------------------------------------------
    //
    //  Internals
    //
    // -------------------------------------------------------------------------

    function _reimburseRest(uint256 cost) private {
        // uint256 cost = _cost(num);

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
}
