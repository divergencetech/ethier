// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

abstract contract PurchaseHandler {
    function _handlePurchase(
        address to,
        uint256 num,
        uint256 cost
    ) internal virtual;
}

abstract contract Seller is PurchaseHandler, ReentrancyGuard {
    /**
    @notice Enforces all purchase limits (counts and costs) before calling
    _handlePurchase(), after which the received funds are disbursed to the
    beneficiary, less any required refunds.
    @param to The final recipient of the item(s).
    @param num The number of items requested for purchase, which MAY be
    reduced when passed to _handlePurchase().
     */
    function _purchase(
        address to,
        uint256 num,
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

    // /**
    // @dev Must return the current cost of a batch of items. This may be constant
    // or, for example, decreasing for a Dutch auction or increasing for a bonding
    // curve.
    // @param num The number of items being purchased.
    //  */
    // function _cost(uint256 num) internal view virtual returns (uint256);

    // function purchaseCost(uint256 num) external view returns (uint256) {
    //     return _cost(num);
    // }

    // -------------------------------------------------------------------------
    //
    //  Hooks
    //
    // -------------------------------------------------------------------------

    function _beforePurchase(
        address to,
        uint256 num,
        uint256 cost
    )
        internal
        virtual
        returns (
            address,
            uint256,
            uint256
        )
    {
        return (to, num, cost);
    }

    function _afterPurchase(
        address to,
        uint256 num,
        uint256 cost
    ) internal virtual {} // solhint-disable-line no-empty-blocks

    // -------------------------------------------------------------------------
    //
    //  Errors
    //
    // -------------------------------------------------------------------------

    error InvalidPurchaseValue(uint256 should);
}

abstract contract InternalCostSeller is Seller {
    /**
    @dev Must return the current cost of a batch of items. This may be constant
    or, for example, decreasing for a Dutch auction or increasing for a bonding
    curve.
    @param num The number of items being purchased.
     */
    function _cost(uint256 num) internal view virtual returns (uint256);

    function cost(uint256 num) external view returns (uint256) {
        return _cost(num);
    }

    function _beforePurchase(
        address to,
        uint256 num,
        uint256
    )
        internal
        virtual
        override
        returns (
            address,
            uint256,
            uint256
        )
    {
        return (to, num, _cost(num));
    }
}
