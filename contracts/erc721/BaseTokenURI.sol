// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import {AccessControlEnumerable} from "../utils/AccessControlEnumerable.sol";

/**
@notice ERC721 extension that overrides the OpenZeppelin _baseURI() function to
return a prefix that can be set by the contract owner.
 */
contract BaseTokenURI is AccessControlEnumerable {
    /// @notice Base token URI used as a prefix by tokenURI().
    string public baseTokenURI;

    constructor(string memory _baseTokenURI) {
        _setBaseTokenURI(_baseTokenURI);
    }

    /// @notice Sets the base token URI prefix.
    /// @dev Only callable by the contract steerer.
    function setBaseTokenURI(string memory _baseTokenURI)
        public
        onlyRole(DEFAULT_STEERING_ROLE)
    {
        _setBaseTokenURI(_baseTokenURI);
    }

    /// @notice Sets the base token URI prefix.
    function _setBaseTokenURI(string memory _baseTokenURI) internal virtual {
        baseTokenURI = _baseTokenURI;
    }

    /**
    @notice Concatenates and returns the base token URI and the token ID without
    any additional characters (e.g. a slash).
    @dev This requires that an inheriting contract that also inherits from OZ's
    ERC721 will have to override both contracts; although we could simply
    require that users implement their own _baseURI() as here, this can easily
    be forgotten and the current approach guides them with compiler errors. This
    favours the latter half of "APIs should be easy to use and hard to misuse"
    from https://www.infoq.com/articles/API-Design-Joshua-Bloch/.
     */
    function _baseURI() internal view virtual returns (string memory) {
        return baseTokenURI;
    }
}
