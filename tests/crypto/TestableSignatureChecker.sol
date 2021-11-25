// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/crypto/SignedAllowance.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/**
@notice Exposes functions allowing testing of SignatureChecker.
 */
contract TestableSignatureChecker {
    using EnumerableSet for EnumerableSet.AddressSet;

    // SignedAllowance adds additional functionality to an AddressSet, allowing
    // for a signature from any set member.
    using SignedAllowance for EnumerableSet.AddressSet;

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
        signers.validateSignature(
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
        signers.validateSignature(data, signature);
        return true;
    }

    /// @dev Reverts if the signature is not valid for keccak256(msg.sender).
    function needsSenderSignature(bytes calldata signature)
        external
        view
        returns (bool)
    {
        signers.validateSignature(msg.sender, signature);
        return true;
    }
}
