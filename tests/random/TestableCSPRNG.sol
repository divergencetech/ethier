// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/random/CSPRNG.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/// @notice Testing contract that exposes functionality of the PRNG library.
contract TestableCSPRNG {
    using CSPRNG for CSPRNG.Source;
    using Strings for uint256;

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
    function _sample(
        bytes32 seed,
        uint16 bits,
        uint16 n
    ) private pure returns (uint256[] memory, State memory) {
        // NB! Read the documentation of PRNG re unpredictability.
        CSPRNG.Source src = CSPRNG.newSource(seed);

        uint256[] memory samples = new uint256[](n);
        for (uint256 i = 0; i < n; i++) {
            samples[i] = src.read(bits);
        }

        State memory state;
        (state.seed, state.counter, state.entropy, state.remain) = src.state();

        return (samples, state);
    }

    function sample(
        bytes32 seed,
        uint16 bits,
        uint16 n
    ) public pure returns (uint256[] memory samples) {
        (samples, ) = _sample(seed, bits, n);
    }

    function sampleState(
        bytes32 seed,
        uint16 bits,
        uint16 n
    ) public pure returns (State memory state) {
        (, state) = _sample(seed, bits, n);
    }

    /// @dev Exposes CSPRNG.bitLength().
    function bitLength(uint256 n) public pure returns (uint256) {
        return CSPRNG.bitLength(n);
    }

    /// @dev Returns n samples in [0,max).
    function readLessThan(
        bytes32 seed,
        uint256 max,
        uint16 n
    ) public pure returns (uint256[] memory) {
        // NB! Read the documentation of PRNG re unpredictability.
        CSPRNG.Source src = CSPRNG.newSource(seed);

        // As all samples have the same upper bound, calculate the bit length
        // once and reuse it.
        uint16 bits = CSPRNG.bitLength(max);

        uint256[] memory samples = new uint256[](n);
        for (uint256 i = 0; i < n; i++) {
            samples[i] = src.readLessThan(max, bits);
        }
        return samples;
    }

    uint256[2] public storedSource;

    /**
    @notice Tests store() and loadSource().
    @dev Makes `beforeStore` calls to read(bits) to ensure a non-zero
    state in the Source, stores it to storedSource and immediately reloads it to
    a different Source. Internal state of both copies is also compared before
    100 additional reads are asserted to be identical.
     */
    function testStoreAndLoad(
        bytes32 seed,
        uint16 bits,
        uint16 beforeStore
    ) public {
        CSPRNG.Source src = CSPRNG.newSource(seed);
        for (uint256 i = 0; i < beforeStore; i++) {
            src.read(bits);
        }

        src.store(storedSource);
        CSPRNG.Source copy = CSPRNG.loadSource(storedSource);

        // Confirm that we've actually round-tripped the internal state via
        // storage and not just copied it.
        require(
            CSPRNG.Source.unwrap(src) != CSPRNG.Source.unwrap(copy),
            "Identical Sources, not a copy"
        );

        (
            uint256 seed0,
            uint256 counter0,
            uint256 entropy0,
            uint256 remain0
        ) = src.state();
        (
            uint256 seed1,
            uint256 counter1,
            uint256 entropy1,
            uint256 remain1
        ) = copy.state();

        require(seed0 == seed1, "Seeds differ");
        require(counter0 == counter1, "Counters differ");
        require(remain0 == remain1, "Remaining bits differ");
        // Test the entropy last as it's derived from the other values so a
        // revert() from one of them is more informative.
        require(entropy0 == entropy1, "Entropy differs");

        // Although unnecessary given the check of state, merely being thorough.
        for (uint256 i = 0; i < 100; i++) {
            assert(src.read(bits) == copy.read(bits));
        }
    }
}
