// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./ISellable.sol";

abstract contract BaseSellable is ISellable {
    mapping(address => bool) internal _sellers;

    function handlePurchase(address to, uint256 num)
        external
        payable
        onlySellers(msg.sender)
    {
        _handlePurchase(to, num);
    }

    function _handlePurchase(address to, uint256 num) internal virtual;

    modifier onlySellers(address caller) {
        require(_sellers[caller], "Unauthorized seller");
        _;
    }

    function setSeller(address seller, bool isAllowed) internal {
        _sellers[seller] = isAllowed;
    }
}
