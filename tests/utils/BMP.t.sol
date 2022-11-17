// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import "forge-std/Test.sol";
import "forge-std/console2.sol";

import {BMP} from "../../contracts/utils/BMP.sol";
import {TestLib} from "../TestLib.sol";

contract BMPTest is Test {
    using TestLib for Vm;

    function testBlock22() public {
        bytes memory raw = hex"ffffff_ff0000_00ff00_0000ff";
        bytes memory bmp = BMP.bmp(raw, 2, 2);
        assertTrue(vm.isValidBMP(bmp));
    }

    function testBlock33() public {
        bytes
            memory raw = hex"ff0000_00ff00_0000ff_ff0000_00ff00_0000ff_ff0000_00ff00_0000ff";
        bytes memory bmp = BMP.bmp(raw, 3, 3);
        console2.logBytes(bmp);
        assertTrue(vm.isValidBMP(bmp));
    }

    function testBlock44() public {
        bytes
            memory raw = hex"ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000_ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000";
        bytes memory bmp = BMP.bmp(raw, 4, 4);
        assertTrue(vm.isValidBMP(bmp));
    }

    function testBlock88() public {
        bytes
            memory raw = hex"ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000_ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000_ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000_ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000_ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000_ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000_ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000_ff0000_00ff00_0000ff_ffffff_ffffff_0000ff_00ff00_ff0000";
        bytes memory bmp = BMP.bmp(raw, 8, 8);
        assertTrue(vm.isValidBMP(bmp));
    }
}
