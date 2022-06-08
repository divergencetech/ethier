// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/thirdparty/chainlink/Chainlink.sol";
import "../../../contracts/thirdparty/chainlink/VRFConsumerHelper.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
@notice A testable version of VRFConsumerHelper that (a) confirms plumbing of
functions, and (b) allows testing of the ethier chainlinktest VRF fake.
@dev Although the contract doesn't have to be Ownable, the exposed withdrawLINK
function may be copy-pasted as an example code, so ensure that it has the
onlyOwner modifier.
 */
contract TestableVRFConsumerHelper is Ownable, VRFConsumerHelper {
    /**
    @notice Values required by VRFConsumerBase and stored here to allow us to
    benchmark gas usage of Chainlink library vs regular approach of reading
    values from storage.
     */
    bytes32 public immutable keyHash;
    uint256 public immutable fee;

    constructor() {
        keyHash = Chainlink.vrfKeyHash();
        fee = Chainlink.vrfFee();
    }

    /// @notice All fulfilled randomness, available for Go testing.
    mapping(bytes32 => uint256) public randomness;

    /**
    @notice Required override from VRFConsumerBase.
    @dev This also acts as a test of the request-fulfill loop, requiring a
    matching request ID and that the randomess is keccak256(requestId).
     */
    function fulfillRandomness(bytes32 requestId, uint256 _randomness)
        internal
        override
    {
        require(requestId == lastRequestId, "Invalid request ID");
        require(
            _randomness == uint256(keccak256(abi.encodePacked(requestId))),
            "Invalid randomness"
        );
        randomness[requestId] = _randomness;
    }

    /// @notice Used to test for sequential IDs from chainlinktest Go package.
    bytes32 public lastRequestId;

    /**
    @notice Requests randomness from VRFConsumerBase using the VRFConsumerHelper
    for convenience.
     */
    function helperRequestRandomness() external {
        lastRequestId = super.requestRandomness();
    }

    /**
    @notice Requests randomness directly from VRFConsumerBase, bypassing the
    Chainlink library, and performs identical functionality to
    helperRequestRandomness() so we can compare gas usage.
     */
    function standardRequestRandomness() external {
        lastRequestId = VRFConsumerBase.requestRandomness(keyHash, fee);
    }

    /// @notice Exposes the internal _withdrawLINK function.
    function withdrawLINK(address recipient, uint256 amount)
        external
        onlyOwner
    {
        _withdrawLINK(recipient, amount);
    }
}
