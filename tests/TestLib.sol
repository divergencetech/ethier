// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import {stdJson, Vm} from "forge-std/Components.sol";
import {Test as Test_} from "forge-std/Test.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

library TestLib {
    using stdJson for string;

    function writeFile(
        Vm vm,
        bytes memory data,
        string memory filename
    ) public {
        vm.writeFile(filename, vm.toString(data));
        string[] memory cmds = new string[](2);
        cmds[0] = "convertHexFile";
        cmds[1] = filename;
        vm.ffi(cmds);
    }

    function mktemp(Vm vm, string memory suffix)
        public
        returns (string memory)
    {
        string[] memory cmds = new string[](2);
        cmds[0] = "mktemp";
        cmds[1] = string.concat("--suffix=", suffix);
        string memory filename = string(vm.ffi(cmds));
        return filename;
    }

    function writeTempFile(Vm vm, bytes memory data)
        public
        returns (string memory)
    {
        return writeTempFile(vm, data, "");
    }

    function writeTempFile(
        Vm vm,
        bytes memory data,
        string memory suffix
    ) public returns (string memory) {
        string memory tmp = mktemp(vm, suffix);
        writeFile(vm, data, tmp);
        return tmp;
    }

    function isValidBMP(Vm vm, bytes memory bmp) public returns (bool) {
        string memory filename = writeTempFile(vm, bmp, "");
        string[] memory cmds = new string[](2);
        cmds[0] = "isValidBMP";
        cmds[1] = filename;
        bytes memory re = vm.ffi(cmds);

        // Checking the echoed exit code of the script
        return keccak256(re) == keccak256("0");
    }
}

contract Test is Test_ {
    function missingRoleError(address account, bytes32 role)
        public
        pure
        returns (bytes memory)
    {
        return
            bytes(
                string.concat(
                    "AccessControl: account ",
                    Strings.toHexString(account),
                    " is missing role ",
                    vm.toString(role)
                )
            );
    }
}
