// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./ERC721Common.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";

/**
@notice Extends ERC721Common functionality with ERC721Enumerable.
@dev This adds a significant gas cost to minting and transfers so only use if
absolutely necessary. If only totalSupply() is needed and the contract is also
an ethier Seller then use totalSold() as an alias.

See: https://shiny.mirror.xyz/OUampBbIz9ebEicfGnQf5At_ReMHlZy0tB4glb9xQ0E
*/
contract ERC721CommonEnumerable is ERC721Common, ERC721Enumerable {
    constructor(string memory name, string memory symbol)
        ERC721Common(name, symbol)
    {}

    /**
    @notice Returns ERC721Common.isApprovedForAll() to guarantee use of OpenSea
    gas-free listing functionality.
    */
    function isApprovedForAll(address owner, address operator)
        public
        view
        virtual
        override(ERC721Common, ERC721, IERC721)
        returns (bool)
    {
        return ERC721Common.isApprovedForAll(owner, operator);
    }

    /// @dev Calls ERC721Common.setApprovalForAll to manage pre-approvals
    /// related to OpenSea's gas-free listing functionality.
    function setApprovalForAll(address operator, bool approved)
        public
        virtual
        override(ERC721Common, ERC721, IERC721)
    {
        ERC721Common.setApprovalForAll(operator, approved);
    }

    /// @notice Overrides _beforeTokenTransfer as required by inheritance.
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 tokenId
    ) internal virtual override(ERC721Common, ERC721Enumerable) {
        super._beforeTokenTransfer(from, to, tokenId);
    }

    /// @notice Overrides supportsInterface as required by inheritance.
    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(ERC721Common, ERC721Enumerable)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
