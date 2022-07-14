// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "./Seller.sol";

/// @notice Abstract seller module that modifies the requested number of tokens
/// based on a capping function that needs to be implemented. If to much funds
/// were sent with the transaction, the rest will be reimbursed.
abstract contract CappedRefund is Seller, Context {
    function _capRequested(address to, uint64 requested)
        internal
        view
        virtual
        returns (uint64);

    /// @notice Caps the number of purchased tokens.
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
        num = _capRequested(to, num);
        return (to, num, cost);
    }

    /// @notice Reimpurses the surplus to the sender.
    /// @dev It is important that this is handled after the actual purchase
    /// bacause the input parameters might still be subject to change in
    /// `_beforePurchase`.
    function _afterPurchase(
        address,
        uint64,
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
            // solhint-disable-next-line avoid-low-level-calls
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
