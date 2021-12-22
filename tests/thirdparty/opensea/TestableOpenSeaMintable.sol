// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/thirdparty/opensea/OpenSeaERC721Mintable.sol";

/// @notice A testable implementation of OpenSeaERC721Mintable.
contract TestableOpenSeaMintable is OpenSeaERC721Mintable {
    constructor(uint256 _numFactoryOptions, string memory baseOptionURI)
        OpenSeaERC721Mintable("", "", _numFactoryOptions, baseOptionURI)
    {}

    /**
    @notice Required override to indicate if an option can currently be minted.
     */
    function factoryCanMint(uint256 optionId)
        public
        view
        override
        returns (bool)
    {
        return canMint[optionId];
    }

    mapping(uint256 => bool) public canMint;

    /// @notice Controls values returned by factoryCanMint().
    function setCanMint(uint256 optionId, bool can) public {
        canMint[optionId] = can;
    }

    /// @notice Records calls to _factoryMint().
    struct Mint {
        uint256 optionId;
        address to;
    }
    Mint[] public mints;

    function numMinted() external view returns (uint256) {
        return mints.length;
    }

    /// @notice Required override to perform actual minting.
    function _factoryMint(uint256 optionId, address to) internal override {
        mints.push(Mint({optionId: optionId, to: to}));
    }

    /**
    @dev Workaround for a bug in geth's abigen / bind package that doesn't
    create types unless they're used in function signatures.
     */
    function abigenBugHack(Mint memory) external pure {}
}
