// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import {Test} from "../TestLib.sol";

import {IERC721A} from "erc721a/contracts/IERC721A.sol";
import {ERC721ACommon} from "../../contracts/erc721/ERC721ACommon.sol";
import {ERC721ATransferRestricted, ERC721ATransferRestrictedBase, TransferRestriction} from "../../contracts/erc721/ERC721ATransferRestricted.sol";

contract TestableERC721ATransferRestricted is ERC721ATransferRestricted {
    constructor()
        ERC721ACommon(
            msg.sender,
            msg.sender,
            "",
            "",
            payable(address(0xdead)),
            0
        )
    {}

    function mint(address to, uint256 amount) public {
        _mint(to, amount);
    }

    function burn(uint256 tokenId) public {
        _burn(tokenId);
    }

    function bypassedTransferFrom(
        address from,
        address to,
        uint256 tokenId
    ) public bypassTransferRestriction {
        transferFrom(from, to, tokenId);
    }
}

contract ERC721ATransferRestrictedGeneralTest is Test {
    TestableERC721ATransferRestricted public token;

    function setUp() public virtual {
        token = new TestableERC721ATransferRestricted();
    }

    function testVandalCannotCallOwnerFunctions(address vandal, address steerer)
        public
    {
        vm.assume(steerer != address(0));
        vm.assume(vandal != address(this));
        vm.assume(vandal != steerer);
        token.grantRole(token.DEFAULT_STEERING_ROLE(), steerer);

        vm.startPrank(vandal, steerer);

        vm.expectRevert(
            missingRoleError(vandal, token.DEFAULT_STEERING_ROLE())
        );
        token.setTransferRestriction(TransferRestriction.None);

        vm.expectRevert(
            missingRoleError(vandal, token.DEFAULT_STEERING_ROLE())
        );
        token.lockTransferRestriction(TransferRestriction.None);
    }

    function toTransferRestriction(uint256 x)
        public
        view
        returns (TransferRestriction)
    {
        return
            TransferRestriction(
                bound(
                    x,
                    uint256(type(TransferRestriction).min),
                    uint256(type(TransferRestriction).max)
                )
            );
    }

    function _checkAndSetup(address alice, address bob) internal {
        vm.assume(alice != bob);
        vm.assume(alice != address(0));
        vm.assume(bob != address(0));

        token.mint(alice, 1);
    }

    function testApprovalClearedAfterOnlyBurn(address alice, address bob)
        public
    {
        _checkAndSetup(alice, bob);

        vm.prank(alice);
        token.setApprovalForAll(bob, true);

        _setRestrictionAndAssertApprovedForAll(
            TransferRestriction.None,
            alice,
            bob,
            true
        );

        _setRestrictionAndAssertApprovedForAll(
            TransferRestriction.OnlyBurn,
            alice,
            bob,
            false
        );
        _setRestrictionAndAssertApprovedForAll(
            TransferRestriction.None,
            alice,
            bob,
            true
        );

        _setRestrictionAndAssertApprovedForAll(
            TransferRestriction.Frozen,
            alice,
            bob,
            false
        );
        _setRestrictionAndAssertApprovedForAll(
            TransferRestriction.None,
            alice,
            bob,
            true
        );
    }

    function _setRestrictionAndAssertApprovedForAll(
        TransferRestriction restriction,
        address owner,
        address operator,
        bool expectedApproval
    ) internal {
        token.setTransferRestriction(restriction);
        assertEq(token.isApprovedForAll(owner, operator), expectedApproval);
    }

    function testRestrictionGetter() public {
        TransferRestriction[3] memory restrictons = [
            TransferRestriction.None,
            TransferRestriction.OnlyBurn,
            TransferRestriction.Frozen
        ];

        for (uint256 i; i < restrictons.length; ++i) {
            token.setTransferRestriction(restrictons[i]);
            assertEq(uint8(token.transferRestriction()), uint8(restrictons[i]));
        }
    }

    function testLock(uint8 setRestriction, uint8 lockRestriction) public {
        TransferRestriction sr = toTransferRestriction(setRestriction);
        TransferRestriction lr = toTransferRestriction(lockRestriction);

        token.setTransferRestriction(sr);

        if (sr != lr) {
            vm.expectRevert(
                abi.encodeWithSelector(
                    ERC721ATransferRestricted
                        .TransferRestrictionCheckFailed
                        .selector,
                    sr
                )
            );
        }
        token.lockTransferRestriction(lr);

        if (sr == lr) {
            vm.expectRevert(
                abi.encodeWithSelector(
                    ERC721ATransferRestricted.TransferRestrictionLocked.selector
                )
            );
        }
        token.setTransferRestriction(sr);
    }
}

contract TransferBehaviourTest is Test {
    TestableERC721ATransferRestricted public token;

    bytes internal _lockedErr =
        abi.encodeWithSelector(
            ERC721ATransferRestrictedBase
                .DisallowedByTransferRestriction
                .selector
        );

    bytes internal _notApprovedErr =
        abi.encodeWithSelector(
            IERC721A.TransferCallerNotOwnerNorApproved.selector
        );

    struct TestCase {
        TransferRestriction restriction;
        bool wantTransfersLocked;
        bool wantMintLocked;
        bool wantBurnLocked;
    }

    TestCase public tt;

    constructor(TestCase memory testCase_) {
        tt = testCase_;
    }

    function setUp() public virtual {
        token = new TestableERC721ATransferRestricted();
    }

    function _checkAndSetup(
        address alice,
        address bob,
        uint8 tokenId
    ) internal {
        vm.assume(alice != bob);
        vm.assume(alice != address(0));
        vm.assume(bob != address(0));

        token.mint(alice, uint256(tokenId) + 1);
        token.setTransferRestriction(tt.restriction);
    }

    function testOwnerTransfer(
        address alice,
        address bob,
        uint8 tokenId
    ) public {
        _checkAndSetup(alice, bob, tokenId);
        bool locked = tt.wantTransfersLocked;

        if (locked) {
            vm.expectRevert(_lockedErr);
        }
        vm.prank(alice);
        token.transferFrom(alice, bob, 0);

        assertEq(
            token.balanceOf(alice),
            token.totalSupply() - (locked ? 0 : 1)
        );
        assertEq(token.balanceOf(bob), locked ? 0 : 1);
    }

    function testApprovedTransfer(
        address alice,
        address bob,
        uint8 tokenId
    ) public {
        _checkAndSetup(alice, bob, tokenId);
        bool locked = tt.wantTransfersLocked;

        if (locked) {
            vm.expectRevert(_lockedErr);
        }
        vm.prank(alice);
        token.approve(bob, 0);

        if (locked) {
            vm.expectRevert(_notApprovedErr);
        }
        vm.prank(bob);
        token.transferFrom(alice, bob, 0);

        assertEq(
            token.balanceOf(alice),
            token.totalSupply() - (locked ? 0 : 1)
        );
        assertEq(token.balanceOf(bob), locked ? 0 : 1);
    }

    function testApprovedForAllTransfer(
        address alice,
        address bob,
        uint8 tokenId
    ) public {
        _checkAndSetup(alice, bob, tokenId);
        bool locked = tt.wantTransfersLocked;

        if (locked) {
            vm.expectRevert(_lockedErr);
        }
        vm.prank(alice);
        token.setApprovalForAll(bob, true);

        if (locked) {
            vm.expectRevert(_notApprovedErr);
        }
        vm.prank(bob);
        token.transferFrom(alice, bob, 0);

        assertEq(
            token.balanceOf(alice),
            token.totalSupply() - (locked ? 0 : 1)
        );
        assertEq(token.balanceOf(bob), locked ? 0 : 1);
    }

    function testMint(address alice, uint8 num) public {
        vm.assume(alice != address(0));
        vm.assume(num > 0);
        token.setTransferRestriction(tt.restriction);

        uint256 totalSupply = token.totalSupply();

        bool locked = tt.wantMintLocked;
        if (locked) {
            vm.expectRevert(_lockedErr);
        }
        token.mint(alice, num);

        assertEq(token.balanceOf(alice), totalSupply + (locked ? 0 : num));
    }

    function testBurn(
        address alice,
        address bob,
        uint8 tokenId
    ) public {
        _checkAndSetup(alice, bob, tokenId);
        uint256 totalSupply = token.totalSupply();
        bool locked = tt.wantBurnLocked;

        if (locked) {
            vm.expectRevert(_lockedErr);
        }
        token.burn(tokenId);

        assertEq(token.balanceOf(alice), totalSupply - (locked ? 0 : 1));
    }

    function testBypassedTransfer(
        address alice,
        address bob,
        uint8 tokenId
    ) public {
        _checkAndSetup(alice, bob, tokenId);

        vm.prank(alice);
        token.bypassedTransferFrom(alice, bob, 0);

        assertEq(token.balanceOf(alice), token.totalSupply() - 1);
        assertEq(token.balanceOf(bob), 1);

        // To make sure the bypass restores the right settings afterwards.
        if (tt.wantTransfersLocked) {
            vm.expectRevert(_lockedErr);
        }
        vm.prank(bob);
        token.transferFrom(bob, alice, 0);
    }
}

contract NoneRestrictionTest is TransferBehaviourTest {
    constructor()
        TransferBehaviourTest(
            TransferBehaviourTest.TestCase({
                restriction: TransferRestriction.None,
                wantTransfersLocked: false,
                wantMintLocked: false,
                wantBurnLocked: false
            })
        )
    {}
}

contract OnlyMintRestrictionTest is TransferBehaviourTest {
    constructor()
        TransferBehaviourTest(
            TransferBehaviourTest.TestCase({
                restriction: TransferRestriction.OnlyMint,
                wantTransfersLocked: true,
                wantMintLocked: false,
                wantBurnLocked: true
            })
        )
    {}
}

contract OnlyBurnRestrictionTest is TransferBehaviourTest {
    constructor()
        TransferBehaviourTest(
            TransferBehaviourTest.TestCase({
                restriction: TransferRestriction.OnlyBurn,
                wantTransfersLocked: true,
                wantMintLocked: true,
                wantBurnLocked: false
            })
        )
    {}
}

contract FrozenRestrictionTest is TransferBehaviourTest {
    constructor()
        TransferBehaviourTest(
            TransferBehaviourTest.TestCase({
                restriction: TransferRestriction.Frozen,
                wantTransfersLocked: true,
                wantMintLocked: true,
                wantBurnLocked: true
            })
        )
    {}
}
