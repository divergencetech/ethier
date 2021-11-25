// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/crypto/SignatureChecker.sol";
import "../../contracts/crypto/MessageGenerator.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/**
@notice Exposes functions allowing testing of SignatureChecker.
 */
contract TestableSignatureChecker {
    using EnumerableSet for EnumerableSet.AddressSet;
    // SignatureChecker adds additional functionality to an AddressSet, allowing
    // for a signature from any set member.
    using SignatureChecker for EnumerableSet.AddressSet;

    EnumerableSet.AddressSet private signers;
    mapping(bytes32 => bool) private usedMessages;

    constructor(address[] memory _signers) {
        for (uint256 i = 0; i < _signers.length; i++) {
            signers.add(_signers[i]);
        }
    }

    /// @dev Reverts if the signature is invalid or the nonce is already used.
    function needsSignature(
        bytes memory data,
        bytes32 nonce,
        bytes calldata signature
    ) external {
        bytes32 message = MessageGenerator.generateMessage(
            abi.encodePacked(data, nonce)
        );
        signers.validateSignature(message, signature, usedMessages);
    }

    /// @dev Reverts if the signature is invalid.
    function needsReusableSignature(bytes memory data, bytes calldata signature)
        external
        view
        returns (bool)
    {
        bytes32 message = MessageGenerator.generateMessage(data);
        signers.validateSignature(message, signature);
        return true;
    }

    /// @dev Reverts if the signature is not valid for keccak256(msg.sender).
    function needsSenderSignature(bytes calldata signature)
        external
        view
        returns (bool)
    {
        bytes32 message = MessageGenerator.generateMessage(msg.sender);
        signers.validateSignature(message, signature);
        return true;
    }
}
