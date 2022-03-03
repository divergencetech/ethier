// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./ERC20Shield.sol";
import "./ERC721Shield.sol";
import "./ERC1155Shield.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
@notice A bare-minimum escrow contract for holding assets that can only be
reclaimed by their rightful owner. This acts as a forced vault wallet that
programatically requires the token owner to adhere to best practice of storing
assets in an address that does not perform transactions unrelated to the assets
themselves. The contract is also unaware of external protocols, e.g. Wyvern, so
shields from external vulnerabilities.
@dev Assets MUST be transferred into escrow with the safeTransferFrom() methods
as this allows for automated accounting of ownership. ERC721 and ERC1155 MUST
use their standard safeTransferFrom() methods on the respective contracts, and
ERC20 MUST approve this contract for spending before using safeTransferERC20().
NOTE THAT ALL ASSETS TRANFERRED BY A DIFFERENT MEANS WILL BE PERMANENTLY LOCKED.
@dev Under this model it is safe for multiple true beneficial owners to share
the same Shield. However this may be difficult to communicate to end users who
are accustomed to the "not your keys, not your tokens" mantra. To overcome this,
the contract can be cheaply cloned with a minimal proxy contract.
 */
contract Shield is Initializable, ERC20Shield, ERC721Shield, ERC1155Shield {
    /**
    @notice A purely "cosmetic" record of the address responsible for deploying
    this contract, which MAY be used by off-chain platforms that display
    galleries/balances (e.g. OpenSea) and wish to give a single address control
    of their profile.
    @dev NOTE that this address has absolutely no administrative control nor
    fallback ownership of assets.
     */
    address public creator;

    /**
    @dev Equivalent to a constructor, ideally called in the same transaction as
    deployment (see ShieldFactory).
     */
    function initialize(address _creator) external initializer {
        creator = _creator;
    }
}
