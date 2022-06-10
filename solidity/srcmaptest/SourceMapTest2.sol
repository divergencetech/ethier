// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

/**
See SourceMapTest.sol for rationale.
 */

contract SourceMapTest2 {
    /// @dev Allows functions to be non-view.
    event Noop();

    function id() external returns (uint256 chainId) {
        assembly {
            chainId := chainid()
        }
        emit Noop();
    }
}
