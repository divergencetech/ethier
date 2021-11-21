// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "./BaseOpenSea.sol";
import "../utils/OwnerPausable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Pausable.sol";

/**
@notice An ERC721 contract with common functionality:
 - OpenSea gas-free listings
 - OpenZeppelin Enumerable and Pausable
 - OpenZeppelin Pausable with functions exposed to Owner only
 */
contract ERC721Common is
    BaseOpenSea,
    ERC721Enumerable,
    ERC721Pausable,
    OwnerPausable
{
    /**
    @param openSeaProxyRegistry Set to
    0xa5409ec958c83c3f309868babaca7c86dcb077c1 for Mainnet and
    0xf57b2c51ded3a29e6891aba85459d600256cf317 for Rinkeby.
     */
    constructor(
        string memory name,
        string memory symbol,
        address openSeaProxyRegistry
    ) ERC721(name, symbol) {
        if (openSeaProxyRegistry != address(0)) {
            BaseOpenSea._setOpenSeaRegistry(openSeaProxyRegistry);
        }
    }

    /// @notice Requires that the token exists.
    modifier tokenExists(uint256 tokenId) {
        require(ERC721._exists(tokenId), "ERC721Common: Token doesn't exist");
        _;
    }

    /// @notice Requires that msg.sender owns or is approved for the token.
    modifier onlyApprovedOrOwner(uint256 tokenId) {
        require(
            _isApprovedOrOwner(msg.sender, tokenId),
            "ERC721Common: Not approved nor owner"
        );
        _;
    }

    /// @notice Overrides _beforeTokenTransfer as required by inheritance.
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 tokenId
    ) internal virtual override(ERC721Enumerable, ERC721Pausable) {
        super._beforeTokenTransfer(from, to, tokenId);
    }

    /// @notice Overrides supportsInterface as required by inheritance.
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721Enumerable)
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
        override
        returns (bool)
    {
        return
            super.isApprovedForAll(owner, operator) ||
            BaseOpenSea.isOwnersOpenSeaProxy(owner, operator);
    }
}
