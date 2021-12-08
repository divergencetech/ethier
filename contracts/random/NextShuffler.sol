// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.9 <0.9.0;

import "./PRNG.sol";

/**
@notice Returns the next value in a shuffled list [0,n), amortising the shuffle
across all calls to _next(). Can be used for randomly allocating a set of tokens
but the caveats in `dev` docs MUST be noted.
@dev Although the final shuffle is uniformly random, it is entirely
deterministic if the seed to the PRNG.Source is known. This MUST NOT be used for
applications that require secure (i.e. can't be manipulated) allocation unless
parties who stand to gain from malicious use have no control over nor knowledge
of the seed at the time that their transaction results in a call to _next().
 */
contract NextShuffler {
    using PRNG for PRNG.Source;

    /// @notice Total number of elements to shuffle.
    uint256 public immutable NUM_TO_SHUFFLE;

    /// @param numToShuffle Total number of elements to shuffle.
    constructor(uint256 numToShuffle) {
        NUM_TO_SHUFFLE = numToShuffle;
    }

    /**
    @dev Number of items already shuffled; i.e. number of historical calls to
    _next(). This is the equivalent of `i` in the Wikipedia description of the
    Fisher–Yates algorithm.
     */
    uint256 private shuffled;

    /**
    @dev A sparse representation of the shuffled list [0,n). List items that
    have been shuffled are stored with their original index as the key and their
    new index + 1 as their value. Note that mappings with numerical values
    return 0 for non-existent keys so we MUST increment the new index to
    differentiate between a default value and a new index of 0. See _get() and
    _set().
     */
    mapping(uint256 => uint256) private _permutation;

    /**
    @notice Returns the current value stored in list index `i`, accounting for
    all historical shuffling.
     */
    function _get(uint256 i) private view returns (uint256) {
        uint256 val = _permutation[i];
        return val == 0 ? i : val - 1;
    }

    /**
    @notice Sets the list index `i` to `val`, equivalent `arr[i] = val` in a
    standard Fisher–Yates shuffle.
     */
    function _set(uint256 i, uint256 val) private {
        _permutation[i] = i == val ? 0 : val + 1;
    }

    /// @notice Emited on each call to _next() to allow for thorough testing.
    event ShuffledWith(uint256 current, uint256 with);

    /**
    @notice Returns the next value in the shuffle list in O(1) time and memory.
    @dev NB: See the `dev` documentation of this contract re security (or lack
    thereof) of deterministic shuffling.
     */
    function _next(PRNG.Source src) internal returns (uint256) {
        require(shuffled < NUM_TO_SHUFFLE, "NextShuffler: finished");

        uint256 j = src.readLessThan(NUM_TO_SHUFFLE - shuffled) + shuffled;
        emit ShuffledWith(shuffled, j);

        uint256 chosen = _get(j);
        _set(j, _get(shuffled));
        shuffled++;
        return chosen;
    }
}
