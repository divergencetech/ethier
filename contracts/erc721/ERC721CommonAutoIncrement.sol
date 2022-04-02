// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./ERC721Common.sol";
import "../utils/Monotonic.sol";

/**
@notice An extension of the ethier ERC721Common contract that auto-increments
the next tokenId and exposes a totalSupply.
 */
contract ERC721CommonAutoIncrement is ERC721Common {
    using Monotonic for Monotonic.Increaser;

    constructor(string memory name, string memory symbol)
        ERC721Common(name, symbol)
    {} // solhint-disable-line no-empty-blocks

    /**
    @notice Total number of tokens minted.
     */
    Monotonic.Increaser public totalSupply;

    /**
    @dev _safeMint() n tokens to the specified address, only incrementing
    totalSupply once.
     */
    function _safeMintN(
        address to,
        uint256 n,
        bytes memory data
    ) internal {
        uint256 tokenId = totalSupply.current();
        uint256 end = tokenId + n;
        for (; tokenId < end; ++tokenId) {
            _safeMint(to, tokenId, data);
        }
        totalSupply.add(n);
    }

    /**
    @dev Alias for _safeMintN(address,uint,bytes) assuming an empty byte buffer.
     */
    function _safeMintN(address to, uint256 n) internal {
        _safeMintN(to, n, "");
    }
}
