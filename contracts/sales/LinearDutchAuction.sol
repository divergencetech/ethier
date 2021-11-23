// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "./Seller.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";

/// @notice A Seller with a linearly decreasing price.
abstract contract LinearDutchAuction is Seller {
    /**
    @param unit The unit of "time" used for decreasing prices, block number or
    timestamp.
    @param startPoint The block or timestamp at which the auction opens.
    @param startPrice The price at `startPoint`.
    @param decreaseInterval The number of units to wait before decreasing the
    price. MUST be non-zero.
    @param decreaseSize The amount by which price decreases after every
    `decreaseInterval`.
    @param numDecreases The maximum number of price decreases before remaining
    constant. The reserve price is therefore implicit and equal to
    startPrice-numDecrease*decreaseSize.
     */
    struct DutchAuctionConfig {
        uint256 startPoint;
        uint256 startPrice;
        uint256 decreaseInterval;
        uint256 decreaseSize;
        // From https://docs.soliditylang.org/en/v0.8.10/types.html#enums "Enums
        // cannot have more than 256 members"; presumably they take 8 bits, so
        // use some of the numDecreases space instead.
        uint248 numDecreases;
        AuctionIntervalUnit unit;
    }

    /**
    @notice The unit of "time" along which the cost decreases.
    @dev If no value is provided then the zero UNSPECIFIED will trigger an
    error.

    TODO: implement Time unit. This requires knowledge of how to precisely
    control time in geth's SimulatedBackend also how block timestamps are
    established + the implications thereof on inviariants (e.g. is a timestamp
    guaranteed to be after a transaction is submitted; this seems unlikely
    because it will change the Tx hash, which is known in advance; so many
    questions).
     */
    enum AuctionIntervalUnit {
        UNSPECIFIED,
        Block
    }

    constructor(
        DutchAuctionConfig memory config,
        Seller.SellerConfig memory sellerConfig,
        address payable _beneficiary
    ) Seller(sellerConfig, _beneficiary) {
        setAuctionConfig(config);
    }

    /// @notice Configuration of price changes.
    DutchAuctionConfig public dutchAuctionConfig;

    /// @notice Sets the auction config.
    function setAuctionConfig(DutchAuctionConfig memory config)
        public
        onlyOwner
    {
        require(
            config.unit != AuctionIntervalUnit.UNSPECIFIED,
            "LinearDutchAuction: unspecified unit"
        );
        require(
            config.decreaseInterval > 0,
            "LinearDutchAuction: zero decrease interval"
        );
        dutchAuctionConfig = config;
    }

    /// @notice Override of Seller.cost() with Dutch-auction logic.
    function cost(uint256 n) public view override returns (uint256) {
        DutchAuctionConfig storage cfg = dutchAuctionConfig;

        // TODO: once the Time unit is added, select between block.number and
        // block.timestamp here.
        uint256 current = block.number;

        require(current >= cfg.startPoint, "LinearDutchAuction: Not started");

        uint256 decreases = Math.min(
            (current - cfg.startPoint) / cfg.decreaseInterval,
            cfg.numDecreases
        );
        return n * (cfg.startPrice - decreases * cfg.decreaseSize);
    }
}
