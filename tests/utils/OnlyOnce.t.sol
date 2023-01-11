// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import "forge-std/Test.sol";
import "forge-std/console2.sol";

import {OnlyOnce} from "../../contracts/utils/OnlyOnce.sol";

// solhint-disable no-empty-blocks
contract OnlyOnceConsumer is OnlyOnce {
    function limitedFunction1() external onlyOnce {}

    function limitedFunction2(uint256) external onlyOnce {}

    function limitedFunction3()
        external
        onlyOnceByIdentifier(OnlyOnceConsumer.limitedFunction3.selector)
    {}
}

// solhint-enable no-empty-blocks

contract OnlyOnceTest is Test {
    OnlyOnceConsumer public c;

    function setUp() public {
        c = new OnlyOnceConsumer();
    }

    function _testCannotCallTwice(function() external func) internal {
        func();
        vm.expectRevert(
            abi.encodeWithSelector(
                OnlyOnce.FunctionAlreadyExecuted.selector,
                func.selector
            )
        );
        func();
    }

    function _testCannotCallTwice(
        function(uint256) external func,
        uint256 param1,
        uint256 param2
    ) internal {
        func(param1);
        vm.expectRevert(
            abi.encodeWithSelector(
                OnlyOnce.FunctionAlreadyExecuted.selector,
                func.selector
            )
        );
        func(param2);
    }

    function testCannotCallFunctionsTwice(uint256 param1, uint256 param2)
        public
    {
        // This test also ensures that calling one limited function does not
        // affect the executability of the other.

        _testCannotCallTwice(c.limitedFunction1);
        _testCannotCallTwice(c.limitedFunction2, param1, param2);
        _testCannotCallTwice(c.limitedFunction3);
    }
}
