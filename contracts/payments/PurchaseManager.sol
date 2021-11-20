// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/**
@notice An abstract contract providing the managePurchase() modifier to:
 - Enforce per-wallet / per-transaction limits
 - Calculate required cost, forwarding to a beneficiary, and refunding extra
 */
abstract contract PurchaseManager is Ownable {
    using Strings for uint256;

    /**
    @dev Note that the address limits are vulnerable to wallet farming.
    @param alsoLimitOrigin Whether to also limit the maxPerAddress by tx.origin
    if different to msg.sender. This stops minting via a contract to get around
    limits, but doesn't protect against wallet farming.
    */
    struct PurchaseConfig {
        uint256 totalInventory;
        uint256 maxPerAddress;
        uint256 maxPerTx;
        bool alsoLimitOrigin;
    }

    constructor(PurchaseConfig memory config, address payable _beneficiary) {
        setPurchaseConfig(config);
        setBeneficiary(_beneficiary);
    }

    /// @notice Configuration of purchase limits.
    PurchaseConfig public purchaseConfig;

    /// @notice Sets the purchase config.
    function setPurchaseConfig(PurchaseConfig memory config) public onlyOwner {
        purchaseConfig = config;
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
    @notice Tracks the number of items already bought by an address, regardless
    of transferring out (in the case of ERC721).
     */
    mapping(address => uint256) public bought;

    /**
    @notice Returns min(n, max(extra items addr can purchase)) and reverts if 0.
    @param zeroMsg The message with which to revert on 0 extra.
     */
    function _capExtra(
        uint256 n,
        address addr,
        string memory zeroMsg
    ) internal view returns (uint256) {
        uint256 extra = purchaseConfig.maxPerAddress - bought[addr];
        require(extra > 0, zeroMsg);
        return Math.min(n, extra);
    }

    /**
    @dev The managePurchase() modifier may adjust the number of items being
    purchased due to per-address limits or to avoid inventory being sold out. To
    communicate this to the modified function, _numPurchasing is set to the
    adjusted number.

    Although this approach isn't ideal because of the additional gas of writing
    and then reading the value, it greatly improves usability of the
    managePurchase() modifier whilst also enforcing the checks, effects,
    interactions pattern.
     */
    uint256 internal _numPurchasing;

    /// @notice Emitted when a buyer is refunded.
    event Refund(address indexed buyer, uint256 amount);

    /// @notice Emitted on all purchases.
    event Revenue(
        address indexed beneficiary,
        uint256 numPurchased,
        uint256 amount
    );

    /**
    @notice Enforces all purchase limits (counts and costs) before executing the
    modified function. After the function is run, the message sender is
    reimbursed for any excess payment.
    @dev This uses the checks, effects, interactions pattern but the SHOULD
    ideally also be modified as nonReentrant.
    @param n The number of items being purchased.
     */
    modifier managePurchase(uint256 n) {
        /**
         * ##### CHECKS
         */
        n = Math.min(n, purchaseConfig.maxPerTx);

        n = _capExtra(n, msg.sender, "Sender limit");
        // Enforce the limit even if proxying through a contract.
        if (purchaseConfig.alsoLimitOrigin && msg.sender != tx.origin) {
            n = _capExtra(n, tx.origin, "Origin limit");
        }

        n = Math.min(n, purchaseConfig.totalInventory - totalSupply());
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
        bought[msg.sender] += n;
        if (purchaseConfig.alsoLimitOrigin && msg.sender != tx.origin) {
            bought[tx.origin] += n;
        }

        _numPurchasing = n;
        _;

        /**
         * ##### INTERACTIONS
         */

        // Ideally we'd be using a PullPayment here, but the user experience is
        // poor when there's a variable cost or the number of items purchased
        // has been capped. We've addressed reentrancy with checks, effects,
        // interactions and also noted in the @dev comment that functions SHOULD
        // also be marked as nonReentrant.

        beneficiary.transfer(_cost);
        emit Revenue(beneficiary, n, _cost);

        if (msg.value > _cost) {
            address payable reimburse = payable(msg.sender);
            uint256 refund = msg.value - _cost;
            reimburse.transfer(refund);
            emit Refund(reimburse, refund);
        }
    }
}
