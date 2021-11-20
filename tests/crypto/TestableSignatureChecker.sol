// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/crypto/SignatureChecker.sol";

/**
@notice Exposes functions allowing testing of SignatureChecker.
 */
contract TestableSignatureChecker is SignatureChecker {
    constructor(address[] memory signers) SignatureChecker(signers) {}

    /// @dev Reverts if the signature is invalid or the nonce is already used.
    function needsSignature(
        bytes memory data,
        bytes32 nonce,
        bytes calldata signature
    ) external validSignature(data, nonce, signature) {}

    /// @dev Reverts if the signature is invalid.
    function needsReusableSignature(bytes memory data, bytes calldata signature)
        external
        view
        validReusableSignature(keccak256(data), signature)
        returns (bool)
    {
        return true;
    }

    /// @dev Reverts if the signature is not valid for keccak256(msg.sender).
    function needsSenderSignature(bytes calldata signature)
        external
        view
        allowedAddress(msg.sender, signature)
        returns (bool)
    {
        return true;
    }
}
