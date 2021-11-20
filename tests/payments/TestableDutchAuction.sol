// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/payments/LinearDutchAuction.sol";

/**
@notice Exposes a buy() function to allow testing of DutchAuction and, by proxy,
PurchaseManager.
@dev Setting the price decrease of the DutchAuction to zero is identical to a
constant PurchaseManager. Creating only a single Testable contract is simpler.
 */
contract TestableDutchAuction is LinearDutchAuction {
    constructor(
        LinearDutchAuction.DutchAuctionConfig memory auctionConfig,
        PurchaseManager.PurchaseConfig memory purchaseConfig,
        address payable beneficiary
    ) LinearDutchAuction(auctionConfig, purchaseConfig, beneficiary) {}

    uint256 private total;

    /**
    @dev Override of PurchaseManager.totalSupply(). Usually this would be
    fulfilled by ERC721Enumerable.
     */
    function totalSupply() public view override returns (uint256) {
        return total;
    }

    /**
    @dev Although this mirrors PurchaseManager.bought, it is used to test the
    _numPurchasing value used to communicate with the modified function. This
    mapping also only counts per tx.origin.
     */
    mapping(address => uint256) public purchased;

    /// @dev Public API for testing of managePurchase().
    function buy(uint256 n) public payable managePurchase(n) {
        total += PurchaseManager._numPurchasing;
        purchased[tx.origin] += _numPurchasing;
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
