// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity ^0.8.15;

import {Test} from "../TestLib.sol";

import {IERC721A} from "erc721a/contracts/IERC721A.sol";
import {Address} from "@openzeppelin/contracts/utils/Address.sol";

import {OperatorFilterRegistryErrorsAndEvents, OperatorFilterRegistry} from "operator-filter-registry/src/OperatorFilterRegistry.sol";
import {OperatorFilterer} from "operator-filter-registry/src/OperatorFilterer.sol";
import {CANONICAL_OPERATOR_FILTER_REGISTRY_ADDRESS, CANONICAL_CORI_SUBSCRIPTION} from "operator-filter-registry/src/lib/Constants.sol";

import {ERC721ACommon} from "../../contracts/erc721/ERC721ACommon.sol";
import {OperatorFilterOS} from "../../contracts/erc721/OperatorFilterOS.sol";

contract TestableToken is OperatorFilterOS {
    constructor(address admin, address steerer)
        ERC721ACommon(
            admin,
            steerer,
            "NAME",
            "SYM",
            payable(address(0xFEE)),
            750
        )
    {} // solhint-disable-line no-empty-blocks

    function mint(address to, uint256 num) public {
        _mint(to, num);
    }
}

contract ERC721ATransferRestrictedGeneralTest is
    Test,
    OperatorFilterRegistryErrorsAndEvents
{
    using Address for address;

    TestableToken public token;

    address public immutable admin = makeAddr("admin");
    address public immutable steerer = makeAddr("steerer");

    OperatorFilterRegistry public immutable registry =
        OperatorFilterRegistry(CANONICAL_OPERATOR_FILTER_REGISTRY_ADDRESS);

    function setUp() public virtual {
        vm.etch(
            CANONICAL_OPERATOR_FILTER_REGISTRY_ADDRESS,
            address(new OperatorFilterRegistry()).code
        );

        vm.prank(CANONICAL_CORI_SUBSCRIPTION);
        registry.register(CANONICAL_CORI_SUBSCRIPTION);

        token = new TestableToken(admin, steerer);
    }

    function testVandalCannotSteer(address vandal) public {
        vm.assume(vandal != address(this));
        vm.assume(vandal != steerer);
        vm.assume(vandal != admin);

        vm.startPrank(vandal, steerer);

        vm.expectRevert(
            missingRoleError(vandal, token.DEFAULT_STEERING_ROLE())
        );
        token.callOperatorFilterRegistry(hex"");
    }

    function testForwardUnsubscribe(bool copyExistingEntries) public {
        assertEq(
            registry.subscriptionOf(address(token)),
            CANONICAL_CORI_SUBSCRIPTION
        );

        vm.prank(steerer);
        token.callOperatorFilterRegistry(
            abi.encodeWithSelector(
                OperatorFilterRegistry.unsubscribe.selector,
                address(token),
                copyExistingEntries
            )
        );
        assertEq(registry.subscriptionOf(address(token)), address(0));
    }

    function testForwardSubscribe(address newSubscription) public {
        vm.assume(newSubscription != address(0));
        vm.assume(newSubscription != address(token));
        vm.assume(newSubscription != CANONICAL_CORI_SUBSCRIPTION);

        vm.prank(newSubscription);
        registry.register(newSubscription);

        assertEq(
            registry.subscriptionOf(address(token)),
            CANONICAL_CORI_SUBSCRIPTION
        );

        vm.prank(steerer);
        token.callOperatorFilterRegistry(
            abi.encodeWithSelector(
                OperatorFilterRegistry.subscribe.selector,
                address(token),
                newSubscription
            )
        );
        assertEq(registry.subscriptionOf(address(token)), newSubscription);
    }

    function testForwardSubscribeWithError() public {
        vm.expectRevert(CannotSubscribeToZeroAddress.selector);
        vm.prank(steerer);
        token.callOperatorFilterRegistry(
            abi.encodeWithSelector(
                OperatorFilterRegistry.subscribe.selector,
                address(token),
                address(0)
            )
        );
    }

    function _updateOperator(address operator, bool filtered) internal {
        address[] memory list = new address[](1);
        list[0] = operator;

        vm.prank(CANONICAL_CORI_SUBSCRIPTION);
        registry.updateOperators(CANONICAL_CORI_SUBSCRIPTION, list, filtered);
    }

    enum TransferType {
        TransferFrom,
        SafeTransferFrom
    }

    enum ApprovalType {
        Approve,
        ApprovalForAll
    }

    struct TransferTest {
        address from;
        address to;
        address blockedOperator;
        address operator;
        TransferType transferType;
        ApprovalType approvalType;
    }

    function _testTransfer(TransferTest memory tt) internal {
        vm.assume(tt.from != address(0));
        vm.assume(tt.to != address(0));
        vm.assume(!tt.to.isContract());

        token.mint(tt.from, 1);
        _updateOperator(tt.blockedOperator, true);

        bytes memory err = abi.encodeWithSelector(
            AddressFiltered.selector,
            tt.blockedOperator
        );
        bool blocked = (tt.operator == tt.blockedOperator &&
            tt.operator != tt.from);

        if (tt.operator != tt.from) {
            // can't approve self

            if (blocked) {
                vm.expectRevert(err);
            }
            vm.prank(tt.from);
            if (tt.approvalType == ApprovalType.Approve) {
                token.approve(tt.operator, 0);
            } else {
                token.setApprovalForAll(tt.operator, true);
            }
        }

        if (blocked) {
            vm.expectRevert(err);
        }
        vm.prank(tt.operator);
        if (tt.transferType == TransferType.TransferFrom) {
            token.transferFrom(tt.from, tt.to, 0);
        } else {
            token.safeTransferFrom(tt.from, tt.to, 0);
        }

        assertEq(token.ownerOf(0), blocked ? tt.from : tt.to);
    }

    struct FuzzParams {
        address from;
        address to;
        address blockedOperator;
        address operator;
        uint8 transferType;
        uint8 approvalType;
    }

    function _toTestCase(FuzzParams memory fuzz)
        internal
        pure
        returns (TransferTest memory)
    {
        return
            TransferTest({
                from: fuzz.from,
                to: fuzz.to,
                blockedOperator: fuzz.blockedOperator,
                operator: fuzz.operator,
                transferType: TransferType(
                    uint8(
                        _bound(
                            fuzz.transferType,
                            0,
                            uint256(type(TransferType).max)
                        )
                    )
                ),
                approvalType: ApprovalType(
                    uint8(
                        _bound(
                            fuzz.approvalType,
                            0,
                            uint256(type(ApprovalType).max)
                        )
                    )
                )
            });
    }

    function testTransfer(FuzzParams memory fuzz) public {
        _testTransfer(_toTestCase(fuzz));
    }

    function testBlockedTransfer(FuzzParams memory fuzz) public {
        TransferTest memory tt = _toTestCase(fuzz);
        tt.operator = tt.blockedOperator;
        _testTransfer(tt);
    }

    function testEOATransfer(FuzzParams memory fuzz) public {
        TransferTest memory tt = _toTestCase(fuzz);
        tt.operator = tt.from;
        tt.blockedOperator = tt.from;
        _testTransfer(tt);
    }
}
