// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/utils/DynamicBuffer.sol";

/**
@notice Exposes functions allowing testing of DynamicBuffer.
 */
contract TestableDynamicBuffer {
    using DynamicBuffer for bytes;

    /**
    @notice Allocates a buffer with a given capacity and safely appends data for
    a given number of times.
     */
    function allocateAndAppendRepeated(
        uint256 capacity,
        string memory data,
        uint256 repetitions
    ) public pure returns (string memory) {
        bytes memory buffer = DynamicBuffer.allocate(capacity);

        for (uint256 idx = 0; idx < repetitions; ++idx) {
            buffer.appendSafe(bytes(data));
        }

        return string(buffer);
    }

    /**
    @notice Allocates a buffer with a given capacity and safely appends data for
    a given number of times encoded as base64.
     */
    function allocateAndAppendRepeatedBase64(
        uint256 capacity,
        bytes memory data,
        uint256 repetitions,
        bool fileSafe,
        bool noPadding
    ) public pure returns (string memory) {
        bytes memory buffer = DynamicBuffer.allocate(capacity);

        for (uint256 idx = 0; idx < repetitions; ++idx) {
            buffer.appendSafeBase64(data, fileSafe, noPadding);
        }

        return string(buffer);
    }
}
