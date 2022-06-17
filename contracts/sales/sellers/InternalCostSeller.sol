// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./Seller.sol";

abstract contract InternalCostSeller is Seller {
    /**
    @dev Must return the current cost of a batch of items. This may be constant
    or, for example, decreasing for a Dutch auction or increasing for a bonding
    curve.
    @param num The number of items being purchased.
     */
    function _cost(uint64 num) internal view virtual returns (uint256);

    function cost(uint64 num) external view returns (uint256) {
        return _cost(num);
    }

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

    function _handlePurchase(address to, uint64 num) internal {
        _handlePurchase(to, num, 0);
    }
}
