// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import "forge-std/Test.sol";
import "forge-std/console2.sol";

import {OnlyOnce} from "../../contracts/utils/OnlyOnce.sol";

contract OnlyOnceConsumer is OnlyOnce {
    function limitedFunction1()
        external
        onlyOnce(OnlyOnceConsumer.limitedFunction1.selector)
    {}

    function limitedFunction2()
        external
        onlyOnce(OnlyOnceConsumer.limitedFunction2.selector)
    {}
}

contract OnlyOnceTest is Test {
    OnlyOnceConsumer c;

    function setUp() public {
        c = new OnlyOnceConsumer();
    }

    function testCannotCallFunctionsTwice() public {
        // This test also ensures that calling one limited function does not
        // affect the executability of the other.

        function() external[2] memory funcs = [
            c.limitedFunction1,
            c.limitedFunction2
        ];

        for (uint256 i; i < funcs.length; ++i) {
            funcs[i]();
            vm.expectRevert(
                abi.encodeWithSelector(
                    OnlyOnce.FunctionAlreadyExecuted.selector,
                    OnlyOnceConsumer.limitedFunction1.selector
                )
            );
            funcs[i]();
        }
    }
}
