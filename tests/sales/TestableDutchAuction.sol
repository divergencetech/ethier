// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/sales/LinearDutchAuction.sol";

/**
@notice Exposes a buy() function to allow testing of DutchAuction and, by proxy,
Seller.
@dev Setting the price decrease of the DutchAuction to zero is identical to a
constant Seller. Creating only a single Testable contract is simpler.
 */
contract TestableDutchAuction is LinearDutchAuction {
    constructor(
        LinearDutchAuction.DutchAuctionConfig memory auctionConfig,
        Seller.SellerConfig memory sellerConfig,
        address payable beneficiary
    ) LinearDutchAuction(auctionConfig, sellerConfig, beneficiary) {}

    uint256 private total;

    /**
    @dev Override of Seller.totalSupply(). Usually this would be
    fulfilled by ERC721Enumerable.
     */
    function totalSupply() public view override returns (uint256) {
        return total;
    }

    /**
    @dev Although this mirrors Seller.bought, it is used to test the
    _numPurchasing value used to communicate with the modified function. This
    mapping also only counts per tx.origin.
     */
    mapping(address => uint256) public purchased;

    /// @dev Public API for testing of managePurchase().
    function buy(uint256 requested) public payable managePurchase(requested) {
        // The number requested may have been capped by the modifier. The actual
        // amount allowed is communicated via the _getNumPurchasing() method.
        uint256 n = Seller._getNumPurchasing();
        total += n;
        purchased[tx.origin] += n;
    }
}

/// @notice Buys on behalf of a sender to circumvent per-address limits.
contract ProxyPurchaser {
    TestableDutchAuction public auction;

    constructor(address _auction) {
        auction = TestableDutchAuction(_auction);
    }

    function buy(uint256 n) public payable {
        auction.buy(n);
    }
}
