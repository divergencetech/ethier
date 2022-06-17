// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./Seller.sol";

/// @notice Extends the basic seller by assuming that the total cost of the
/// purchase can be computed by an internal function and is not supplied
/// externally.
abstract contract InternalCostSeller is Seller {
    /// @notice Computes the total cost of purchasing `num` tokens.
    /// @dev This is intended to be overridden by derived contracts.
    function _cost(uint64 num) internal view virtual returns (uint256);

    /// @notice Returns the total cost of purchasing `num` tokens.
    /// @dev Intended for third-party integrations.
    function cost(uint64 num) external view returns (uint256) {
        return _cost(num);
    }

    /// @dev Replaces the cost of the purchase with the computed value.
    function _beforePurchase(
        address to,
        uint64 num,
        uint256
    )
        internal
        virtual
        override
        returns (
            address,
            uint64,
            uint256
        )
    {
        return (to, num, _cost(num));
    }

    /// @dev Convenience function without cost that is now computed internally
    /// instead.
    function _handlePurchase(address to, uint64 num) internal {
        _handlePurchase(to, num, 0);
    }
}
