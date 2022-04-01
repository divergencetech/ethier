// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/random/NextShuffler.sol";
import "../../contracts/random/PRNG.sol";

contract TestableNextShuffler is NextShuffler {
    // solhint-disable-next-line no-empty-blocks
    constructor(uint256 numToShuffle) NextShuffler(numToShuffle) {}

    uint256[] public permutation;

    function permute(uint64 seed) external {
        permute(keccak256(abi.encodePacked(seed)));
    }

    function permute(bytes32 seed) public {
        PRNG.Source src = PRNG.newSource(seed);
        for (uint256 i = 0; i < NextShuffler.numToShuffle; i++) {
            permutation.push(NextShuffler._next(src));
        }
    }
}
