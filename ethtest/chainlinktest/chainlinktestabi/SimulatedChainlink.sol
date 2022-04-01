// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/thirdparty/chainlink/Chainlink.sol";
import "@chainlink/contracts/src/v0.8/VRFConsumerBase.sol";
import "@chainlink/contracts/src/v0.8/VRFRequestIDBase.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/**
@notice Deploys simulated Chainlink assets for a VRF mock compatible with the
ethier chainlinktest Go package and the Chainlink.sol chain-agnostic library.
 */
contract SimulatedChainlink {
    address public immutable linkToken;
    SimulatedVRFCoordinator public immutable vrfCoordinator;

    constructor() {
        linkToken = address(new SimulatedLinkToken());
        vrfCoordinator = new SimulatedVRFCoordinator();
        vrfCoordinator.transferOwnership(msg.sender);

        require(
            linkToken == Chainlink.linkToken(),
            "SimulatedChainLink: unexpected LINK token address"
        );
        require(
            address(vrfCoordinator) == Chainlink.vrfCoordinator(),
            "SimulatedChainLink: unexpected VRFCoordinator address"
        );
    }
}

/// @notice A minimal contract to mock the LINK token for VRFConsumerBase.
contract SimulatedLinkToken is ERC20 {
    // solhint-disable-next-line no-empty-blocks
    constructor() ERC20("ChainLink Token", "LINK") {}

    /**
    @notice Emitted on a valid call to transferAndCall(); will be watched by the
    chainlinktest Go package, resulting in fulfilment in a different
    transaction.
     */
    event RandomnessRequested(address indexed by);

    function transferAndCall(
        address to,
        uint256 value,
        bytes calldata data
    ) external returns (bool success) {
        require(
            to == address(Chainlink.vrfCoordinator()),
            "Simulated LINK token: incorrect VRF Coordinator"
        );
        require(
            value == Chainlink.vrfFee(),
            "Simulated LINK token: incorrect fee for VRF"
        );
        require(
            keccak256(data) ==
                keccak256(abi.encodePacked(Chainlink.vrfKeyHash(), uint256(0))),
            "Simulated LINK token: invalid data"
        );

        transfer(to, value);

        emit RandomnessRequested(msg.sender);
        return true;
    }

    /// @notice Increases the recipient's balance by the specified amount.
    function faucet(address recipient, uint256 amount) external {
        _mint(recipient, amount);
    }
}

/**
@notice A minimal contract to mock a VRFCoordinator that works with
SimualtedLinkToken.
@dev Is Ownable to ensure that calls are from the chainlinktest Go package; akin
to a private function.
 */
contract SimulatedVRFCoordinator is Ownable, VRFRequestIDBase {
    mapping(address => uint256) private nonce;

    /**
    @notice Recreates the requestId computed by VRFConsumerBase, hashes it for
    simulated randomness, and fulfills with the original caller.
     */
    function fulfill(VRFConsumerBase consumer) public onlyOwner {
        address consumerAddr = address(consumer);

        uint256 vRFSeed = makeVRFInputSeed(
            Chainlink.vrfKeyHash(),
            0,
            consumerAddr,
            nonce[consumerAddr]
        );
        nonce[consumerAddr]++;

        bytes32 requestId = makeRequestId(Chainlink.vrfKeyHash(), vRFSeed);
        consumer.rawFulfillRandomness(
            requestId,
            uint256(keccak256(abi.encodePacked(requestId)))
        );
    }
}
