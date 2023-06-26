// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import {Pausable} from "@openzeppelin/contracts/security/Pausable.sol";
import {AccessControlEnumerable} from "./AccessControlEnumerable.sol";

/// @notice A Pausable contract that can only be toggled by a member of the
/// STEERING role.
contract AccessControlPausable is AccessControlEnumerable, Pausable {
    /// @notice Pauses the contract.
    function pause() public virtual onlyRole(DEFAULT_STEERING_ROLE) {
        Pausable._pause();
    }

    /// @notice Unpauses the contract.
    function unpause() public virtual onlyRole(DEFAULT_STEERING_ROLE) {
        Pausable._unpause();
    }
}
