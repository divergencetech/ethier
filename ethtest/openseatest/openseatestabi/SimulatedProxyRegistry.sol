// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "../../../contracts/thirdparty/opensea/ProxyRegistry.sol";

/**
@notice A minimal simulated OpenSea Wyvern proxy registry for use with ethier's
ethtest.SimulatedBackend Go testing.
 */
contract SimulatedProxyRegistry is ProxyRegistry {
    function setProxyFor(address owner, address proxy) public {
        proxies[owner] = OwnableDelegateProxy(proxy);
    }
}
