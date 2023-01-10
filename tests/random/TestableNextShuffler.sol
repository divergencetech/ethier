// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/random/NextShuffler.sol";
import "../../contracts/random/PRNG.sol";

contract TestableNextShuffler {
    using PRNG for PRNG.Source;
    using NextShuffler for NextShuffler.State;

    NextShuffler.State public state;

    uint256[] public permutation;

    /// @notice Emited on each call to _next() to allow for thorough testing.
    event ShuffledWith(uint256 current, uint256 with);

    constructor(uint256 numToShuffle) {
        state.init(numToShuffle);
    }

    /**
     * @notice An instrumented call to `state.next` for testing.
     */
    function _next(PRNG.Source src) internal returns (uint256) {
        uint256 shuffled = state.shuffled;
        (uint256 choice, uint256 rand) = state.nextWithRand(src);
        emit ShuffledWith(shuffled, shuffled + rand);
        return choice;
    }

    function permute(uint64 seed) external {
        permute(keccak256(abi.encodePacked(seed)));
    }

    function permute(bytes32 seed) public {
        PRNG.Source src = PRNG.newSource(seed);
        for (uint256 i = 0; i < state.numToShuffle; i++) {
            permutation.push(_next(src));
        }
    }

    function reset() public {
        state.reset();
    }
}
