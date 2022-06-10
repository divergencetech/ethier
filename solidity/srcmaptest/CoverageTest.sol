// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/utils/Address.sol";

/**
@dev Arbitrary code with branches to test trace-based coverage with
deterministic behaviour, regardless of inputs.
 */
contract CoverageTest {
    /// @dev Allows functions to be non-view/pure by emitting this event.
    event Noop();

    function foo() external {
        uint256 a = 0;
        uint256 b = 1;
        if (a < b) {
            ++a;
        } else {
            ++b;
            return;
        }
        --a;
        emit Noop();
    }

    function bar() external {
        if (Address.isContract(msg.sender)) {
            emit Noop();
        } else if (block.number == 0) {
            emit Noop();
        } else {
            uint256 x = 0;
            ++x;
            if (x == 1) {
                emit Noop();
            } else {
                --x;
            }

            assembly {
                function foo(y) {
                    y := add(y, 1)
                }

                x := add(x, 1)
                x := add(x, 1)
                if eq(x, 42) {
                    foo(x)
                }
            }
            if (x == 3) {
                --x;
            } else {
                ++x;
            }
        }
    }
}
