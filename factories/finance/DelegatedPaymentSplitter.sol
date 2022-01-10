// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts-upgradeable/finance/PaymentSplitterUpgradeable.sol";

/**
@notice This contract functions identically to a standard OpenZeppelin
PaymentSplitter except that it can be cheaply cloned and deployed via the
PaymentSplitterFactory. The upgradeable functionality is not used, but is
required for cloning with an EIP-1677 minimal contract proxy.
@dev Cloning only replicates the implementation logic, but not the data
associated with each clone. See EIP-1677 for details.

NOTE: there is likely no need to import this contract directly; instead see the
ethier documentation for the deployed factory addresses.
 */
contract DelegatedPaymentSplitter is PaymentSplitterUpgradeable {
    /**
    @dev Initializes the PaymentSplitter, akin to a constructor. MUST be called
    by the Factory, in the same transaction as deployment, as there are no
    protections in place.
     */
    function initialize(address[] memory payees, uint256[] memory shares)
        external
        payable
        initializer
    {
        __PaymentSplitter_init(payees, shares);
    }
}
