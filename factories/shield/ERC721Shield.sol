// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/interfaces/IERC721.sol";
import "@openzeppelin/contracts/interfaces/IERC721Receiver.sol";

/**
@dev The ERC721 component of a Shield contract.
 */
contract ERC721Shield is IERC721Receiver {
    /**
    @notice True owners of ERC721 tokens transferred to this contract with the
    safeTransferFrom() method.
     */
    mapping(IERC721 => mapping(uint256 => address)) private erc721Owners;

    /**
    @notice Equivalent of ERC721 ownerOf, but namespaced by the token.
    @dev Only valid if token.ownerOf(tokenId) returns this contract.
    @return The true benficial owner of the token.
     */
    function ownerOf(IERC721 token, uint256 tokenId)
        external
        view
        returns (address)
    {
        return erc721Owners[token][tokenId];
    }

    /**
    @notice Emitted by onERC721Received, propagating the `from` address as
    `owner` in the event log.
     */
    event ERC721Shielded(
        IERC721 indexed token,
        address indexed owner,
        uint256 tokenId
    );

    /**
    @notice Emitted when an ERC721 token is returned to its true owner.
     */
    event ERC721Unshielded(
        IERC721 indexed token,
        address indexed owner,
        uint256 tokenId
    );

    /**
    @dev Records the rightful owner of the token and returns the
    onERC721Received selector, in compliance with safeTransferFrom() standards.
    @param from The address recorded as the token owner.
     */
    function onERC721Received(
        address,
        address from,
        uint256 tokenId,
        bytes memory
    ) public override returns (bytes4) {
        erc721Owners[IERC721(msg.sender)][tokenId] = from;
        emit ERC721Shielded(IERC721(msg.sender), from, tokenId);
        return this.onERC721Received.selector;
    }

    /**
    @notice Reclaim an ERC721 token as its rightful owner.
    @dev NOTE In the case that the token was received by a non-safe method (i.e.
    the owner is not known), it is permanently locked.
    @param data Data parameter piped to token.safeTransferFrom().
     */
    function reclaimERC721(
        IERC721 token,
        uint256 tokenId,
        bytes memory data
    ) external {
        // CHECKS
        require(
            msg.sender == erc721Owners[token][tokenId],
            "ERC721Shield: not owner"
        );
        // EFFECTS
        erc721Owners[token][tokenId] = address(0);
        // INTERACTIONS
        token.safeTransferFrom(address(this), msg.sender, tokenId, data);

        emit ERC721Unshielded(token, msg.sender, tokenId);
    }
}
