// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./Chainlink.sol";
import "@chainlink/contracts/src/v0.8/VRFConsumerBase.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
@notice Convenience wrapper around VRFConsumerBase that abstracts the need to
know constants such as contract addresses.
@dev This contract should be used in an identical manner to Chainlink's standard
VRFConsumerBase, along with all the same best practices.
 */
abstract contract VRFConsumerHelper is VRFConsumerBase {
    constructor()
        VRFConsumerBase(Chainlink.vrfCoordinator(), Chainlink.linkToken())
    {}

    /**
    @notice Calls standard VRFConsumerBase.requestRandomness() with
    chain-specific constants.
     */
    function requestRandomness() internal returns (bytes32 requestId) {
        return
            super.requestRandomness(Chainlink.vrfKeyHash(), Chainlink.vrfFee());
    }

    /// @notice Withdraws LINK tokens, sending them to the recipient.
    function _withdrawLINK(address recipient, uint256 amount) internal {
        require(
            IERC20(Chainlink.linkToken()).transfer(recipient, amount),
            "VRFConsumerHelper: withdrawal failed"
        );
    }
}
