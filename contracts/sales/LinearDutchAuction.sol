// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
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
        uint96 startPrice; // sufficient bits to store up to 1e10 eth
        uint96 decreaseSize;
        uint64 totalInventory;
        uint64 startPoint;
        uint64 decreaseInterval;
        // From https://docs.soliditylang.org/en/v0.8.10/types.html#enums "Enums
        // cannot have more than 256 members"; presumably they take 8 bits, so
        // use some of the numDecreases space instead.
        uint64 numDecreases;
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

    /// @notice Configuration of price changes.
    DutchAuctionConfig public dutchAuctionConfig;

    /// @param expectedReserve See setAuctionConfig().
    constructor(DutchAuctionConfig memory config, uint256 expectedReserve) {
        _setAuctionConfig(config, expectedReserve);
    }

    /**
     * @notice Sets the auction config.
     * @param expectedReserve A safety check that the reserve, as calculated from
     * the config, is as expected.
     */
    function _setAuctionConfig(
        DutchAuctionConfig memory config,
        uint256 expectedReserve
    ) internal {
        require(
            config.decreaseInterval > 0,
            "LinearDutchAuction: zero decrease interval"
        );
        // Underflow might occur is size/num decreases is too large.
        require(
            config.startPrice - config.decreaseSize * config.numDecreases ==
                expectedReserve,
            "LinearDutchAuction: incorrect reserve"
        );
        require(
            config.unit != AuctionIntervalUnit.UNSPECIFIED,
            "LinearDutchAuction: unspecified unit"
        );
        dutchAuctionConfig = config;
    }

    /**
    @notice Sets the config startPoint. A startPoint of zero disables the
    auction.
    @dev The auction can be toggle on and off with this function, without the
    cost of having to update the entire config.
     */
    function _setAuctionStartPoint(uint64 startPoint) internal {
        dutchAuctionConfig.startPoint = startPoint;
    }

    /**
    @notice Override of Seller.cost() with Dutch-auction logic.
    @dev The second parameter, metadata propagated from the call to _purchase(),
    is ignored.
    **/
    function _cost(uint256 num) internal view override returns (uint256) {
        DutchAuctionConfig storage cfg = dutchAuctionConfig;

        uint256 current;
        if (cfg.unit == AuctionIntervalUnit.Block) {
            current = block.number;
        } else if (cfg.unit == AuctionIntervalUnit.Time) {
            // solhint-disable-next-line not-rely-on-time
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
        return num * (cfg.startPrice - decreases * cfg.decreaseSize);
    }
}
