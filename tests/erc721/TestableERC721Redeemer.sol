// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../contracts/erc721/ERC721Redeemer.sol";
import "@openzeppelin/contracts/token/ERC721/presets/ERC721PresetMinterPauserAutoId.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/// @notice Exposes the ERC721Redeemer library for testing.
contract TestableERC721Redeemer {
    using EnumerableSet for EnumerableSet.AddressSet;
    using ERC721Redeemer for ERC721Redeemer.Claims;
    using ERC721Redeemer for ERC721Redeemer.SingleClaims;
    using Strings for uint256;

    /**
    @notice The token against which claims are redeemed.
     */
    ERC721PresetMinterPauserAutoId public immutable token;

    constructor() {
        token = new ERC721PresetMinterPauserAutoId("", "", "");
    }

    /**
    @notice Exposes minting to the test code.
     */
    function mint(address to, uint256 n) public {
        for (uint256 i = 0; i < n; i++) {
            token.mint(to);
        }
    }

    /**
    @notice Record of already-redeemed claims, used by ERC721Redeemer.redeem().
     */
    ERC721Redeemer.Claims private claims;

    /**
    @notice Equivalent of `claims` but only a single claim can be made against
    each.
     */
    ERC721Redeemer.SingleClaims private singleClaims;

    /**
    @notice Record of per-address redemptions, agnostic to specific token IDs;
    value under test.
     */
    mapping(address => uint256) public _redeemed;

    /**
    @notice Record of all successful redeemers, allowing for enumeration of the
    `redeemed` mapping.
     */
    EnumerableSet.AddressSet private _redeemers;

    /**
    @notice Exposes the Claims.redeem() function publicly.
    @dev DO NOT use this pattern in production as it contains a vulnerability by
    allowing the user to set the max allowance, n. Real implementations MUST set
    n internally.
     */
    function redeemMaxN(uint256 n, uint256[] calldata tokenIds) public {
        _redeemed[msg.sender] += claims.redeem(n, msg.sender, token, tokenIds);
        _redeemers.add(msg.sender);
    }

    /**
    @notice Exposes the SingleClaims.redeem() function publicly.
     */
    function redeemFromSingle(uint256[] calldata tokenIds) public {
        _redeemed[msg.sender] += singleClaims.redeem(
            msg.sender,
            token,
            tokenIds
        );
        _redeemers.add(msg.sender);
    }

    /**
    @notice Returns the entire set of redeemers and respective number of
    redemptions.
    @dev This is used by both redeemMaxN and redeemFromSingle. Different
    deployments of this contract MUST be used to ensure hermetic tests.
     */
    function allRedeemed()
        public
        view
        returns (address[] memory redeemers, uint256[] memory numRedeemed)
    {
        uint256 n = _redeemers.length();
        redeemers = new address[](n);
        numRedeemed = new uint256[](n);

        for (uint256 i = 0; i < n; i++) {
            redeemers[i] = _redeemers.at(i);
            numRedeemed[i] = _redeemed[redeemers[i]];
        }
    }

    /**
    @notice Exposes the regular claimed() function for testing.
     */
    function claimed(uint256 tokenId) public view returns (uint256) {
        return claims.claimed(tokenId);
    }

    /**
    @notice Exposes the single-claim claimed() function for testing.
     */
    function singleClaimed(uint256 tokenId) public view returns (bool) {
        return singleClaims.claimed(tokenId);
    }
}
