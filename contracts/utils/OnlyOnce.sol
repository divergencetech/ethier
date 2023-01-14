// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "forge-std/console2.sol";

contract OnlyOnce {
    bytes32 private constant _LEFTMOST_FOUR_BYTES_MASK =
        0xFFFFFFFF00000000000000000000000000000000000000000000000000000000;

    /**
     * @notice Thrown if a function can only be executed once.
     */
    error FunctionAlreadyExecuted(bytes32 selector);

    /**
     * @notice Keeps track of function executions.
     */
    mapping(bytes32 => bool) private _functionAlreadyExecuted;

    /**
     * @notice Ensures that a given identifier has not been marked as already
     * executed and marks it as such.
     */
    function _ensureOnlyOnce(bytes32 identifier) private {
        if (_functionAlreadyExecuted[identifier]) {
            revert FunctionAlreadyExecuted(identifier);
        }
        _functionAlreadyExecuted[identifier] = true;
    }

    /**
     * @notice Ensures that the modified function can only be executed once.
     * @param identifier A generic UNIQUE identifier to mark the wrapped
     * function as used. Typically the EVM function selector of the wrapped
     * function (padded to the right with zeroes).
     */
    modifier onlyOnceByIdentifier(bytes32 identifier) {
        _ensureOnlyOnce(identifier);
        _;
    }

    /**
     * @notice Ensures that the modified function can only be executed once.
     * @dev This modifier MUST only be used on functions that are external (not
     * public nor internal). The modifier uses the function selector of the
     * current calldata context as identifier, which can have unintended
     * side-effects for internally used functions.
     */
    modifier onlyOnce() {
        bytes32 selector;
        assembly {
            selector := and(calldataload(0), _LEFTMOST_FOUR_BYTES_MASK)
        }
        _ensureOnlyOnce(selector);
        _;
    }
}
