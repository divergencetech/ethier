![testing](https://github.com/divergencetech/ethier/actions/workflows/test.yml/badge.svg)
![lint](https://github.com/divergencetech/ethier/actions/workflows/lint.yml/badge.svg)

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