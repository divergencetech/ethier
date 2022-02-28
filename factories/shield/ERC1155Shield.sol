// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/interfaces/IERC1155.sol";
import "@openzeppelin/contracts/interfaces/IERC1155Receiver.sol";

/**
@dev The ERC1155 component of a Shield contract.
 */
contract ERC1155Shield is IERC1155Receiver {
    /**
    @notice True balances of ERC1155 holders, as transferred to this contract
    with the safeTransferFrom() method.
    @dev Mapping from token contract => id => holder => balance.
     */
    mapping(IERC1155 => mapping(uint256 => mapping(address => uint256)))
        private erc1155Balances;

    /**
    @notice Equivalent of ERC1155 balanceof, but namespaced by the token.
    @return The true balance of the beneficial owner of the token.
     */
    function balanceOf(
        IERC1155 token,
        address account,
        uint256 id
    ) external view returns (uint256) {
        return erc1155Balances[token][id][account];
    }

    /**
    @notice Emitted by onERC1155Received, propagating the `from` address as
    `owner` in the event log.
     */
    event ERC1155Shielded(
        IERC1155 indexed token,
        uint256 indexed id,
        address indexed owner,
        uint256 amount
    );

    /**
    @notice Emitted when an ERC1155 token is returned to its true owner.
     */
    event ERC1155Unshielded(
        IERC1155 indexed token,
        uint256 indexed id,
        address indexed owner,
        uint256 amount
    );

    /**
    @dev Records the rightful owner of the tokens and returns the
    onERC1155Received selector, in compliance with safeTransferFrom() standards.
    @param from The address recorded as the token owner.
     */
    function onERC1155Received(
        address,
        address from,
        uint256 id,
        uint256 value,
        bytes memory
    ) public override returns (bytes4) {
        erc1155Balances[IERC1155(msg.sender)][id][from] += value;
        emit ERC1155Shielded(IERC1155(msg.sender), id, from, value);
        return this.onERC1155Received.selector;
    }

    /**
    @dev Batch equivalent of onERC1155Received().
     */
    function onERC1155BatchReceived(
        address,
        address from,
        uint256[] memory ids,
        uint256[] memory values,
        bytes memory data
    ) public override returns (bytes4) {
        require(
            ids.length == values.length,
            "ERC1155Shield: arrays lengths differ"
        );
        for (uint256 i = 0; i < ids.length; i++) {
            onERC1155Received(address(0), from, ids[i], values[i], data);
        }
        return this.onERC1155BatchReceived.selector;
    }

    /**
    @notice Reclaim an ERC1155 token as its rightful owner. In the case that the
    token was received by a non-safe method (i.e. the owner is not known), the
    contract `controller` assumes ownership.
    @param data Data parameter piped to token.safeTransferFrom().
     */
    function reclaimERC1155(
        IERC1155 token,
        uint256 id,
        uint256 amount,
        bytes memory data
    ) external {
        // CHECKS
        require(
            erc1155Balances[token][id][msg.sender] >= amount,
            "ERC1155Shield: insufficient balance"
        );
        // EFFECTS
        erc1155Balances[token][id][msg.sender] -= amount;
        // INTERACTIONS
        token.safeTransferFrom(address(this), msg.sender, id, amount, data);

        emit ERC1155Unshielded(token, id, msg.sender, amount);
    }

    /**
    @dev Batch equivalent of reclaimERC1155().
     */
    function batchReclaimERC1155(
        IERC1155 token,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) external {
        require(
            ids.length == amounts.length,
            "ERC1155Shield: arrays lengths differ"
        );

        for (uint256 i = 0; i < ids.length; i++) {
            // CHECKS
            require(
                erc1155Balances[token][ids[i]][msg.sender] >= amounts[i],
                "ERC1155Shield: insufficient balance"
            );
            // EFFECTS
            erc1155Balances[token][ids[i]][msg.sender] -= amounts[i];
        }

        // INTERACTIONS
        for (uint256 i = 0; i < ids.length; i++) {
            token.safeTransferFrom(
                address(this),
                msg.sender,
                ids[i],
                amounts[i],
                data
            );
            emit ERC1155Unshielded(token, ids[i], msg.sender, amounts[i]);
        }
    }

    /**
    @notice Retrusn true i.f.f. interfaceId is that of IERC1155Receiver or
    IERC165.
     */
    function supportsInterface(bytes4 interfaceId)
        external
        pure
        virtual
        returns (bool)
    {
        return
            interfaceId == type(IERC1155Receiver).interfaceId ||
            interfaceId == type(IERC165).interfaceId;
    }
}
