// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

interface ISellable {
    // Todo should we have this return a success flag?
    function handlePurchase(address to, uint64 num) external payable;
}
