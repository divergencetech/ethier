// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

/**
These contracts and libraries result in the CHAINID OpCode in various places,
which we use to test the Go SourceMap implementation because it's an otherwise
obscure OpCode that won't result in false positives.
 */

library SourceMapTestLib {
    function id() external view returns (uint256 chainId) {
        assembly {
            chainId := chainid()
        }
    }
}

contract SourceMapTest0 {
    /// @dev Allows functions to be non-view.
    event Noop();

    function id() external returns (uint256 chainId) {
        assembly {
            chainId := chainid()
        }
        emit Noop();
    }

    function idPlusOne() external returns (uint256 chainIdPlusOne) {
        assembly {
            chainIdPlusOne := chainid()
        }
        chainIdPlusOne++;
        emit Noop();
    }

    function fromLib() external returns (uint256) {
        emit Noop();
        return SourceMapTestLib.id();
    }
}

contract SourceMapTest1 {
    event Noop();
    
    function id() external returns (uint256 chainId) {
        assembly {
            chainId := chainid()
        }
        emit Noop();
    }
}
