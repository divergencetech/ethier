// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/**
@title SignatureChecker
@notice Additional functions for EnumerableSet.Addresset that require a valid
ECDSA signature, signed by any member of the set.
 */
library SignatureChecker {
    using EnumerableSet for EnumerableSet.AddressSet;

    /**
    @notice Requires that the message has not been used previously and that the
    recovered signer is contained in the signers AddressSet.
    @param signers Set of addresses from which signatures are accepted.
    @param usedMessages Set of already-used messages.
    @param signature ECDSA signature of message.
     */
    function validateSignature(
        EnumerableSet.AddressSet storage signers,
        bytes32 message,
        bytes calldata signature,
        mapping(bytes32 => bool) storage usedMessages
    ) internal {
        require(!usedMessages[message], "SignatureChecker: Message already used");
        usedMessages[message] = true;
        _validate(signers, message, signature);
    }

    /**
    @notice Requires that the recovered signer is contained in the signers
    AddressSet.
    */
    function validateSignature(
        EnumerableSet.AddressSet storage signers,
        bytes32 message,
        bytes calldata signature
    ) internal view {
        _validate(signers, message, signature);
    }

    /**
    @notice Common validator logic, requiring that the recovered signer is
    contained in the _signers AddressSet.
     */
    function _validate(
        EnumerableSet.AddressSet storage signers,
        bytes32 hash,
        bytes calldata signature
    ) private view {
        require(
            signers.contains(ECDSA.recover(hash, signature)),
            "SignatureChecker: Invalid signature"
        );
    }
}
