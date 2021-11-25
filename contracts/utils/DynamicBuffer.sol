// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)

pragma solidity >=0.8.0;

/// @title DynamicBuffer
/// @author David Huber (@cxkoda) and Simon Fremaux (@dievardump). See also
///         https://raw.githubusercontent.com/dievardump/solidity-dynamic-buffer
/// @notice This library is used to allocate a big amount of container memory
//          which will be subsequently filled without needing to reallocate
///         memory.
/// @dev First, allocate memory.
///      Then use `DynamicBuffer.appendBytes(buffer, theBytes)`
library DynamicBuffer {
    function allocate(uint256 capacity)
        internal
        pure
        returns (bytes memory container, bytes memory buffer)
    {
        assembly {
            // Get next-free memory address
            container := mload(0x40)

            // Allocate memory by setting a new next-free address
            {
                // Add 2 x 32 bytes in size for the two length fields
                // Add 32 bytes safety space for 32B chunked copy
                let size := add(capacity, 0x60)
                let newNextFree := add(container, size)
                mstore(0x40, newNextFree)
            }

            // Set the correct container length
            {
                let length := add(capacity, 0x40)
                mstore(container, length)
            }

            // The buffer starts at idx 1 in the container (0 is length)
            buffer := add(container, 0x20)

            // Init content with length 0
            mstore(buffer, 0)
        }

        return (container, buffer);
    }

    /// @notice Appends data_ to buffer_, and update buffer_ length
    /// @param buffer_ the buffer to append the data to
    /// @param data_ the data to append
    /// @dev Does not perform out-of-bound checks (container capacity)
    ///      for efficiency.
    function appendBytes(bytes memory buffer_, bytes memory data_)
        internal
        pure
    {
        assembly {
            let length := mload(data_)
            for {
                let data := add(data_, 32)
                let dataEnd := add(data, length)
                let buf := add(buffer_, add(mload(buffer_), 32))
            } lt(data, dataEnd) {
                data := add(data, 32)
                buf := add(buf, 32)
            } {
                // Copy 32B chunks from data to buffer.
                // This may read over data array boundaries and copy invalid
                // bytes, which doesn't matter in the end since we will
                // later set the correct buffer length.
                mstore(buf, mload(data))
            }

            // Update buffer length
            mstore(buffer_, add(mload(buffer_), length))
        }
    }

    /// @notice Appends data_ to buffer_, and update buffer_ length
    /// @param buffer_ the buffer to append the data to
    /// @param data_ the data to append
    /// @dev Performs out-of-bound checks and calls `appendBytes`.
    function appendBytesSafe(bytes memory buffer_, bytes memory data_)
        internal
        pure
    {
        uint256 capacity;
        uint256 length;
        assembly {
            capacity := sub(mload(sub(buffer_, 0x20)), 0x20)
            length := mload(buffer_)
        }

        // Account for the 32B chunk copying scheme in `appendBytes`
        uint256 copySize = data_.length;
        uint256 remainder = data_.length % 32;
        if (remainder > 0) {
            copySize += 32 - remainder;
        }

        require(
            length + copySize <= capacity,
            "DynamicBuffer: Appending out of bounds."
        );
        appendBytes(buffer_, data_);
    }
}
