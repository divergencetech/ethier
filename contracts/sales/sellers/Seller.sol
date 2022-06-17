// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/// @notice Interface to handle purchases in `Seller`s.
/// @dev This handles the actual effect of a purchase. This can be anything
/// from minting ERC721 tokens to transfering funds, etc.
abstract contract PurchaseHandler {
    function _handlePurchase(
        address to,
        uint64 num,
        uint256 cost
    ) internal virtual;
}

/// @notice Abstract base contract for all `Seller`s.
/// @dev The intention of this contract is to provide an extensible base
/// for various kinds of Seller modules that can be flexibly composed
/// to build more complex sellers - allowing effective code reuse.
/// can be used like modules to build up a final seller.
/// Derived contracts are intended to implement their logic by overriding
/// and extending the `_beforePurchase` and `_afterPurchase` hooks (calling
/// the parent implementation(s) to compose logic).
/// The former is intended to perform manipulations and checks of the
/// input data; the latter to update the internal state of the module.
/// Final sellers will compose these modules and expose an addition external
/// purchase function for buyers.
abstract contract Seller is PurchaseHandler, ReentrancyGuard {
    /// @notice Internal function handling a give purchase, performing checks
    /// and input manipulations depending on the logic in the hooks.
    /// @param to The receiver of the purchase
    /// @param num Number of requested purchases
    /// @param cost Total cost of the purchase
    /// @dev This function is intended to be wrapped in an external method for
    /// final sellers. Since we cannot foresee what logic will be implemented
    /// in the hooks, we added a reentrancy guard for safety.
    function _purchase(
        address to,
        uint64 num,
        uint256 cost
    ) internal virtual nonReentrant {
        (to, num, cost) = _beforePurchase(to, num, cost);

        // TODO: Check how to parse custom errors in go
        // if (msg.value < cost) revert InvalidPurchaseValue(cost);
        if (msg.value < cost) {
            revert(
                string(
                    abi.encodePacked(
                        "Seller: Costs ",
                        Strings.toString(cost / 1e9),
                        " GWei"
                    )
                )
            );
        }

        _handlePurchase(to, num, cost);
        _afterPurchase(to, num, cost);
        assert(address(this).balance == 0);
    }

    /// @notice Returns hash that identifies the current seller.
    /// @dev This is intended for integrations with other services.
    function sellerType() external virtual returns (bytes32);

    // -------------------------------------------------------------------------
    //
    //  Hooks
    //
    // -------------------------------------------------------------------------

    /// @notice Hook that is called before handling a purchase.
    /// @dev The intent of this hook is to manipulate the input data and perform
    /// checks before actually handling the purchase.
    /// @param to The receiver of the purchase
    /// @param num Number of requested purchases
    /// @param cost Total cost of the purchase
    /// @dev This function MUST return sensible values, since these will be used
    /// to perfom the purchase.
    /// @dev Don't perform updates of the internal state in this hook! The purchase
    /// parameters might still be subject to change in another module.
    function _beforePurchase(
        address to,
        uint64 num,
        uint256 cost
    )
        internal
        virtual
        returns (
            address,
            uint64,
            uint256
        )
    {
        return (to, num, cost);
    }

    /// @notice Hook that is called after handling a purchase.
    /// @dev The intent of this hook is to the internal state of the seller
    /// (module) if necessary. It is critical that the updates happen here and
    /// not in `_beforePurchase` because only after the purchase the input values
    /// can be considered fixed.
    function _afterPurchase(
        address to,
        uint64 num,
        uint256 cost
    ) internal virtual {} // solhint-disable-line no-empty-blocks

    // -------------------------------------------------------------------------
    //
    //  Errors
    //
    // -------------------------------------------------------------------------

    error InvalidPurchaseValue(uint256 should);
}
