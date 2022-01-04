// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

/**
@notice Convenience library for Chainlink constants without pre-deployment
knowlege of the chain.
@dev Chain IDs:
 - Ethereum Mainnet 1
 - Rinkeby 4
 - Polygon 137
 - Mumbai 80001
 - geth's SimulatedBackend 1337 but only compatible if using ethier'
   chainlinktest package
 */
library Chainlink {
    /// @notice Returns the LINK token address for the current chain.
    function linkToken() internal view returns (address addr) {
        assembly {
            switch chainid()
            case 1 {
                addr := 0x514910771AF9Ca656af840dff83E8264EcF986CA
            }
            case 4 {
                addr := 0x01BE23585060835E02B77ef475b0Cc51aA1e0709
            }
            case 137 {
                addr := 0xb0897686c545045aFc77CF20eC7A532E3120E0F1
            }
            case 80001 {
                addr := 0x326C977E6efc84E512bB9C30f76E30c160eD06FB
            }
            case 1337 {
                // The geth SimulatedBackend iff used with the ethier
                // chainlinktest package.
                addr := 0x55B04d60213bcfdC383a6411CEff3f759aB366d6
            }
        }
    }

    /// @notice Returns the VRF coordinator address for the current chain.
    function vrfCoordinator() internal view returns (address addr) {
        assembly {
            switch chainid()
            case 1 {
                addr := 0xf0d54349aDdcf704F77AE15b96510dEA15cb7952
            }
            case 4 {
                addr := 0xb3dCcb4Cf7a26f6cf6B120Cf5A73875B7BBc655B
            }
            case 137 {
                addr := 0x3d2341ADb2D31f1c5530cDC622016af293177AE0
            }
            case 80001 {
                addr := 0x8C7382F9D8f56b33781fE506E897a4F1e2d17255
            }
            case 1337 {
                // The geth SimulatedBackend iff used with the ethier
                // chainlinktest package.
                addr := 0x5FfD760b2B48575f3869722cd816d8b3f94DDb48
            }
        }
    }

    /// @notice Returns the VRF key hash for the current chain.
    function vrfKeyHash() internal view returns (bytes32 keyHash) {
        assembly {
            switch chainid()
            case 1 {
                keyHash := 0xAA77729D3466CA35AE8D28B3BBAC7CC36A5031EFDC430821C02BC31A238AF445
            }
            case 4 {
                keyHash := 0x2ed0feb3e7fd2022120aa84fab1945545a9f2ffc9076fd6156fa96eaff4c1311
            }
            case 137 {
                keyHash := 0xf86195cf7690c55907b2b611ebb7343a6f649bff128701cc542f0569e2c549da
            }
            case 80001 {
                keyHash := 0x6e75b569a01ef56d18cab6a8e71e6600d6ce853834d4a5748b720d06f878b3a4
            }
            case 1337 {
                // The geth SimulatedBackend iff used with the ethier
                // chainlinktest package.
                keyHash := keccak256(0x1337, 2)
            }
        }
    }

    /**
    @notice Returns the VRF fee, in LINK denomination, for the current chain.
     */
    function vrfFee() internal view returns (uint256 fee) {
        uint256 chainId;
        assembly {
            chainId := chainid()
        }

        // All LINK implementations have 18 decimal places, just like ETH, so
        // we can use the Solidity suffix to ensure the correct multiplier
        // whilst still being readable.
        if (chainId == 1 || chainId == 1337) {
            // 1 = mainnet
            //
            // 1337 = The geth SimulatedBackend iff used with the ethier
            // chainlinktest package. The same as Ethereum Mainnet to enable
            // testing of this library.
            return 2 ether;
        }
        if (chainId == 137 || chainId == 80001) {
            // Polygon main- and test- nets
            return 0.0001 ether;
        }
        if (chainId == 4) {
            // Rinkeby
            return 0.1 ether;
        }
    }
}
