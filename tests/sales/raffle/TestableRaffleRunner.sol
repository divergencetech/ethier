// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/sales/raffle/RaffleRunner.sol";

/**
@dev Similarly to OpenZeppelin escrow contracts, a Raffle MUST be a separate
contract, responsible for its own funds.
 */
contract TestableRaffleRunner is RaffleRunner {
    constructor(uint256 maxWinners, uint256 entryCost)
        RaffleRunner(maxWinners, entryCost)
    {}

    function reserve(uint128 num) public payable {
        /**
        NOTE: in an actual deployment, this function MUST be protected with
        whatever logic is appropriate for defining who can reserve a guaranteed
        entry.

        e.g.
        import "@divergencetech/ethier/contracts/crypto/SignatureChecker.sol";
        …
        using SignatureChecker for EnumerableSet.AddressSet;
        …
        signers.validateSignature(msg.sender, sig);
         */

        /**
        TEST CODE ONLY. An exposed reserve() function MUST implement logic to
        determine whether the sender is allowed to reserve a ticket gauranteed
        to win.
         */
        raffle.reserve{value: msg.value}(msg.sender, num);
    }

    mapping(address => uint128) public wins;

    function redeemRaffleWins(address winner, uint128 num) internal override {
        // e.g. Mint an ERC721
        wins[winner] += num;
    }
}
