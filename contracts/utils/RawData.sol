// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.16 <0.9.0;

/**
 * @notice Utility library to work with raw bytes data.
 */
library RawData {
    /**
     * @notice Return the byte at the given index interpreted as bool.
     * @dev Any non-zero value is interpreted as true.
     */
    function getBool(bytes memory data, uint256 idx)
        internal
        pure
        returns (bool value)
    {
        return data[idx] != 0;
    }

    /**
     * @notice Clones a bytes array.
     */
    function clone(bytes memory data) internal pure returns (bytes memory) {
        uint256 len = data.length;
        bytes memory buf = new bytes(len);

        uint256 nFullWords = (len - 1) / 32;

        // At the end of data we might still have a few bytes that don't make
        // up a full 32-bytes word.
        // ... [nTailBytes | 32 - nTailBytes -> dirty]
        // So if we again copied a full word for efficiency it would also
        // include some dirty bytes that need to be cleaned first.

        uint256 nTailBytes = len - nFullWords * 32;
        uint256 mask = 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff <<
                ((32 - nTailBytes) * 8);

        assembly {
            let src := add(data, 0x20)
            let dst := add(buf, 0x20)
            for {
                let end := add(src, mul(0x20, nFullWords))
            } lt(src, end) {
                src := add(src, 0x20)
                dst := add(dst, 0x20)
            } {
                mstore(dst, mload(src))
            }

            mstore(dst, and(mload(src), mask))
        }

        return buf;
    }

    /**
     * @notice Reads a big-endian-encoded, 16-bit, unsigned interger from a
     * given offset in a bytes array.
     * @param data The bytes array
     * @param offset The index of the byte in the array at which we start reading.
     * @dev Equivalent to `(uint(data[offset]) << 8) + uint(data[offset + 1])`
     */
    function getUint16(bytes memory data, uint256 offset)
        internal
        pure
        returns (uint16 value)
    {
        assembly {
            value := shr(240, mload(add(data, add(0x20, offset))))
        }
    }

    /**
     * @notice Removes and returns the first byte of an array.
     */
    function popByteFront(bytes memory data)
        internal
        pure
        returns (bytes memory, bytes1)
    {
        bytes1 ret = data[0];
        uint256 len = data.length - 1;
        assembly {
            data := add(data, 1)
            mstore(data, len)
        }
        return (data, ret);
    }

    /**
     * @notice Removes and returns the first DWORD (4bytes) of an array.
     */
    function popDWORDFront(bytes memory data)
        internal
        pure
        returns (bytes memory, bytes4)
    {
        bytes4 ret;
        uint256 len = data.length - 4;
        assembly {
            ret := mload(add(data, 0x20))
            data := add(data, 4)
            mstore(data, len)
        }
        return (data, ret);
    }

    /**
     * @notice Writes an uint32 in little-ending encoding to a given location in
     * bytes array.
     */
    function writeUint32LE(
        bytes memory buf,
        uint256 pos,
        uint32 data
    ) internal pure {
        buf[pos] = bytes1(uint8(data));
        buf[pos + 1] = bytes1(uint8(data >> 8));
        buf[pos + 2] = bytes1(uint8(data >> 16));
        buf[pos + 3] = bytes1(uint8(data >> 24));
    }

    /**
     * @notice Writes an uint16 in little-ending encoding to a given location in
     * bytes array.
     */
    function writeUint16LE(
        bytes memory buf,
        uint256 pos,
        uint16 data
    ) internal pure {
        buf[pos] = bytes1(uint8(data));
        buf[pos + 1] = bytes1(uint8(data >> 8));
    }

    /**
     * @notice Returns a slice of a bytes array.
     * @dev The old array can no longer be used.
     * Intended syntax: `data = data.slice(from, len)`
     */
    function slice(
        bytes memory data,
        uint256 from,
        uint256 len
    ) internal pure returns (bytes memory) {
        assembly {
            data := add(data, from)
            mstore(data, len)
        }
        return data;
    }
}
