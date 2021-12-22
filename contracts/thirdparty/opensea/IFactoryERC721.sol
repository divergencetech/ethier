// SPDX-License-Identifier: MIT
// Source: https://github.com/ProjectOpenSea/opensea-creatures/blob/d0999de0a735b3b5c41feb7663ab88dc6407158a/contracts/IFactoryERC721.sol
// Copyright ProjectOpenSea

/**
This is copied verbatim from the OpenSea source (under MIT license) except for
the change of the interface name to include the I prefix. The source repository
is a tutorial for creating OpenSea-compatible factory contracts, so we DO NOT
want to add the entire thing as an ethier dependency.
 */

pragma solidity ^0.8.0;

/**
 * This is a generic factory contract that can be used to mint tokens. The configuration
 * for minting is specified by an _optionId, which can be used to delineate various
 * ways of minting.
 */
interface IFactoryERC721 {
    /**
     * Returns the name of this factory.
     */
    function name() external view returns (string memory);

    /**
     * Returns the symbol for this factory.
     */
    function symbol() external view returns (string memory);

    /**
     * Number of options the factory supports.
     */
    function numOptions() external view returns (uint256);

    /**
     * @dev Returns whether the option ID can be minted. Can return false if the developer wishes to
     * restrict a total supply per option ID (or overall).
     */
    function canMint(uint256 _optionId) external view returns (bool);

    /**
     * @dev Returns a URL specifying some metadata about the option. This metadata can be of the
     * same structure as the ERC721 metadata.
     */
    function tokenURI(uint256 _optionId) external view returns (string memory);

    /**
     * Indicates that this is a factory contract. Ideally would use EIP 165 supportsInterface()
     */
    function supportsFactoryInterface() external view returns (bool);

    /**
     * @dev Mints asset(s) in accordance to a specific address with a particular "option". This should be
     * callable only by the contract owner or the owner's Wyvern Proxy (later universal login will solve this).
     * Options should also be delineated 0 - (numOptions() - 1) for convenient indexing.
     * @param _optionId the option id
     * @param _toAddress address of the future owner of the asset(s)
     */
    function mint(uint256 _optionId, address _toAddress) external;
}
