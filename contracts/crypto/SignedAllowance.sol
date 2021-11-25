// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Divergent Technologies Ltd (github.com/divergencetech)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "./SignatureChecker.sol";
import "./MessageGenerator.sol";

/**
@title SignedAllowance
@notice Additional functions for EnumerableSet.Addresset that require a valid
ECDSA signature of a standardized message, signed by any member of the set.
@dev This is a convenience library combining the functionalities of 
{SignatureChecker} and {MessageGenerator}.
 */
library SignedAllowance {
    using EnumerableSet for EnumerableSet.AddressSet;

    // SignatureChecker adds additional functionality to an AddressSet, allowing
    // for a signature from any set member.
    using SignatureChecker for EnumerableSet.AddressSet;

    /**
    @notice Requires that the message has not been used previously and that the
    recovered signer is contained in the signers AddressSet.
    @param signers Set of addresses from which signatures are accepted.
    @param usedMessages Set of already-used messages.
    @param signature ECDSA signature of message.
     */
    function validateSignature(
        EnumerableSet.AddressSet storage signers,
        bytes memory data,
        bytes calldata signature,
        mapping(bytes32 => bool) storage usedMessages
    ) internal {
        bytes32 message = MessageGenerator.generateMessage(data);
        return signers.validateSignature(message, signature, usedMessages);
    }

    /**
    @notice Requires that the message has not been used previously and that the
    recovered signer is contained in the signers AddressSet.
     */
    function validateSignature(
        EnumerableSet.AddressSet storage signers,
        bytes memory data,
        bytes calldata signature
    ) internal view {
        bytes32 message = MessageGenerator.generateMessage(data);
        return signers.validateSignature(message, signature);
    }

    /**
    @notice Requires that the message has not been used previously and that the
    recovered signer is contained in the signers AddressSet.
     */
    function validateSignature(
        EnumerableSet.AddressSet storage signers,
        address addr,
        bytes calldata signature
    ) internal view {
        bytes32 message = MessageGenerator.generateMessage(addr);
        return signers.validateSignature(message, signature);
    }
}
