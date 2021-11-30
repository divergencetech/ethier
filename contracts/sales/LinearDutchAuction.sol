// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "./Seller.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";

/// @notice A Seller with a linearly decreasing price.
abstract contract LinearDutchAuction is Seller {
    /**
    @param unit The unit of "time" used for decreasing prices, block number or
    timestamp. NOTE: See the comment on AuctionIntervalUnit re use of Time as a
    unit.
    @param startPoint The block or timestamp at which the auction opens. A value
    of zero disables the auction. See setAuctionStartPoint().
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

    NOTE: The Block unit is more reliable as it has an explicit progression
    (simply incrementing). Miners are allowed to have a time drift into the
    future although which predisposes to unexpected behaviour by which "future"
    costs are encountered. See the ConsenSys 15-second rule:
    https://consensys.net/blog/developers/solidity-best-practices-for-smart-contract-security/
     */
    enum AuctionIntervalUnit {
        UNSPECIFIED,
        Block,
        Time
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
        require(
            config.startPrice >= config.decreaseSize * config.numDecreases,
            "LinearDutchAuction: negative reserve"
        );
        dutchAuctionConfig = config;
    }

    /**
    @notice Sets the config startPoint. A startPoint of zero disables the
    auction.
    @dev The auction can be toggle on and off with this function, without the
    cost of having to update the entire config.
     */
    function setAuctionStartPoint(uint256 startPoint) public onlyOwner {
        dutchAuctionConfig.startPoint = startPoint;
    }

    /// @notice Override of Seller.cost() with Dutch-auction logic.
    function cost(uint256 n) public view override returns (uint256) {
        DutchAuctionConfig storage cfg = dutchAuctionConfig;

        uint256 current;
        if (cfg.unit == AuctionIntervalUnit.Block) {
            current = block.number;
        } else if (cfg.unit == AuctionIntervalUnit.Time) {
            current = block.timestamp;
        }

        require(
            cfg.startPoint != 0 && current >= cfg.startPoint,
            "LinearDutchAuction: Not started"
        );

        uint256 decreases = Math.min(
            (current - cfg.startPoint) / cfg.decreaseInterval,
            cfg.numDecreases
        );
        return n * (cfg.startPrice - decreases * cfg.decreaseSize);
    }
}
