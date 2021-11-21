// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/**
@notice An abstract contract providing the _purchase() function to:
 - Enforce per-wallet / per-transaction limits
 - Calculate required cost, forwarding to a beneficiary, and refunding extra
 */
abstract contract Seller is Ownable, ReentrancyGuard {
    using Strings for uint256;

    /**
    @dev Note that the address limits are vulnerable to wallet farming.
    @param alsoLimitOrigin Whether to also limit the maxPerAddress by tx.origin
    if different to msg.sender. This stops minting via a contract to get around
    limits, but doesn't protect against wallet farming.
    */
    struct SellerConfig {
        uint256 totalInventory;
        uint256 maxPerAddress;
        uint256 maxPerTx;
        bool alsoLimitOrigin;
    }

    constructor(SellerConfig memory config, address payable _beneficiary) {
        setSellerConfig(config);
        setBeneficiary(_beneficiary);
    }

    /// @notice Configuration of purchase limits.
    SellerConfig public sellerConfig;

    /// @notice Sets the seller config.
    function setSellerConfig(SellerConfig memory config) public onlyOwner {
        sellerConfig = config;
    }

    /// @notice Recipient of revenues.
    address payable public beneficiary;

    /// @notice Sets the recipient of revenues.
    function setBeneficiary(address payable _beneficiary) public onlyOwner {
        beneficiary = _beneficiary;
    }

    /**
    @dev Must return the current cost of a batch of items. This may be constant
    or, for example, decreasing for a Dutch auction or increasing for a bonding
    curve.
    @param n The number of items being purchased.
     */
    function cost(uint256 n) public view virtual returns (uint256);

    /**
    @dev Must return the number of items already sold. Naming is in keeping with
    the ERC721Enumerable function that provides the expected functionality.
     */
    function totalSupply() public view virtual returns (uint256);

    /**
    @dev Called by _purchase() after all limits have been put in place; must
    perform all contract-specific sale logic, e.g. ERC721 minting.
    @param n The number of items allowed to be purchased, which MAY be less than
    to the number passed to _purchase() but SHALL be greater than zero.
     */
    function _handlePurchase(uint256 n) internal virtual;

    /**
    @notice Tracks the number of items already bought by an address, regardless
    of transferring out (in the case of ERC721).
    @dev This isn't public as it may be skewed due to differences in msg.sender
    and tx.origin, which it treats in the same way such that
    sum(_bought)>=totalSupply().
     */
    mapping(address => uint256) private _bought;

    /**
    @notice Returns min(n, max(extra items addr can purchase)) and reverts if 0.
    @param zeroMsg The message with which to revert on 0 extra.
     */
    function _capExtra(
        uint256 n,
        address addr,
        string memory zeroMsg
    ) internal view returns (uint256) {
        uint256 extra = sellerConfig.maxPerAddress - _bought[addr];
        if (extra == 0) {
            revert(string(abi.encodePacked("Seller: ", zeroMsg)));
        }
        return Math.min(n, extra);
    }

    /// @notice Emitted when a buyer is refunded.
    event Refund(address indexed buyer, uint256 amount);

    /// @notice Emitted on all purchases of non-zero amount.
    event Revenue(
        address indexed beneficiary,
        uint256 numPurchased,
        uint256 amount
    );

    /**
    @notice Enforces all purchase limits (counts and costs) before calling
    _handlePurchase(), after which the received funds are disbursed to the
    beneficiary, less any required refunds.
    @param requested The number of items requested for purchase, which MAY be
    reduced when passed to _handlePurchase().
     */
    function _purchase(uint256 requested) internal nonReentrant {
        /**
         * ##### CHECKS
         */
        uint256 n = Math.min(requested, sellerConfig.maxPerTx);

        n = _capExtra(n, msg.sender, "Sender limit");
        // Enforce the limit even if proxying through a contract.
        if (sellerConfig.alsoLimitOrigin && msg.sender != tx.origin) {
            n = _capExtra(n, tx.origin, "Origin limit");
        }

        n = Math.min(n, sellerConfig.totalInventory - totalSupply());
        require(n > 0, "Sold out");

        uint256 _cost = cost(n);
        if (msg.value < _cost) {
            revert(
                string(
                    abi.encodePacked(
                        "Costs ",
                        (_cost / 1e9).toString(),
                        " GWei"
                    )
                )
            );
        }

        /**
         * ##### EFFECTS
         */
        _bought[msg.sender] += n;
        if (sellerConfig.alsoLimitOrigin && msg.sender != tx.origin) {
            _bought[tx.origin] += n;
        }

        _handlePurchase(n);

        /**
         * ##### INTERACTIONS
         */

        // Ideally we'd be using a PullPayment here, but the user experience is
        // poor when there's a variable cost or the number of items purchased
        // has been capped. We've addressed reentrancy with both a nonReentrant
        // modifier and the checks, effects, interactions pattern.

        if (_cost > 0) {
            beneficiary.transfer(_cost);
            emit Revenue(beneficiary, n, _cost);
        }

        if (msg.value > _cost) {
            address payable reimburse = payable(msg.sender);
            uint256 refund = msg.value - _cost;
            reimburse.transfer(refund);
            emit Refund(reimburse, refund);
        }
    }
}
