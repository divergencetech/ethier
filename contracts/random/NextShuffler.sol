// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.9 <0.9.0;

import "./PRNG.sol";

/**
 * @notice Returns the next value in a shuffled list [0,n), amortising the
 * shuffle across all calls to _next(). Can be used for randomly allocating a
 * set of tokens but the caveats in `dev` docs MUST be noted.
 * @dev Although the final shuffle is uniformly random, it is entirely
 * deterministic if the seed to the PRNG.Source is known. This MUST NOT be used
 * for applications that require secure (i.e. can't be manipulated) allocation
 * unless parties who stand to gain from malicious use have no control over nor
 * knowledge of the seed at the time that their transaction results in a call to
 * next().
 */
library NextShuffler {
    using PRNG for PRNG.Source;

    struct State {
        // Number of items already shuffled.
        // This is the equivalent of `i` in the Wikipedia description of the
        // Fisher–Yates algorithm.
        uint256 shuffled;
        // The number of items that will be shuffled.
        uint256 numToShuffle;
        // A sparse representation of the shuffled list [0,n). List items that
        // have been shuffled are stored with their original index as the key
        // and their new index + 1 as their value. Note that mappings with
        // numerical values return 0 for non-existent keys so we MUST increment
        // the new index to differentiate between a default value and a new
        // index of 0. See _get() and _set().
        mapping(uint256 => uint256) permutation;
    }

    /**
     * @notice Initialises a shuffler at a given location in storage.
     */
    function init(State storage state, uint256 numToShuffle) internal {
        state.numToShuffle = numToShuffle;
    }

    /**
     * @notice Returns the current value stored in list index `i`, accounting
     * for all historical shuffling.
     */
    function _get(State storage state, uint256 i)
        private
        view
        returns (uint256)
    {
        uint256 val = state.permutation[i];
        return val == 0 ? i : val - 1;
    }

    /**
     * @notice Sets the list index `i` to `val`, equivalent `arr[i] = val` in a
     * standard Fisher–Yates shuffle.
     */
    function _set(
        State storage state,
        uint256 i,
        uint256 val
    ) private {
        state.permutation[i] = i == val ? 0 : val + 1;
    }

    /**
     * @notice Returns the next value in the shuffle list in O(1) time and
     * memory.
     * @dev NB: See the `dev` documentation of this contract re security (or
     * lack thereof) of deterministic shuffling.
     * @param rand Uniformly distributed random number in
     * [0, state.numToShuffle - state.shuffled)
     */
    function next(State storage state, uint256 rand)
        internal
        returns (uint256)
    {
        uint256 shuffled = state.shuffled;
        require(!finished(state), "NextShuffler: finished");

        unchecked {
            // Cannot overflow if rand is supplied as specified.
            rand += shuffled;
        }

        uint256 chosen = _get(state, rand);

        // Even though a full swap of the elements in the list is not needed for
        // the algoritm to work, we do it anyway because it allows us to restart
        // the shuffling.
        _set(state, rand, _get(state, shuffled));
        _set(state, shuffled, chosen);

        unchecked {
            ++state.shuffled;
        }
        return chosen;
    }

    /**
     * @notice Returns the next value in the shuffle list in O(1) time and
     * memory together with the random number that was used for the drawing.
     * @dev NB: See the `dev` documentation of this contract re security (or
     * lack thereof) of deterministic shuffling.
     * @dev This is intended to be used if the random number drawn from `src`
     * that was used for shuffling needs to be reused for something else, e.g.
     * to thoroughly test the algorithm.
     */
    function nextAndRand(State storage state, PRNG.Source src)
        internal
        returns (uint256 choice, uint256 rand)
    {
        src.readLessThan(state.numToShuffle - state.shuffled);
        choice = next(state, rand);
    }

    /**
     * @notice Returns the next value in the shuffle list in O(1) time and
     * memory.
     * @dev NB: See the `dev` documentation of this contract re security (or
     * lack thereof) of deterministic shuffling.
     */
    function next(State storage state, PRNG.Source src)
        internal
        returns (uint256)
    {
        (uint256 choice, ) = nextAndRand(state, src);
        return choice;
    }

    /**
     * @notice Returns a flag that indicates if the entire list has been
     * shuffled.
     */
    function finished(State storage state) internal view returns (bool) {
        return state.shuffled >= state.numToShuffle;
    }

    /**
     * @notice Restarts the shuffler, such that all elements can be drawn again.
     * @dev Restarting does not clear the internal permutation. Running the
     * shuffle again with same seed after restarting might, therefore,
     * yield different results.
     */
    function restart(State storage state) internal {
        state.shuffled = 0;
    }

    /**
     * @notice Resets the shuffler.
     */
    function reset(State storage state) internal {
        uint256 shuffled = state.shuffled;
        for (uint256 i; i < shuffled; ++i) {
            state.permutation[i] = 0;
        }
        restart(state);
    }
}
