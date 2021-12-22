// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/erc721/ERC721CommonEnumerable.sol";
import "../../contracts/erc721/BaseTokenURI.sol";

/// @notice Exposes a functions modified with the modifiers under test.
contract TestableERC721CommonEnumerable is
    ERC721CommonEnumerable,
    BaseTokenURI
{
    constructor() ERC721CommonEnumerable("Token", "JRR") BaseTokenURI("") {}

    function mint(uint256 tokenId) public {
        ERC721._safeMint(msg.sender, tokenId);
    }

    /// @dev For testing the tokenExists() modifier.
    function mustExist(uint256 tokenId) public view tokenExists(tokenId) {}

    /// @dev For testing the onlyApprovedOrOwner() modifier.
    function mustBeApprovedOrOwner(uint256 tokenId)
        public
        onlyApprovedOrOwner(tokenId)
    {}

    function _baseURI()
        internal
        view
        override(ERC721, BaseTokenURI)
        returns (string memory)
    {
        return BaseTokenURI._baseURI();
    }
}
