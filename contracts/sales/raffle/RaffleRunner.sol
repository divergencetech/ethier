// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "./IRaffleRunner.sol";
import "./internal/Raffle.sol";

/**
@dev Do not deploy a Raffle directly; instead inherit from RaffleRunner to
deploy an independent Raffle contract. Similarly to OpenZeppelin escrow
contracts, a Raffle MUST be a separate contract, responsible for its own funds.

TODO: copy comments from Raffle re rationale.
 */
abstract contract RaffleRunner is IRaffleRunner {
    Raffle public immutable raffle;

    constructor(uint256 maxWinners, uint256 entryCost) {
        raffle = new Raffle(maxWinners, entryCost);
        raffle.transferOwnership(msg.sender);
    }

    /**
    @dev Accepts winning numbers, requiring that the sender is the raffle,
    before propagating the parameters to the internal redeemRaffleWins(),
    implemented by the user of this contract.
    */
    function protectedRaffleRedemption(address winner, uint128 num) external {
        require(
            msg.sender == address(raffle),
            "RaffleRunner: only raffle can redeem"
        );
        redeemRaffleWins(winner, num);
    }

    function redeemRaffleWins(address winner, uint128 num) internal virtual;
}
