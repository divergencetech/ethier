// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../BaseSellable.sol";
import "../../utils/OwnerPausable.sol";

/// @notice A base contract for handling sellable content.
/// @dev Sales can only be made via authorized sellers.
abstract contract CommonSellable is BaseSellable, OwnerPausable {
    function changeSellers(
        address[] calldata rmSellerList,
        address[] calldata addSellerList
    ) external onlyOwner {
        for (uint256 idx = 0; idx < rmSellerList.length; ++idx) {
            _removeSeller(rmSellerList[idx]);
        }
        for (uint256 idx = 0; idx < addSellerList.length; ++idx) {
            _addSeller(addSellerList[idx]);
        }
    }

    // function handlePurchase(address to, uint64 num)
    //     public
    //     payable
    //     virtual
    //     override
    //     whenNotPaused
    // {
    //     handlePurchase(to, num);
    // }
}
