// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/sales/presets/LinearDutchAuctionRefundSeller.sol";
import "../SellableMock.sol";

/// @notice A concrete FixedPriceSeller for testing the cost() function.
contract TestableDutchAuction is LinearDutchAuctionRefundSeller {
    constructor(
        Config memory sellerConfig_,
        AuctionConfig memory auctionConfig_,
        uint256 expectedReserve
    )
        LinearDutchAuctionRefundSeller(
            sellerConfig_,
            auctionConfig_,
            expectedReserve,
            new SellableMock()
        )
    {} // solhint-disable-line no-empty-blocks

    /// @dev Returns the current timestamp for testing of time-based auctions.
    function timestamp() public view returns (uint256) {
        // solhint-disable-next-line not-rely-on-time
        return block.timestamp;
    }

    function sellerConfig() external view returns (Config memory cfg) {
        cfg.totalInventory = totalInventory();
        cfg.maxPerTx = maxPerTx();
        cfg.maxPerAddress = maxPerAddress();
    }

    // This is a convenience interface to maintain compatibility with already
    // implemented tests.
    function own(address owner) external view returns (uint64) {
        return SellableMock(address(sellable)).balanceOf(owner);
    }
}

/// @notice Buys on behalf of a sender to circumvent per-address limits.
contract ProxyPurchaser {
    TestableDutchAuction public auction;

    constructor(address _auction) {
        auction = TestableDutchAuction(_auction);
    }

    function purchase(address to, uint64 n) public payable {
        auction.purchase(to, n);
    }
}

/// @notice A malicious contract that attempts to reenter the buy() function.
/// @dev Naming things is hard. Is Reenterer a word?
contract ReentrantProxyPurchaser {
    TestableDutchAuction public auction;

    constructor(address _auction) {
        auction = TestableDutchAuction(_auction);
    }

    function purchase(address to, uint64 n) public payable {
        auction.purchase{value: msg.value}(to, n);
    }

    receive() external payable {
        // Attempt reentrance when receiving a refund.
        // solhint-disable-next-line avoid-tx-origin
        auction.purchase{value: msg.value}(tx.origin, 1);
    }
}
