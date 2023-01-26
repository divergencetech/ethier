// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import {ERC721ACommon} from "../../contracts/erc721/ERC721ACommon.sol";
import {IERC4906Events} from "../../contracts/erc721/ERC4906.sol";
import {Test} from "../TestLib.sol";

contract TestableERC721ACommon is ERC721ACommon {
    constructor(address admin, address steerer)
        ERC721ACommon(admin, steerer, "Test", "T", payable(address(1)), 0)
    {} // solhint-disable-line no-empty-blocks

    function mint(uint256 num) public {
        _mint(msg.sender, num);
    }
}

contract ERC4906Test is Test, IERC4906Events {
    TestableERC721ACommon public token;
    address public admin = makeAddr("admin");
    address public steerer = makeAddr("steerer");

    function setUp() public {
        token = new TestableERC721ACommon(admin, steerer);
    }

    function testEmitBatch(address sender, uint256 numMint) public {
        numMint = bound(numMint, 1, 512);
        token.mint(numMint);

        if (sender != steerer) {
            vm.expectRevert(
                missingRoleError(sender, token.DEFAULT_STEERING_ROLE())
            );
        } else {
            vm.expectEmit(true, true, true, true, address(token));
            emit BatchMetadataUpdate(0, numMint);
        }

        vm.prank(sender);
        token.refreshMetadata();
    }

    function testEmitBatch(uint256 numMint) public {
        // To make sure that the happy path is covered
        testEmitBatch(steerer, numMint);
    }
}
