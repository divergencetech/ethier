// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/interfaces/IERC20.sol";

/**
@dev The ERC20 component of a Shield contract.
 */
contract ERC20Shield {
    /**
    @notice True balances of ERC20 holders, as transferred to this contract
    with the safeTransferFrom() method.
     */
    mapping(IERC20 => mapping(address => uint256)) private erc20Balances;

    /**
    @notice Equivalent of ERC20 balanceof, but namespaced by the token.
    @return The true balance of the beneficial owner of the token.
     */
    function balanceOf(IERC20 token, address account)
        external
        view
        returns (uint256)
    {
        return erc20Balances[token][account];
    }

    /**
    @notice Emitted by safeTransferERC20From, propagating the `from` address as
    `owner` in the event log.
     */
    event ERC20Shielded(
        IERC20 indexed token,
        address indexed owner,
        uint256 amount
    );

    /**
    @notice Emitted when an ERC20 balance is returned to its true owner.
     */
    event ERC20Unshielded(
        IERC20 indexed token,
        address indexed owner,
        uint256 amount
    );

    /**
    @notice "Safely" transfers a balance from the message sender to this
    contract, where "safe" implies the same as ERC721/ERC1155 (i.e. the
    recipient address is "aware" of the standard and the balance is
    recoverable).
    @dev This contract MUST already be approved to spend the message sender's
    balance of the token.
     */
    function safeTransferERC20From(IERC20 token, uint256 amount) external {
        token.transferFrom(msg.sender, address(this), amount);
        erc20Balances[token][msg.sender] += amount;
        emit ERC20Shielded(token, msg.sender, amount);
    }

    /**
    @notice Reclaim an ERC20 balance as its rightful owner.
     */
    function reclaimERC20(IERC20 token, uint256 amount) external {
        // CHECKS
        require(
            erc20Balances[token][msg.sender] >= amount,
            "ERC20Shield: insufficient balance"
        );
        // EFFECTS
        erc20Balances[token][msg.sender] -= amount;
        // INTERACTIONS
        token.transferFrom(address(this), msg.sender, amount);

        emit ERC20Unshielded(token, msg.sender, amount);
    }
}
