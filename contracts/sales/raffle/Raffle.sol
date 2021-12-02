// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "./IRaffleRunner.sol";
import "../../random/PRNG.sol";
import "../../utils/OwnerPausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Address.sol";

/**
@notice Provides a raffle mechanism with additional guaranteed pre-reservation.
This is well-suited to use as an NFT mint pass as entry and refunds require
minimal gas while redemption can be performed at a time when the gas price is
lower.
@dev Entries are fungible, removing the need to track tickets held by a
particular address.

Intended usage: like OpenZeppelin's escrow contracts, this contract should be a
standalone contract that only interacts with the RaffleRunner that instantiated
it. This guarantees that all Ether will be handled according to the `Raffle`
rules, and there is no need to check for payable functions or transfers in the
inheritance tree. 
 */
contract Raffle is OwnerPausable, ReentrancyGuard {
    using Address for address payable;
    using PRNG for PRNG.Source;

    /**
    @notice The RaffleRunner that deployed this Raffle. Used to redeem wins.
     */
    IRaffleRunner public immutable runner;

    /**
    @notice The maximum allowable number of winners, including pre-reserved
    places.
     */
    uint256 public immutable maxWinners;

    /// @notice Cost to reserve a guaranteed spot or to enter the raffle.
    uint256 public immutable entryCost;

    /**
    @param maxWinners_ Maximum number of possible raffle "wins" where reserved,
    i.e. guaruanteed, wins are also counted.
     */
    constructor(uint256 maxWinners_, uint256 entryCost_) {
        runner = IRaffleRunner(msg.sender);
        maxWinners = maxWinners_;
        entryCost = entryCost_;
    }

    /**
    @notice Tracks the number of times a particular address has entered / won.
    @dev Raffle tickets are fungible so there's no need to track specific
    instances; this also makes pre-reservation as simple as simultaneously
    incrementing a Entrant's entries and wins by the number of pre-reserved
    tickets.
     */
    struct Entrant {
        uint128 entries;
        uint128 wins;
    }

    /// @notice Record of entrants.
    mapping(address => Entrant) public entrants;

    /**
    @notice Total number of winners thus far.
    @dev Incremented in line with every Entrant.wins in the entrants map, but
    not decremented upon redemption. If it costs to enter the raffle then the
    beneficiaries MUST NOT receive more than _totalWinners * cost payment at any
    time.
     */
    uint256 private totalWinners;

    /// @notice Requires payment for the specified number of entries.
    modifier costs(uint128 num) {
        require(msg.value == num * entryCost, "Raffle: wrong payment");
        _;
    }

    /**
    @notice Requires that the raffle is in a state to accept new entrants.
     */
    modifier canEnter() {
        require(!Pausable.paused(), "Raffle: closed");
        requireEntropyNotSet();
        _;
    }

    /**
    @notice Reserves the specified number of guaranteed winning tickets for the
    entrant.
    @dev MUST be called from the RaffleRunner, which is responsible for all
    logic determining who may reserve.
     */
    function reserve(address entrant, uint128 num)
        external
        payable
        canEnter
        costs(num)
    {
        require(
            msg.sender == address(runner),
            "Raffle: only runner can reserve"
        );
        uint256 total = totalWinners + num;
        require(total <= maxWinners, "Raffle: too many reserved");

        totalWinners = total;
        Entrant storage e = entrants[entrant];
        e.entries += num;
        e.wins += num;
    }

    /**
    @notice Together, these store a bi-directional mapping from address to
    uint32 identifier, with O(1) lookup.
    @dev Using these IDs instead of raw address in the _entries array reduces
    gas for 5 entries by 19% and for choosing 1024 winners from 5000 entries by
    25%. Reduction to uint16 has no extra gain for entry and minimal further
    gain for shuffling, but limits to 65k participating addresses. See _idFor().
     */
    address[] private addresses;
    mapping(address => uint32) private ids;

    /**
    @notice Returns the already-assigned or new ID for the address.
    @dev Mappings will always return a default value even if a key doesn't
    exist. We therefore store the index from the _ids array, incremented by one
    to differentiate between the first entrant and new ones.
     */
    function idFor(address entrant) private returns (uint32) {
        uint32 id = ids[entrant];
        if (id != 0) {
            return id - 1;
        }

        addresses.push(entrant);
        ids[entrant] = uint32(addresses.length);
        return uint32(addresses.length) - 1;
    }

    /// @notice Raffle tickets from which random winners will be selected.
    uint32[] private entries;

    /// @notice Enter the entrant the specified number of times.
    function enter(address entrant, uint128 num)
        external
        payable
        canEnter
        costs(num)
    {
        entrants[entrant].entries += num;
        uint32 id = idFor(entrant);
        for (uint128 i = 0; i < num; i++) {
            entries.push(id);
        }
    }

    /**
    @notice Entropy source for randomising raffle-ticket shuffling.
    @dev This MUST only be set after entries are closed, to deny entrants the
    ability to cheat. The source of the entropy is left to the user of this
    contract.
     */
    bytes32 public entropy;

    /// @dev Requires that `entropy` == 0.
    function requireEntropyNotSet() private view {
        require(uint256(entropy) == 0, "Raffle: entropy set");
    }

    /**
    @notice Sets the randomisation entropy source.
    @dev The whenPaused modifier avoids a race condition with calls to enter().
     */
    function setEntropy(bytes32 entropy_) external onlyOwner whenPaused {
        requireEntropyNotSet();
        entropy = entropy_;
    }

    /**
    @notice Shuffle the entrants and assign the winners.
    @dev It's safe to allow anyone to run this function as the entropy is set
    securely and shuffling is therefore deterministic. Partial shuffling is
    explicitly forbidden as this allows for manipulation since we don't store
    the state of the PRNG between calls. For example, a malicious caller may
    calculate that resetting the state after n steps will result in their entry
    being chosen at step n+1.
     */
    function shuffle() external {
        require(uint256(entropy) != 0, "Raffle: entropy not set");
        uint256 available = maxWinners - totalWinners;
        require(available > 0, "Raffle: all allocated");

        // j is the notation used in the Wikipedia entry for the Fisher–Yates
        // shuffle.
        uint256 j = entries.length;
        if (j < available) {
            assignWinners(j);
            return;
        }

        // To ensure unbiased, uniform selection of winners we MUST NOT use
        // modulus and instead utilise rejection sampling. The most efficient
        // way to do this is with the least number of required bits sampled on
        // each attempt. This results in an expected number of samples of
        // 1.5*available.
        uint16 bits = PRNG.bitLength(j);
        PRNG.Source src = PRNG.newSource(entropy);

        uint256 winner;
        uint32 swapTmp;
        for (uint256 i = 0; i < available; i++) {
            // Rejection sampling from the last j entrants.
            for (winner = j; winner >= j; winner = src.read(bits)) {}
            winner += i; // Don't swap out existing winners
            swapTmp = entries[i];
            entries[i] = entries[winner];
            entries[winner] = swapTmp;

            if (j & (j - 1) == 0) {
                // j is a power of 2
                bits--;
            }
            j--;
        }

        assignWinners(available);
    }

    /// @notice Flag indicating that shuffle() has assigned winners.
    bool public winnersAssigned;

    /// @notice Increments winning entrants' win count.
    function assignWinners(uint256 n) private {
        for (uint256 i = 0; i < n; i++) {
            entrants[addresses[entries[i]]].wins++;
        }
        totalWinners += n;
        winnersAssigned = true;
    }

    /// @notice Emitted when an entrant is refunded.
    event Refund(address indexed entrant, uint256 amount);

    /// @notice Emitted when an entrant redeems wins, regardless of how many.
    event Redemption(address indexed entrant, uint128 wins);

    /// @notice Requires that winners have been assigned via shuffling.
    modifier whenWinnersAssigned() {
        require(winnersAssigned, "Raffle: winners not assigned");
        _;
    }

    /**
    @notice Redeems wins for and/or reimburses the entrant.
    @dev The message sender is ignored so it's safe to allow calls on behalf of
    an entrant. This allows entrants with zero balance to more easily receive
    their refunds without requiring an additional transfer.
     */
    function redeem(address entrant) external whenWinnersAssigned nonReentrant {
        Entrant storage e = entrants[entrant];

        uint128 wins = e.wins;
        uint256 refund = uint256(e.entries - e.wins) * entryCost;

        /**
         * ##### EFFECTS
         */
        e.wins = 0;
        e.entries = 0;

        // Although this is an effect, it might be an interaction depending on
        // how the user of RaffleRunner implements redeemRaffleWins(), so place
        // it at the end of effects. This is also nonReentrant but we should
        // still employ a defence-in-depth approach.
        if (wins > 0) {
            runner.protectedRaffleRedemption(entrant, wins);
        }
        emit Redemption(entrant, wins);

        /**
         * ##### INTERACTIONS
         */
        if (refund > 0) {
            address payable reimburse = payable(entrant);

            // TODO: abstract this logic from Seller.sol and use here instead
            // of copy-paste.
            (bool success, bytes memory returnData) = reimburse.call{
                value: refund
            }("");
            require(success, string(returnData));

            emit Refund(reimburse, refund);
        }
    }

    /// @notice Emitted when excess places are purchased.
    event ExcessPurchase(address to, uint128 num);

    /**
    @notice Allows purchase of remaining places should totalWinners fall short
    of maxWinners.
     */
    function purchaseExcess(address to, uint128 num)
        external
        payable
        whenWinnersAssigned
        costs(num)
    {
        require(
            totalWinners + num <= maxWinners,
            "Raffle: insufficient excess"
        );

        totalWinners += num;
        emit ExcessPurchase(to, num);

        runner.protectedRaffleRedemption(to, num);
    }

    /// @notice Total already sent by withdraw().
    uint256 public withdrawn;

    /// @notice Emitted by withdraw when funds are sent.
    event Withdrawal(address indexed to, uint256 amount);

    /**
    @notice Sends all remaining revenues of winning entries to the provided
    address.
    @dev Equivalent to withdraw(to, 2^256-1).
     */
    function withdrawAll(address payable to) external {
        withdraw(to, ~uint256(0));
    }

    /**
    @notice Sends partial revenues of winning entries to the provided address.
     */
    function withdraw(address payable to, uint256 max)
        public
        onlyOwner
        whenWinnersAssigned
        nonReentrant
    {
        /**
         * ##### CHECKS
         */
        uint256 send = totalWinners * entryCost - withdrawn;
        require(send > 0, "Raffle: all withdrawn");

        if (send > max) {
            send = max;
        }

        /**
         * ##### EFFECTS
         */

        withdrawn += send;
        emit Withdrawal(to, send);

        /**
         * ##### INTERACTIONS
         */
        to.sendValue(send);
    }
}
