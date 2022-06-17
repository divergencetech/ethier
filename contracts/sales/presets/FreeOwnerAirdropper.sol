// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../sellers/FixedSupplyRefund.sol";
import "../sellers/FixedPrice.sol";
import "../sellers/SellableCallbacker.sol";
import "../../utils/OwnerPausable.sol";

contract FreeOwnerAirdropper is FixedSupply, SellableCallbacker, OwnerPausable {
    constructor(uint64 totalInventory, ISellable sellable)
        FixedSupply(totalInventory)
        SellableCallbacker(sellable)
    {} // solhint-disable-line no-empty-blocks

    struct Receiver {
        address to;
        uint64 num;
    }

    function airdrop(Receiver[] calldata receivers)
        external
        payable
        onlyOwner
        whenNotPaused
    {
        for (uint256 idx = 0; idx < receivers.length; ++idx) {
            _purchase(receivers[idx].to, receivers[idx].num, 0);
        }
    }

    bytes32 private constant TYPE =
        keccak256("ETHIER.sellers.FreeOwnerAirdropper");

    function sellerType() external view virtual override returns (bytes32) {
        return TYPE;
    }
}
