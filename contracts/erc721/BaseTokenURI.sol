// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import {AccessControlEnumerable} from "../utils/AccessControlEnumerable.sol";
import {ERC721A, ERC721ACommon} from "./ERC721ACommon.sol";

/**
 * @notice ERC721 extension that implements a commonly used _baseURI() function
 * to return an URL prefix that can be set by the contract steerer.
 */
contract BaseTokenURI is AccessControlEnumerable {
    /**
     * @notice Base token URI used as a prefix by tokenURI().
     */
    string private _baseTokenURI;

    constructor(string memory baseTokenURI_) {
        _setBaseTokenURI(baseTokenURI_);
    }

    /**
     * @notice Sets the base token URI prefix.
     * @dev Only callable by the contract steerer.
     */
    function setBaseTokenURI(string memory baseTokenURI_)
        public
        onlyRole(DEFAULT_STEERING_ROLE)
    {
        _setBaseTokenURI(baseTokenURI_);
    }

    /**
     * @notice Sets the base token URI prefix.
     */
    function _setBaseTokenURI(string memory baseTokenURI_) internal virtual {
        _baseTokenURI = baseTokenURI_;
    }

    /**
     * @notice Returns the `baseTokenURI`.
     */
    function baseTokenURI() public view virtual returns (string memory) {
        return _baseTokenURI;
    }

    /**
     * @notice Returns the base token URI * without any additional characters (e.g. a slash).
     */
    function _baseURI() internal view virtual returns (string memory) {
        return _baseTokenURI;
    }
}

/**
 * @notice ERC721ACommon extension that adds BaseTokenURI.
 */
abstract contract ERC721ACommonBaseTokenURI is ERC721ACommon, BaseTokenURI {
    /**
     * @notice Overrides supportsInterface as required by inheritance.
     */
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

    /**
     * @dev Inheritance resolution.
     */
    function _baseURI()
        internal
        view
        virtual
        override(ERC721A, BaseTokenURI)
        returns (string memory)
    {
        return BaseTokenURI._baseURI();
    }
}
