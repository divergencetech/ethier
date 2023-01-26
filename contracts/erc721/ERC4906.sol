// SPDX-License-Identifier: CC0-1.0
pragma solidity ^0.8.0;

import {IERC165, ERC165} from "@openzeppelin/contracts/utils/introspection/ERC165.sol";

interface IERC4906Events {
    /// @dev This event emits when the metadata of a token is changed.
    /// So that the third-party platforms such as NFT market could
    /// timely update the images and related attributes of the NFT.
    event MetadataUpdate(uint256 _tokenId);

    /// @dev This event emits when the metadata of a range of tokens is changed.
    /// So that the third-party platforms such as NFT market could
    /// timely update the images and related attributes of the NFTs.
    event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId);
}

/// @title EIP-721 Metadata Update Extension
// solhint-disable-next-line no-empty-blocks
interface IERC4906 is IERC165, IERC4906Events {

}

contract ERC4906 is IERC4906, ERC165 {
    function _refreshMetadata(uint256 tokenId) internal {
        emit MetadataUpdate(tokenId);
    }

    function _refreshMetadata(uint256 fromTokenId, uint256 toTokenId) internal {
        emit BatchMetadataUpdate(fromTokenId, toTokenId);
    }

    /// @dev See {IERC165-supportsInterface}.
    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(ERC165, IERC165)
        returns (bool)
    {
        return
            interfaceId == bytes4(0x49064906) ||
            ERC165.supportsInterface(interfaceId);
    }
}
