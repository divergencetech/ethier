// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("", "") {}

    function mint(uint256 amount) external {
        _mint(msg.sender, amount);
    }
}

contract MockERC721 is ERC721 {
    constructor() ERC721("", "") {}

    function mint(uint256 tokenId) external {
        _safeMint(msg.sender, tokenId);
    }
}

contract MockERC1155 is ERC1155 {
    constructor() ERC1155("") {}

    function mint(uint256 id, uint256 amount) external {
        bytes memory data;
        _mint(msg.sender, id, amount, data);
    }
}
