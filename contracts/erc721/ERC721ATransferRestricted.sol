// SPDX-License-Identifier: MIT
// Copyright (c) 2023 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.16 <0.9.0;

import {ERC721ATransferRestrictedBase, TransferRestriction} from "./ERC721ATransferRestrictedBase.sol";

/**
 * @notice Extension of ERC721 transfer restrictions with manual restriction
 * setter.
 */
abstract contract ERC721ATransferRestricted is ERC721ATransferRestrictedBase {
    // =========================================================================
    //                           Error
    // =========================================================================

    error TransferRestrictionLocked();
    error TransferRestrictionCheckFailed(TransferRestriction want);

    // =========================================================================
    //                           Storage
    // =========================================================================

    /**
     * @notice The current restrictions.
     */
    TransferRestriction private _transferRestriction;

    /**
     * @notice Flag to lock in the current transfer restriction.
     */
    bool private _locked;

    // =========================================================================
    //                           Steering
    // =========================================================================

    /**
     * @notice Sets the transfer restrictions.
     */
    function _setTransferRestriction(TransferRestriction restriction)
        internal
        virtual
    {
        _transferRestriction = restriction;
    }

    /**
     * @notice Sets the transfer restrictions.
     * @dev Only callable by a contract steerer.
     */
    function setTransferRestriction(TransferRestriction restriction)
        external
        onlyRole(DEFAULT_STEERING_ROLE)
    {
        if (_locked) {
            revert TransferRestrictionLocked();
        }

        _setTransferRestriction(restriction);
    }

    /**
     * @notice Locks the current transfer restrictions.
     * @dev Only callable by a contract steerer.
     * @param restriction must match the current transfer restriction as
     * additional security measure.
     */
    function lockTransferRestriction(TransferRestriction restriction)
        external
        onlyRole(DEFAULT_STEERING_ROLE)
    {
        if (restriction != _transferRestriction) {
            revert TransferRestrictionCheckFailed(_transferRestriction);
        }

        _locked = true;
    }

    /**
     * @notice Returns the stored transfer restrictions.
     */
    function transferRestriction()
        public
        view
        virtual
        override
        returns (TransferRestriction)
    {
        return _transferRestriction;
    }
}
