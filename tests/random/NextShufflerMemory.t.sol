// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.15;

import "forge-std/Test.sol";

import {NextShufflerMemory} from "../../contracts/random/NextShufflerMemory.sol";
import {CSPRNG} from "../../contracts/random/CSPRNG.sol";

contract NextShufflerMemoryTest is Test {
    using NextShufflerMemory for NextShufflerMemory.State;

    function testShuffleMustContainAllNumbers(uint8 length, bytes32 entropy)
        public
    {
        uint256[] memory shuffled = shuffle(length, entropy);
        assertContainsAll(shuffled, length);
    }

    function testShuffleWithDifferentEntropy(bytes32 entropy1, bytes32 entropy2)
        public
    {
        vm.assume(entropy1 != entropy2);
        uint256 length = 1000;
        assertNeq(shuffle(length, entropy1), shuffle(length, entropy2));
    }

    function testCannotShuffleMoreThanAvailable(uint8 length, bytes32 entropy)
        public
    {
        NextShufflerMemory.State memory state = NextShufflerMemory.allocate(
            length,
            entropy
        );
        for (uint256 i; i < length; ++i) {
            state.next();
        }
        vm.expectRevert("NextShuffler: finished");
        state.next();
    }

    function shuffle(uint256 length, bytes32 entropy)
        public
        pure
        virtual
        returns (uint256[] memory)
    {
        NextShufflerMemory.State memory state = NextShufflerMemory.allocate(
            length,
            entropy
        );
        uint256[] memory ret = new uint256[](length);
        for (uint256 i; i < length; ++i) {
            ret[i] = state.next();
        }
        return ret;
    }

    function assertContainsAll(uint256[] memory shuffled, uint256 length)
        public
    {
        assertEq(shuffled.length, length);
        bool[] memory seen = new bool[](length);
        for (uint256 i; i < length; ++i) {
            uint256 x = shuffled[i];
            assertFalse(seen[x]);
            seen[x] = true;
        }
    }

    function assertNeq(uint256[] memory a, uint256[] memory b) public {
        assertTrue(
            keccak256(abi.encodePacked(a)) != keccak256(abi.encodePacked(b))
        );
    }
}

contract NextShufflerMemoryWithExternalRandomnessTest is
    NextShufflerMemoryTest
{
    using NextShufflerMemory for NextShufflerMemory.State;
    using CSPRNG for CSPRNG.Source;

    function shuffle(uint256 length, bytes32 entropy)
        public
        pure
        override
        returns (uint256[] memory)
    {
        NextShufflerMemory.State memory state = NextShufflerMemory.allocate(
            length,
            entropy
        );
        CSPRNG.Source source = CSPRNG.newSource(entropy);

        uint256[] memory ret = new uint256[](length);
        for (uint256 i; i < length; ++i) {
            ret[i] = state.next(
                source.readLessThan(state.numToShuffle - state.shuffled)
            );
        }
        return ret;
    }
}
