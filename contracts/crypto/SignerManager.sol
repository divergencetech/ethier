// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/**
@title SignerManager
@notice Manges addition and removal of a core set of addresses from which
valid ECDSA signatures can be accepted; see SignatureChecker.
 */
contract SignerManager is Ownable {
    using EnumerableSet for EnumerableSet.AddressSet;

    /**
    @dev Addresses from which signatures can be accepted.
     */
    EnumerableSet.AddressSet internal signers;

    /**
    @notice Add an address to the set of accepted signers.
     */
    function addSigner(address signer) external onlyOwner {
        signers.add(signer);
    }

    /**
    @notice Remove an address previously added with addSigner().
     */
    function removeSigner(address signer) external onlyOwner {
        signers.remove(signer);
    }
}
