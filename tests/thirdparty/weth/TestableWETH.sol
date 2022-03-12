// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/thirdparty/weth/WETH.sol";

contract TestableWETH {
    /**
    @notice Calls transferFrom(src, dst, wad) on the test double of the wETH
    contract, as deployed by ethier's wethtest package.
    @dev This has the effect of testing both that that package properly deploys
    the contract and that the WETH library for general use returns the IwETH
    at the correct address.
     */
    function transferFrom(
        address src,
        address dst,
        uint256 wad
    ) external {
        WETH.instance().transferFrom(src, dst, wad);
    }
}
