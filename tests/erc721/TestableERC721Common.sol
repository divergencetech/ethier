// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/erc721/ERC721Common.sol";

/// @notice Exposes a functions modified with the modifiers under test.
contract TestableERC721Common is ERC721Common {
    constructor() ERC721Common("Token", "JRR", address(0)) {}

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
}
