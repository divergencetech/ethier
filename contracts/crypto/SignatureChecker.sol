// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
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
    @notice Requires that the nonce has not been used previously and that the
    recovered signer is contained in the signers AddressSet.
    @param signers Set of addresses from which signatures are accepted.
    @param usedNonces Set of already-used nonces.
    @param signature ECDSA signature of keccak256(abi.encodePacked(data,nonce)).
     */
    function validateSignature(
        EnumerableSet.AddressSet storage signers,
        bytes memory data,
        bytes32 nonce,
        bytes calldata signature,
        mapping(bytes32 => bool) storage usedNonces
    ) internal {
        require(!usedNonces[nonce], "SignatureChecker: Nonce already used");
        usedNonces[nonce] = true;
        _validate(signers, keccak256(abi.encodePacked(data, nonce)), signature);
    }

    /**
    @notice Requires that the recovered signer is contained in the signers
    AddressSet.
    */
    function validateSignature(
        EnumerableSet.AddressSet storage signers,
        bytes32 hash,
        bytes calldata signature
    ) internal view {
        _validate(signers, hash, signature);
    }

    /**
    @notice Hashes addr and requires that the recovered signer is contained in
    the signers AddressSet.
    @dev Equivalent to validate(sha3(addr), signature);
     */
    function validateSignature(
        EnumerableSet.AddressSet storage signers,
        address addr,
        bytes calldata signature
    ) internal view {
        _validate(signers, keccak256(abi.encodePacked(addr)), signature);
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
