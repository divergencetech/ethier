// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "./OpenSeaGasFreeListing.sol";
import "../utils/OwnerPausable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Pausable.sol";
import "@openzeppelin/contracts/utils/Context.sol";

/**
@notice An ERC721 contract with common functionality:
 - OpenSea gas-free listings
 - OpenZeppelin Pausable
 - OpenZeppelin Pausable with functions exposed to Owner only
 */
contract ERC721Common is Context, ERC721Pausable, OwnerPausable {
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

    /// @notice Overrides _beforeTokenTransfer as required by inheritance.
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 tokenId
    ) internal virtual override(ERC721Pausable) {
        super._beforeTokenTransfer(from, to, tokenId);
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

    /**
    @notice Returns true if either standard isApprovedForAll() returns true or
    the operator is the OpenSea proxy for the owner.
     */
    function isApprovedForAll(address owner, address operator)
        public
        view
        virtual
        override
        returns (bool)
    {
        return
            super.isApprovedForAll(owner, operator) ||
            OpenSeaGasFreeListing.isApprovedForAll(owner, operator);
    }
}
