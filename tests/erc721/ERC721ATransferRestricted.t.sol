// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import {Test} from "../TestLib.sol";
import {console2} from "forge-std/console2.sol";

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
    {} // solhint-disable-line no-empty-blocks

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

    enum TransferFunction {
        TransferFrom,
        SafeTransferFrom,
        BypassedTransferFrom
    }

    struct TransferTestFuzzParams {
        address from;
        address to;
        uint8 tokenId;
        uint8 transferFunction; // Unfortunately fuzzing for enums is buggy
    }

    struct TransferTestCase {
        TransferTestFuzzParams fuzz;
        address transferCaller;
        bool callerApproved;
        bool callerApprovedForAll;
    }

    function _expectRevertIfLocked(bool locked, bytes memory error) internal {
        if (locked) {
            vm.expectRevert(error);
        }
    }

    function _testTransferFrom(TransferTestCase memory ttt) internal {
        _checkAndSetup(ttt.fuzz.from, ttt.fuzz.to, ttt.fuzz.tokenId);
        TransferFunction transferFunction = TransferFunction(
            ttt.fuzz.transferFunction % uint8(type(TransferFunction).max)
        );

        bool differentCaller = ttt.transferCaller != ttt.fuzz.from;
        if (differentCaller) {
            if (ttt.callerApproved) {
                _expectRevertIfLocked(tt.wantTransfersLocked, _lockedErr);
                vm.prank(ttt.fuzz.from);
                token.approve(ttt.transferCaller, ttt.fuzz.tokenId);
            }
            if (ttt.callerApprovedForAll) {
                _expectRevertIfLocked(tt.wantTransfersLocked, _lockedErr);
                vm.prank(ttt.fuzz.from);
                token.setApprovalForAll(ttt.transferCaller, true);
            }
        }

        bool approved = token.getApproved(ttt.fuzz.tokenId) ==
            ttt.transferCaller ||
            token.isApprovedForAll(ttt.fuzz.from, ttt.transferCaller);

        bool locked = tt.wantTransfersLocked &&
            transferFunction != TransferFunction.BypassedTransferFrom;

        bytes memory err;
        if (locked) {
            err = _lockedErr;
        }
        if (differentCaller && !approved) {
            err = _notApprovedErr;
        }
        bool fails = err.length > 0;
        _expectRevertIfLocked(fails, err);

        vm.prank(ttt.transferCaller);
        if (transferFunction == TransferFunction.TransferFrom) {
            token.transferFrom(ttt.fuzz.from, ttt.fuzz.to, ttt.fuzz.tokenId);
        }
        if (transferFunction == TransferFunction.SafeTransferFrom) {
            token.safeTransferFrom(
                ttt.fuzz.from,
                ttt.fuzz.to,
                ttt.fuzz.tokenId
            );
        }
        if (transferFunction == TransferFunction.BypassedTransferFrom) {
            token.bypassedTransferFrom(
                ttt.fuzz.from,
                ttt.fuzz.to,
                ttt.fuzz.tokenId
            );
        }

        assertEq(
            token.ownerOf(ttt.fuzz.tokenId),
            fails ? ttt.fuzz.from : ttt.fuzz.to
        );
    }

    function testOwnerTransfer(TransferTestFuzzParams memory fuzz) public {
        _testTransferFrom(
            TransferTestCase({
                fuzz: fuzz,
                transferCaller: fuzz.from,
                callerApproved: false,
                callerApprovedForAll: false
            })
        );
    }

    function testVandalTransfer(
        TransferTestFuzzParams memory fuzz,
        address vandal
    ) public {
        vm.assume(vandal != fuzz.from);
        _testTransferFrom(
            TransferTestCase({
                fuzz: fuzz,
                transferCaller: vandal,
                callerApproved: false,
                callerApprovedForAll: false
            })
        );
    }

    function testApprovedTransfer(
        TransferTestFuzzParams memory fuzz,
        address caller
    ) public {
        vm.assume(caller != fuzz.from);
        _testTransferFrom(
            TransferTestCase({
                fuzz: fuzz,
                transferCaller: caller,
                callerApproved: true,
                callerApprovedForAll: false
            })
        );
    }

    function testApprovedForAllTransfer(
        TransferTestFuzzParams memory fuzz,
        address caller
    ) public {
        vm.assume(caller != fuzz.from);
        _testTransferFrom(
            TransferTestCase({
                fuzz: fuzz,
                transferCaller: caller,
                callerApproved: false,
                callerApprovedForAll: true
            })
        );
    }

    function testMint(address alice, uint8 num) public {
        vm.assume(alice != address(0));
        vm.assume(num > 0);
        token.setTransferRestriction(tt.restriction);
        uint256 startSupply = token.totalSupply();

        _expectRevertIfLocked(tt.wantMintLocked, _lockedErr);
        token.mint(alice, num);

        assertEq(
            token.balanceOf(alice),
            startSupply + (tt.wantMintLocked ? 0 : num)
        );
    }

    function testBurn(address alice, uint8 tokenId) public {
        vm.assume(alice != address(0));
        token.mint(alice, uint256(tokenId) + 1);
        token.setTransferRestriction(tt.restriction);
        uint256 startSupply = token.totalSupply();

        _expectRevertIfLocked(tt.wantBurnLocked, _lockedErr);
        token.burn(tokenId);

        assertEq(
            token.balanceOf(alice),
            startSupply - (tt.wantBurnLocked ? 0 : 1)
        );
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
    {} // solhint-disable-line no-empty-blocks
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
    {} // solhint-disable-line no-empty-blocks
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
    {} // solhint-disable-line no-empty-blocks
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
    {} // solhint-disable-line no-empty-blocks
}
