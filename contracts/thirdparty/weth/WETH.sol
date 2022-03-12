// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

/**
@notice Interface of wETH as deployed on Ethereum mainnet.
 */
interface IwETH {
    function balanceOf(address) external view returns (uint256);

    function allowance(address, address) external view returns (uint256);

    receive() external payable;

    function deposit() external payable;

    function withdraw(uint256) external;

    function totalSupply() external view returns (uint256);

    function approve(address, uint256) external returns (bool);

    function transfer(address, uint256) external returns (bool);

    function transferFrom(
        address,
        address,
        uint256
    ) external returns (bool);
}

/**
@notice Convenience library for wETH addresses without pre-deployment
knowlege of the chain.
@dev Chain IDs:
 - Ethereum Mainnet 1
 - Rinkeby 4
 - geth's SimulatedBackend 1337 but only compatible if using ethier'
   wethtest package
 */
library WETH {
    /// @notice Returns the wETH token address for the current chain.
    function wethAddress() internal view returns (address payable addr) {
        assembly {
            switch chainid()
            case 1 {
                addr := 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2
            }
            case 4 {
                addr := 0xc778417E063141139Fce010982780140Aa0cD5Ab
            }
            // TODO: confirm that WMATIC functions similarly before exposing
            // addresses here.
            case 1337 {
                // The geth SimulatedBackend iff used with the ethier
                // wethtest package.
                addr := 0x2336a902f2727C77867A5905dE392fEd3Ff3604b
            }
        }
    }

    /// @notice Returns a wETH token interface for the current chain.
    function instance() internal view returns (IwETH) {
        return IwETH(wethAddress());
    }
}
