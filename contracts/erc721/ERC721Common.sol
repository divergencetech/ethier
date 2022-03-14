// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./ERC721PreApproval.sol";
import "../utils/OwnerPausable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Pausable.sol";

/**
@notice An ERC721 contract with common functionality:
 - OpenSea gas-free listings
 - OpenZeppelin Pausable
 - OpenZeppelin Pausable with functions exposed to Owner only
 */
contract ERC721Common is ERC721Pausable, ERC721PreApproval, OwnerPausable {
    constructor(string memory name, string memory symbol)
        ERC721(name, symbol)
    {}

    /// @notice Requires that the token exists.
    modifier tokenExists(uint256 tokenId) {
        require(ERC721._exists(tokenId), "ERC721Common: Token doesn't exist");
        _;
    }

    /// @notice Requires that msg.sender owns or is approved for the token.
    modifier onlyApprovedOrOwner(uint256 tokenId) {
        require(
            _isApprovedOrOwner(_msgSender(), tokenId),
            "ERC721Common: Not approved nor owner"
        );
        _;
    }

    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 tokenId
    ) internal virtual override(ERC721Pausable, ERC721PreApproval) {
        super._beforeTokenTransfer(from, to, tokenId);
    }

    function setApprovalForAll(address operator, bool approved)
        public
        virtual
        override(ERC721, ERC721PreApproval)
    {
        ERC721PreApproval.setApprovalForAll(operator, approved);
    }

    function isApprovedForAll(address owner, address operator)
        public
        view
        virtual
        override(ERC721, ERC721PreApproval)
        returns (bool)
    {
        return ERC721PreApproval.isApprovedForAll(owner, operator);
    }

    /// @notice Overrides supportsInterface as required by inheritance.
    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(ERC721)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
