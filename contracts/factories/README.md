# Factories and Deployers

## Motivation

Deploying contracts on the Ethereum mainnet is expensive. Many contracts,
however, are identical in their code and only their parameters change. This was
the motivation behind [EIP-1167](https://eips.ethereum.org/EIPS/eip-1167)'s
Minimal Proxy Contract; TL;DR deploy a single "implementation" contract and a
series of "proxy" contracts that delegate all of their calls to the
implementaiton.

### Gas savings

All deployment savings are calculated using the geth SimulatedBackend, from
within ethier tests. All values are calculated at a gas price of 100 gwei and an
ETH price of $3,500.

* `PaymentSplitter`: 1.3M gas units; save >$450 per deployment

## Contracts

### Factories

The OpenZeppelin contracts all have `*Upgradable` equivalents that are amenable
to being used behind proxies. The ethier `./factories/` directory contains
factory contracts that deploy a single instance of the proxy-friendly contract
and then expose `deploy()` and `deployDeterministic()` methods to create
proxies.

There is no need to import any of the code in `./factories/`. To deploy a new
proxy manually, use the exposed functions in the verified Etherscan contracts
(see below). To automate deployment, use a _Deployer_ contract.

### Deployers

Each factory has an equivalent _Deployer_ library in `./contracts/factories/`.
The `instance()` function returns the current chain's factory by checking the
chain ID—no more passing different addresses to your constructor depending on
the network to which you're deploying!

Deploying a PaymentSplitter is as easy as:

```Solidity
pragma solidity >=0.8.0 <0.9.0;

import "@divergencetech/ethier/contracts/factories/PaymentSplitterDeployer.sol";

contract LookMaNoAddress {
    address payable revenues;

    constructor(address[] memory payees, uint256[] memory shares) {
        revenues = payable(
            PaymentSplitterDeployer.instance().deploy(payees, shares)
        );
    }
}
```

### Verified Factories

* `PaymentSplitter`
  * (Mainnet)[https://etherscan.io/address/0xf034d6a4b1a64f0e6038632d87746ca24b79d325#code]
  * (Rinkeby)[https://rinkeby.etherscan.io/address/0x633dc916D9f59cf4aA117dE2Bb8edF7752270EC0#code]