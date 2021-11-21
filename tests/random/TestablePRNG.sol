// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/random/PRNG.sol";

/// @notice Testing contract that exposes functionality of the PRNG library.
contract TestablePRNG {
    using PRNG for PRNG.Source;

    /**
    @notice Representation of the internal state of a PRNG.Source.
    @dev This is for testing purposes only and SHOULD NOT be used. The
    representation is never guaranteed to be stable and MUST NOT be treated as
    part of the public API.
     */
    struct State {
        uint256 seed;
        uint256 counter;
        uint256 entropy;
        uint256 remain;
    }

    /**
    @dev Returns n samples of specified number of bits from the given seed. Also
    returns the internal PRNG state (for testing purposes only) after sampling.
     */
    function sample(
        bytes32 seed,
        uint16 bits,
        uint16 n
    ) public pure returns (uint256[] memory, State memory) {
        // NB! Read the documentation of PRNG re unpredictability.
        PRNG.Source src = PRNG.newSource(seed);

        uint256[] memory samples = new uint256[](n);
        for (uint256 i = 0; i < n; i++) {
            samples[i] = src.read(bits);
        }

        State memory state;
        (state.seed, state.counter, state.entropy, state.remain) = src.state();

        return (samples, state);
    }

    /// @dev Exposes PRNG.bitLength().
    function bitLength(uint256 n) public pure returns (uint256) {
        return PRNG.bitLength(n);
    }

    /// @dev Returns n samples in [0,max).
    function readLessThan(
        bytes32 seed,
        uint256 max,
        uint16 n
    ) public pure returns (uint256[] memory) {
        // NB! Read the documentation of PRNG re unpredictability.
        PRNG.Source src = PRNG.newSource(seed);

        // As all samples have the same upper bound, calculate the bit length
        // once and reuse it.
        uint16 bits = PRNG.bitLength(max);

        uint256[] memory samples = new uint256[](n);
        for (uint256 i = 0; i < n; i++) {
            samples[i] = src.readLessThan(max, bits);
        }
        return samples;
    }
}
