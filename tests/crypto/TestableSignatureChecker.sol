// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/crypto/SignatureChecker.sol";
import "../../contracts/crypto/SignerManager.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/**
@notice Exposes functions allowing testing of SignatureChecker.
 */
contract TestableSignatureChecker is SignerManager {
    using EnumerableSet for EnumerableSet.AddressSet;
    // SignatureChecker adds additional functionality to an AddressSet, allowing
    // for a signature from any set member.
    using SignatureChecker for EnumerableSet.AddressSet;

    mapping(bytes32 => bool) private usedMessages;

    /// @dev Reverts if the signature is invalid or the nonce is already used.
    function needsSignature(
        bytes memory data,
        bytes32 nonce,
        bytes calldata signature
    ) external {
        signers.requireValidSignature(
            abi.encodePacked(data, nonce),
            signature,
            usedMessages
        );
    }

    /// @dev Reverts if the signature is invalid.
    function needsReusableSignature(bytes memory data, bytes calldata signature)
        external
        view
        returns (bool)
    {
        signers.requireValidSignature(data, signature);
        return true;
    }

    /// @dev Reverts if the signature is not valid for msg.sender.
    function needsSenderSignature(bytes calldata signature)
        external
        view
        returns (bool)
    {
        signers.requireValidSignature(msg.sender, signature);
        return true;
    }
}
