// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import "forge-std/Test.sol";
import "forge-std/console2.sol";

import {OnlyOnce} from "../../contracts/utils/OnlyOnce.sol";

// solhint-disable no-empty-blocks
contract OnlyOnceConsumer is OnlyOnce {
    bytes32 public constant SHARED_IDENTIFIER = keccak256("shared");

    function autoLimited1() external onlyOnce {}

    function autoLimited2(uint256) external onlyOnce {}

    function explicitlyLimited1()
        external
        onlyOnceByIdentifier(OnlyOnceConsumer.explicitlyLimited1.selector)
    {}

    function explicitlyLimited2()
        external
        onlyOnceByIdentifier(OnlyOnceConsumer.explicitlyLimited2.selector)
    {}

    function sharedLimited1()
        external
        onlyOnceByIdentifier(SHARED_IDENTIFIER)
    {}

    function sharedLimited2()
        external
        onlyOnceByIdentifier(SHARED_IDENTIFIER)
    {}
}

// solhint-enable no-empty-blocks

contract OnlyOnceTest is Test {
    OnlyOnceConsumer public c;

    function setUp() public {
        c = new OnlyOnceConsumer();
    }

    function _expectRevertWithAlreadyExecuted(bytes32 identifier) internal {
        vm.expectRevert(
            abi.encodeWithSelector(
                OnlyOnce.FunctionAlreadyExecuted.selector,
                identifier
            )
        );
    }

    function _testCannotCallTwice(function() external func) internal {
        func();
        _expectRevertWithAlreadyExecuted(func.selector);
        func();
    }

    function _testCannotCallTwice(
        function(uint256) external func,
        uint256 param1,
        uint256 param2
    ) internal {
        func(param1);
        _expectRevertWithAlreadyExecuted(func.selector);
        func(param2);
    }

    function testCannotCallFunctionsTwice(uint256 param1, uint256 param2)
        public
    {
        // This test also ensures that calling one limited function does not
        // affect the executability of the other.

        _testCannotCallTwice(c.autoLimited1);
        _testCannotCallTwice(c.autoLimited2, param1, param2);
        _testCannotCallTwice(c.explicitlyLimited1);
        _testCannotCallTwice(c.explicitlyLimited2);
    }

    function testSharedIdentifier() public {
        c.sharedLimited1();
        _expectRevertWithAlreadyExecuted(c.SHARED_IDENTIFIER());
        c.sharedLimited2();
    }
}
