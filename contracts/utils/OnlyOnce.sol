// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

contract OnlyOnce {
    /**
     * @notice Thrown if a function can only be executed once.
     */
    error FunctionAlreadyExecuted(bytes4 selector);

    /**
     * @notice Keeps track of function executions.
     */
    mapping(bytes4 => bool) private _functionAlreadyExecuted;

    /**
     * @notice Ensures that the modified function can only be executed once.
     * @param selector The selector of the wrapped function.
     */
    modifier onlyOnce(bytes4 selector) {
        if (_functionAlreadyExecuted[selector]) {
            revert FunctionAlreadyExecuted(selector);
        }
        _functionAlreadyExecuted[selector] = true;
        _;
    }
}
