// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

/**
@title MessageGenerator
@notice Generates messages that can be signed off-chain using ECDSA. 
Corresponding can later be verified on-chain using {SignatureChecker}.
 */
library MessageGenerator {
    /**
    @notice Generates a message for a given data input.
    @dev For multiple data fields, a standard concatenation using 
    `abi.encodePacked` is commonly used to build data.
     */
    function generateMessage(bytes memory data)
        internal
        pure
        returns (bytes32)
    {
        return keccak256(data);
    }

    /**
    @notice Generates a message for a given address.
     */
    function generateMessage(address address_) internal pure returns (bytes32) {
        return generateMessage(abi.encodePacked(address_));
    }
}
