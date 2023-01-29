// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import {ERC721ACommon, ERC721A} from "../../contracts/erc721/ERC721ACommon.sol";
import {BaseTokenURI} from "../../contracts/erc721/BaseTokenURI.sol";
import {AccessControlEnumerable} from "../../contracts/utils/AccessControlEnumerable.sol";

/// @notice Exposes a functions modified with the modifiers under test.
contract TestableERC721ACommon is ERC721ACommon, BaseTokenURI {
    // solhint-disable-next-line no-empty-blocks
    constructor(
        address admin,
        address steerer,
        address payable royaltyReciever,
        uint96 royaltyBasisPoints
    )
        ERC721ACommon(
            admin,
            steerer,
            "Token",
            "JRR",
            royaltyReciever,
            royaltyBasisPoints
        )
        BaseTokenURI("")
    {} // solhint-disable-line no-empty-blocks

    function mint() public {
        mintN(1);
    }

    function mintN(uint256 num) public {
        ERC721A._safeMint(msg.sender, num);
    }

    function burn(uint256 tokenId) public {
        ERC721A._burn(tokenId);
    }

    /// @dev For testing the tokenExists() modifier.
    // solhint-disable-next-line no-empty-blocks
    function mustExist(uint256 tokenId) public view tokenExists(tokenId) {}

    /// @dev For testing the onlyApprovedOrOwner() modifier.
    function mustBeApprovedOrOwner(uint256 tokenId)
        public
        onlyApprovedOrOwner(tokenId)
    {} // solhint-disable-line no-empty-blocks

    function _baseURI()
        internal
        view
        override(ERC721A, BaseTokenURI)
        returns (string memory)
    {
        return BaseTokenURI._baseURI();
    }

    /// @notice Overrides supportsInterface as required by inheritance.
    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(ERC721ACommon, AccessControlEnumerable)
        returns (bool)
    {
        return
            ERC721ACommon.supportsInterface(interfaceId) ||
            AccessControlEnumerable.supportsInterface(interfaceId);
    }
}
