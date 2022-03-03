// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/interfaces/IERC721.sol";
import "@openzeppelin/contracts/interfaces/IERC721Receiver.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";

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
    @dev The value to which the release block is set when an asset is frozen;
    i.e. (2^256 - 1).
     */
    uint256 private constant FROZEN = ~uint256(0);

    // TODO(aschlosberg): abstract ReleaseCriteria and associated functionality
    // (i.e. freeze / unfreeze / restrictions) into ShieldLib for sharing across
    // ERC20/721/1155.

    /**
    @dev Criteria dictating when an asset MAY be reclaimed. When an asset is in
    a frozen state, earliestBlock is set to max(uint256), completely locking it
    within this contract. When an asset is unfrozen it begins a thawing period
    and earliestBlock is set to (block.number + thawPeriod). Reclaiming an asset
    MUST NOT be allowed if block.number < earliestBlock.
     */
    struct ReleaseCriteria {
        uint256 earliestBlock;
        uint256 thawPeriod;
    }

    /**
    @dev Criteria for release of a specific token.
     */
    mapping(IERC721 => mapping(uint256 => ReleaseCriteria))
        public erc721Release;

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
    @dev Requires that msg.sender is the rightful owner of the specific token.
     */
    modifier onlyERC721Owner(IERC721 token, uint256 tokenId) {
        require(
            msg.sender == erc721Owners[token][tokenId],
            "ERC721Shield: not owner"
        );
        _;
    }

    // TODO(aschlosberg): events for ERC721Frozen and ERC721Unfrozen.

    /**
    @notice Sets the token's release block to effectively infinite (2^256 - 1)
    and stores the thawing period for when unfreeze() is called.
    @dev freeze() can be called at any time, even while already frozen or during
    the thaw period, but it MUST NOT reduce the shortest time to release and
    will otherwise revert.
    @param thawPeriod When unfreeze() is called, the token's release block will
    be set to (block.number + thawPeriod).
     */
    function freeze(
        IERC721 token,
        uint256 tokenId,
        uint256 thawPeriod
    ) external onlyERC721Owner(token, tokenId) {
        ReleaseCriteria memory release = erc721Release[token][tokenId];
        uint256 soonest = release.earliestBlock == FROZEN
            ? release.thawPeriod
            : release.earliestBlock - block.number;
        require(thawPeriod >= soonest, "ERC721Shield: thaw reduction");

        erc721Release[token][tokenId] = ReleaseCriteria({
            earliestBlock: FROZEN,
            thawPeriod: thawPeriod
        });
    }

    /**
    @notice Sets the token's release block to current plus `thawPeriod` as
    passed to freeze(). After that many blocks, the token can be transfered to
    its rightful owner via reclaimERC721().
     */
    function unfreeze(IERC721 token, uint256 tokenId)
        external
        onlyERC721Owner(token, tokenId)
    {
        ReleaseCriteria storage release = erc721Release[token][tokenId];
        require(release.earliestBlock == FROZEN, "ERC721Shield: not frozen");
        release.earliestBlock = block.number + release.thawPeriod;
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
    ) external onlyERC721Owner(token, tokenId) {
        // CHECKS
        require(
            erc721Release[token][tokenId].earliestBlock <= block.number,
            "ERC721Shield: frozen"
        );
        // EFFECTS
        erc721Owners[token][tokenId] = address(0);
        // INTERACTIONS
        token.safeTransferFrom(address(this), msg.sender, tokenId, data);

        emit ERC721Unshielded(token, msg.sender, tokenId);
    }
}
