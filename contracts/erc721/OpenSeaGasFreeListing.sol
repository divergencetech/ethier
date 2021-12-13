// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

// Inspired by BaseOpenSea by Simon Fremaux (@dievardump) but without the need
// to pass specific addresses depending on deployment network.
// https://gist.github.com/dievardump/483eb43bc6ed30b14f01e01842e3339b/

/// @notice Library to achieve gas-free listings on OpenSea.
library OpenSeaGasFreeListing {

    /**
    @notice Convinience function to get the right chainId
     */
    function getChainId() internal view returns (uint256) {
        uint256 id;
        assembly {
            id := chainid()
        }
        return id;
    }

    /**
    @notice Returns whether the operator is an OpenSea proxy for the owner, thus
    allowing it to list without the token owner paying gas.
    @dev ERC{721,1155}.isApprovedForAll should be overriden to also check if
    this function returns true.
     */
    function isApprovedForAll(address owner, address operator)
        internal
        view
        returns (bool)
    {
        uint256 chainId = getChainId();
        if (chainId == 1 || chainId == 4) {
            return isApprovedForAllWyvern(owner, operator, chainId);
        } else {
            return isApprovedForAllZeroEx(operator, chainId);
        }
    }

    /**
    @notice Returns whether the operator is an OpenSea proxy for the owner on a Wyvern exchange, thus
    allowing it to list without the token owner paying gas in Ethereum mainnet and rinkeby. 
    @dev Assumes that the passed in chainId is 1 or 4.
     */
    function isApprovedForAllWyvern(address owner, address operator, uint256 chainId)
        internal
        view
        returns (bool) 
    {
        ProxyRegistry registry; 
        assembly {
            switch chainId
            case 1 {
                // mainnet
                registry := 0xa5409ec958c83c3f309868babaca7c86dcb077c1
            }
            case 4 {
                // rinkeby
                registry := 0xf57b2c51ded3a29e6891aba85459d600256cf317
            }
        }

        return
            address(registry) != address(0) &&
            address(registry.proxies(owner)) == operator;
    }

    /**
    @notice Returns whether the operator is an OpenSea proxy for the owner on a 0x exchange, thus
    allowing it to list without the token owner paying gas in Polygon mainnet and mumbai.
    @dev Assumes that the passed in chainId is 137 or 80001.
     */
    function isApprovedForAllZeroEx(address operator, uint256 chainId)
        internal
        pure
        returns (bool) 
    {
        address registry;
        assembly {
            switch chainId
            case 137 {
                // polygon
                registry:= 0x58807baD0B376efc12F5AD86aAc70E78ed67deaE
            }
            case 80001 {
                // mumbai
                registry:= 0xff7Ca10aF37178BdD056628eF42fD7F799fAc77c
            }
        }
        return 
            address(registry) != address(0) && 
            address(registry) == operator;        
    }
}

contract OwnableDelegateProxy {}

contract ProxyRegistry {
    mapping(address => OwnableDelegateProxy) public proxies;
}
