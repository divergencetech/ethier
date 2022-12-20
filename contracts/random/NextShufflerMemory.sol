// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.9 <0.9.0;

/**
 * @notice Computes the values in a shuffled list [0,n).
 * @dev Although the final shuffle is uniformly random, it is entirely
 * deterministic if the seed to the PRNG.Source is known. This MUST NOT be used
 * for applications that require secure (i.e. can't be manipulated) allocation
 * unless parties who stand to gain from malicious use have no control over nor
 * knowledge of the seed at the time that their transaction results in a call
 * to next().
 * @dev The library is a heavily modified version of `NextShuffler` to allow
 * in-memory shuffling.
 */
library NextShufflerMemory {
    struct State {
        uint256[] permutation;
        bytes32 entropy;
        uint256 shuffled;
        uint256 numToShuffle;
    }

    /// @notice Allocates the internal state of the shuffler.
    function allocate(uint256 numToShuffle, bytes32 entropy)
        internal
        pure
        returns (State memory)
    {
        return
            State({
                permutation: new uint256[](numToShuffle),
                entropy: entropy,
                shuffled: 0,
                numToShuffle: numToShuffle
            });
    }

    /**
     * @notice Returns the current value stored in list index `i`, accounting
     * forall historical shuffling.
     */
    function get(State memory state, uint256 i)
        internal
        pure
        returns (uint256)
    {
        uint256 val = state.permutation[i];
        return val == 0 ? i : val - 1;
    }

    /**
     * @notice Sets the list index `i` to `val`, equivalent `arr[i] = val` in a
     * standard Fisher–Yates shuffle.
     */
    function set(
        State memory state,
        uint256 i,
        uint256 val
    ) internal pure {
        state.permutation[i] = val + 1;
    }

    /**
     * @notice Returns the next value in the shuffle list in O(1) time and
     * memory using mod sampling.
     * @dev NB: See the `dev` documentation of this contract re security (or
     * lack thereof) of deterministic shuffling.
     * @dev Even though random sampling using modulo is biased towards lower
     * values it can safely be neglected if state.numToShuffle << 2^256.
     */
    function next(State memory state) internal pure returns (uint256) {
        require(state.shuffled < state.numToShuffle, "NextShuffler: finished");

        unchecked {
            uint256 rand = _getRandom(state) %
                (state.numToShuffle - state.shuffled);
            return _next(state, rand);
        }
    }

    /**
     * @notice Generates a random number form the seed and number of shuffled
     * elements.
     */
    function _getRandom(State memory state)
        private
        pure
        returns (uint256 rand)
    {
        assembly {
            rand := keccak256(add(state, 0x20), 0x40)
        }
    }

    /**
     * @notice Returns the next value in the shuffle list in O(1) time and
     * memory using a supplied random number.
     * @dev NB: See the `dev` documentation of this contract re security (or
     * lack thereof) of deterministic shuffling.
     * @param rand A uniform random number in
     * [0, state.numToShuffle - state.shuffled)
     */
    function next(State memory state, uint256 rand)
        internal
        pure
        returns (uint256)
    {
        require(state.shuffled < state.numToShuffle, "NextShuffler: finished");
        require(
            rand < state.numToShuffle - state.shuffled,
            "NextShuffler: random number to large"
        );

        return _next(state, rand);
    }

    /**
     * @notice Returns the next value in the shuffle list in O(1) time and
     * memory.
     * @param rand A random number in [0, state.numToShuffle - state.shuffled)
     */
    function _next(State memory state, uint256 rand)
        internal
        pure
        returns (uint256)
    {
        unchecked {
            // Cannot overflow if rand is supplied as specified.
            rand += state.shuffled;
        }
        uint256 chosen = get(state, rand);
        set(state, rand, get(state, state.shuffled));
        unchecked {
            ++state.shuffled;
        }
        return chosen;
    }
}
