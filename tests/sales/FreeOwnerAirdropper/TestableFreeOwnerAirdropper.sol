// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/sales/presets/FreeOwnerAirdropper.sol";
import "../SellableMock.sol";

/// @notice A concrete FixedPriceSeller for testing the cost() function.
contract TestableFreeOwnerAirdropper is FreeOwnerAirdropper {
    constructor(uint64 totalInventory)
        FreeOwnerAirdropper(totalInventory, new SellableMock())
    {}
}
