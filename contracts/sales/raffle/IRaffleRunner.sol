// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.9 <0.9.0;

interface IRaffleRunner {
    /**
    @dev Implemented by RaffleRunner to ensure call is from its raffle, after
    which, the parameters are forwarded to an internal, virtual function that
    needs to be implemented by the user of RaffleRunner.
    */
    function protectedRaffleRedemption(address winner, uint128 num) external;
}
