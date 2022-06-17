// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./ISellable.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/// @notice A base contract for handling sellable content.
/// @dev Sales can only be made via authorized sellers.
abstract contract BaseSellable is ISellable {
    using EnumerableSet for EnumerableSet.AddressSet;

    EnumerableSet.AddressSet private _sellers;

    function handlePurchase(address to, uint64 num) public payable onlySellers {
        _beforePurchase(to, num);
        _handlePurchase(to, num);
        _afterPurchase(to, num);
    }

    function _handlePurchase(address to, uint64 num) internal virtual;

    modifier onlySellers() {
        require(_sellers.contains(msg.sender), "Unauthorized seller");
        _;
    }

    function _addSeller(address seller) internal {
        _sellers.add(seller);
    }

    function _removeSeller(address seller) internal {
        _sellers.remove(seller);
    }

    function getSellers() public view returns (address[] memory sellers_) {
        uint256 len = _sellers.length();
        sellers_ = new address[](len);
        for (uint256 idx = 0; idx < len; ++idx) {
            sellers_[idx] = _sellers.at(idx);
        }
    }

    // -------------------------------------------------------------------------
    //
    //  Hooks
    //
    // -------------------------------------------------------------------------

    function _beforePurchase(address to, uint64 num) internal virtual {} // solhint-disable-line no-empty-blocks

    function _afterPurchase(address to, uint64 num) internal virtual {} // solhint-disable-line no-empty-blocks
}
