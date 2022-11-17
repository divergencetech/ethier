// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import "forge-std/Test.sol";
import "forge-std/console2.sol";

import {RawData} from "../../contracts/utils/RawData.sol";

contract RawDataLibTest is Test {
    using RawData for bytes;

    function testGetBool(
        bytes memory data,
        uint256 loc,
        bool val
    ) public {
        vm.assume(data.length > 0);
        loc = bound(loc, 0, data.length - 1);
        data[loc] = bytes1(val ? 0x01 : 0x00);

        assertEq(data.getBool(loc), val);
    }

    function testGetUint16Fuzzed(
        bytes memory data,
        uint256 loc,
        uint16 val
    ) public {
        vm.assume(data.length > 1);

        loc = bound(loc, 0, data.length - 2);
        data[loc] = bytes1(uint8(val >> 8));
        data[loc + 1] = bytes1(uint8(val));

        assertEq(data.getUint16(loc), val);
    }

    function testGetUint16() public {
        bytes memory data = hex"12abcdef";
        assertEq(data.getUint16(1), 0xabcd);
    }

    function testWriteUint32LE() public {
        bytes memory data = hex"0123456789abcdef";
        uint32 val = 0xdeadface;
        data.writeUint32LE(1, val);
        assertEq(data, hex"01cefaaddeabcdef");
    }

    function testWriteUint16LE() public {
        bytes memory data = hex"0123456789abcdef";
        uint16 val = 0xdead;
        data.writeUint16LE(2, val);
        assertEq(data, hex"0123adde89abcdef");
    }

    function testPopDWORD() public {
        bytes memory data = hex"0123456789abcdef";
        bytes4 dword;
        (data, dword) = data.popDWORDFront();
        assertEq(uint32(dword), 0x01234567);
        assertEq(data, hex"89abcdef");
    }

    function testSlice() public {
        bytes memory data = hex"0123456789abcdef";
        data = data.slice(2, 3);
        assertEq(data, hex"456789");
    }

    function testClone() public {
        bytes
            memory data = hex"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcd";
        bytes memory clone = data.clone();
        assertEq(clone, data);

        uint256 ptr1;
        uint256 ptr2;
        assembly {
            ptr1 := data
            ptr2 := clone
        }

        assertTrue(ptr1 != ptr2);
    }
}
