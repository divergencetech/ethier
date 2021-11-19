//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/**
@title SignatureChecker
@notice Modifiers to require that a function has received a valid ECDSA
signature, signed by one of a set of addresses.
 */
contract SignatureChecker {
    using EnumerableSet for EnumerableSet.AddressSet;

    /// @notice Addresses from which a signature is considered valid.
    EnumerableSet.AddressSet internal _signers;

    /**
    @param signers Initial set of signers from which signatures are considered
    valid.
     */
    constructor(address[] memory signers) {
        for (uint256 i = 0; i < signers.length; i++) {
            _signers.add(signers[i]);
        }
    }

    /**
    @notice Set of signature nonces already seen and no longer considered valid.
    */
    mapping(bytes32 => bool) internal _usedNonces;

    /**
    @notice Requires that the nonce has not been used previously and that the
    recovered signer is contained in the _signers AddressSet.
    @param signature ECDSA signature of keccak256(abi.encodePacked(data,nonce)).
     */
    modifier validSignature(
        bytes memory data,
        bytes32 nonce,
        bytes calldata signature
    ) {
        require(!_usedNonces[nonce], "SignatureChecker: Nonce already used");
        _usedNonces[nonce] = true;
        _validate(keccak256(abi.encodePacked(data, nonce)), signature);

        _;
    }

    /**
    @notice Requires that the recovered signer is contained in the _signers
    AddressSet.
    */
    modifier validReusableSignature(bytes32 hash, bytes calldata signature) {
        _validate(hash, signature);
        _;
    }

    /**
    @notice Hashes addr and requires that the recovered signer is contained in
    the _signers AddressSet.
    @dev Equivalent to validReusableSignature(sha3(addr), signature);
     */
    modifier allowedAddress(address addr, bytes calldata signature) {
        _validate(keccak256(abi.encodePacked(addr)), signature);
        _;
    }

    /**
    @notice Common validator logic, requiring that the recovered signer is
    contained in the _signers AddressSet.
     */
    function _validate(bytes32 hash, bytes calldata signature) private view {
        require(
            _signers.contains(ECDSA.recover(hash, signature)),
            "SignatureChecker: Invalid signature"
        );
    }
}
