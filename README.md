![testing](https://github.com/divergencetech/ethier/actions/workflows/test.yml/badge.svg)
![linting](https://github.com/divergencetech/ethier/actions/workflows/lint.yml/badge.svg)

# Motivation

_ethier_ (pronounced "easier" with a lisp) intends to:

1. Gradually replace the reliance on JavaScript in Ethereum development with Go
   as it is (a) faster due to in-process backends for testing, and (b) more
   robust due to type safety. Although unlikely, ethier's "North Star" is a
   replacement for Truffle/Hardhat.
2. Provide reusable Solidity functionality not covered by OpenZeppelin and,
   where appropriate, provide respective Go bindings with round-trip testing.

## Versioning, stability, and production readiness

ethier uses [Semantic Versioning 2.0.0](https://semver.org). As the major
version is currently zero, the _API is open to change without warning_.

Contracts are very thoroughly tested but have not been subject to audit nor
widespread use. Early adopters are not only welcome, but will be greatly
appreciated.

## Why NPM if we're moving away from JavaScript?

Although ethier intends to use Go as much as possible, users may not, and NPM
is the de facto standard in Ethereum development. While this gives us a weird
mashup of go.mod and package.json, it's fit for purpose.

# Getting started

## Installation

1. Assuming `solc` and `go` are already installed:
```
go install github.com/divergencetech/ethier/ethier@latest
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```
2. Ensure that the `go/bin` directory is in your `$PATH`. This can be confirmed
by running `which ethier && echo GOOD`; if the word `GOOD` is printed then the
`ethier` binary has been found.

## Usage

### Generating Go bindings

This example assumes that the file is in the `contracts` directory with all
Solidity files present. If moving to a different directory, simply change the
relative paths.

```Go
package contracts

//go:generate ethier gen MyContract.sol
```

Run `go generate ./...` to generate Go bindings, including deployment functions.
The above example will generate `contracts/generated.go` with bindings to
`MyContract.sol`. These bindings can be used for (1) testing and/or (2) connecting to a
gateway node (e.g. Infura or Alchemy), depending on the
[`ContractBackend`](https://pkg.go.dev/github.com/ethereum/go-ethereum/accounts/abi/bind#ContractBackend)
being used:

1. For tests, use ethier's
[`ethtest.SimulatedBackend`](https://pkg.go.dev/github.com/divergencetech/ethier/ethtest#NewSimulatedBackendTB),
which extends the standard geth simulator with convenience behaviour like auto
commitment of transactions.

2. For a gateway, use the
[`ethclient`](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient)
package.

### Example test

```Go
package contracts

import (
   "testing" 

   "github.com/divergencetech/ethier/ethtest"
)

// The test backend creates as many accounts as needed, each representing a different
// "actor" in the test scenarios. A useful pattern is to simply enumerate them the iota
// pattern (which automatically increments) and add a `numAccounts` at the end.
const (
   deployer = iota
   vandal
   numAccounts
)

func TestMyContract(t *testing.T){
   sim := ethtest.NewSimulatedBackend(t, numAccounts)

   // The DeployMyContract function is automatically generated when running `go generate`.
   // addr and tx generally aren't useful, but are documented here for completeness
   addr, tx, contract, err := DeployMyContract(sim.Acc(deployer), sim /*,,, [constructor arguments]*/)
   if err != nil {
      t.Fatalf("DeployMyContract(%v) error %v", …, err)
   }

   // NOTE: If connecting to a deployed contract above, use NewMyContract() and substitute `sim`
   // for an *ethclient.Client`.

   t.Run("protect something sensitive", func(t *testing.T){
      // The test-actor pattern in the consts above makes tests self-documenting.
      _, err := contract.DoSomethingImportant(sim.Acc(vandal))
      // Confirm that there's an error because the vandal shouldn't be allowed to do anything
      // important!!! See the ethtest/revert package.
   })
}
```

See `tests/` for further usage examples. Remember to add `generated.go` to your
`.gitignore` file.