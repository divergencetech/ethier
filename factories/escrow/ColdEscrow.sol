// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/interfaces/IERC721.sol";
import "@openzeppelin/contracts/interfaces/IERC721Receiver.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
@notice A bare-minimum escrow contract for holding assets that can only be
reclaimed by their rightful owner. This acts as a forced vault wallet that
programatically requires the token owner to adhere to best practice of storing
assets in an address that does not perform transactions unrelated to the assets
themselves. The contract is also unaware of external protocols, e.g. Wyvern, so
shields from external vulnerabilities.
@dev Assets SHOULD be transferred into escrow with the safeTransferFrom() method
as this allows for automated accounting of ownership. Under this model it is
safe for multiple owners to share the same escrow. In the event that an asset is
transferred using the non-safe method, a fallback owner in the form of a
`controller` address, set immutably at deployment, takes ownership. The
controller exists also to signify the administrative EOA should platforms wish
to grant "profile" permissions for displaying galleries.
@dev Although safe for use as a multi-homed escrow for different owners, this
may be difficult to communicate to end users who are accustomed to the "not your
keys, not your tokens" mantra. To overcome this, the contract can be cheaply
cloned with a minimal proxy contract.
 */
contract ColdEscrow is IERC721Receiver, Initializable {
    /**
    @notice The address to which ownership of assets defaults in the event that
    they are received via non-safe transfer() methods. This address SHOULD also
    be given administrative rights over off-chain systems associated with
    ownership of a token (e.g. OpenSea profile, allow-list grants, membership
    gating).
     */
    address public controller;

    /**
    @dev Equivalent to a constructor. This MUST be called before receiving
    assets, ideally in the same transaction as deployment (see
    ColdEscrowFactory).
     */
    function initialize(address _controller) external initializer {
        controller = _controller;
    }

    /**
    @notice Rightful owners of ERC721 tokens transferred to this contract with
    the safeTransferFrom() method.
     */
    mapping(IERC721 => mapping(uint256 => address)) public erc721Owners;

    /**
    @dev Records the rightful owner of the token and returns the
    onERC721Received selector, in compliance with safeTransferFrom() standards.
    @param from The address recorded as the token owner.
     */
    function onERC721Received(
        address,
        address from,
        uint256 tokenId,
        bytes memory
    ) public virtual override returns (bytes4) {
        erc721Owners[IERC721(msg.sender)][tokenId] = from;
        return this.onERC721Received.selector;
    }

    /**
    @notice Reclaim an ERC721 token as its rightful owner. In the case that the
    token was received by a non-safe method (i.e. the owner is not known), the
    contract `controller` assumes ownership.
     */
    function reclaim(IERC721 token, uint256 tokenId) external {
        address owner = erc721Owners[token][tokenId];
        if (owner == address(0)) {
            owner = controller;
        }
        require(msg.sender == owner, "ColdEscrow: not owner");

        erc721Owners[token][tokenId] = address(0);
        token.safeTransferFrom(address(this), msg.sender, tokenId);
    }
}
