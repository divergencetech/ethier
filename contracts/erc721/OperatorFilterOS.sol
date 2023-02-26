// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import {Address} from "@openzeppelin/contracts/utils/Address.sol";
import {DefaultOperatorFilterer} from "operator-filter-registry/src/DefaultOperatorFilterer.sol";
import {ERC721A, ERC721ACommon} from "./ERC721ACommon.sol";

/**
 * @notice ERC721ACommon extension that adds Opensea's operator filtering.
 */
abstract contract OperatorFilterOS is ERC721ACommon, DefaultOperatorFilterer {
    using Address for address;

    /**
     * @notice Calling the operator filter registry with given calldata.
     * @dev The registry contract did not foresee role-based contract access
     * control -- only the contract itself, or its (EIP-173) owner is allowed to
     * change subscription settings. To work around this, we enforce
     * authorisation here and forward arbitrary calldata to the registry.
     * Use with care!
     */
    function callOperatorFilterRegistry(bytes calldata cdata)
        external
        onlyRole(DEFAULT_STEERING_ROLE)
        returns (bytes memory)
    {
        return address(OPERATOR_FILTER_REGISTRY).functionCall(cdata);
    }

    // =========================================================================
    //                           Operator filtering
    // =========================================================================

    function setApprovalForAll(address operator, bool approved)
        public
        virtual
        override
        onlyAllowedOperatorApproval(operator)
    {
        super.setApprovalForAll(operator, approved);
    }

    function approve(address operator, uint256 tokenId)
        public
        payable
        virtual
        override
        onlyAllowedOperatorApproval(operator)
    {
        super.approve(operator, tokenId);
    }

    function transferFrom(
        address from,
        address to,
        uint256 tokenId
    ) public payable virtual override onlyAllowedOperator(from) {
        super.transferFrom(from, to, tokenId);
    }

    function safeTransferFrom(
        address from,
        address to,
        uint256 tokenId
    ) public payable virtual override onlyAllowedOperator(from) {
        super.safeTransferFrom(from, to, tokenId);
    }

    function safeTransferFrom(
        address from,
        address to,
        uint256 tokenId,
        bytes memory data
    ) public payable virtual override onlyAllowedOperator(from) {
        super.safeTransferFrom(from, to, tokenId, data);
    }
}
