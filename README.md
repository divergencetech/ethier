# Motivation

_ethier_ (pronounced "easier" with a lisp) intends to:

1. Gradually replace the reliance on JavaScript in Ethereum development with Go
   as it is (a) faster due to in-process backends for testing, and (b) more
   robust due to type safety. Although unlikely, ethier's "North Star" is a
   replacement for Truffle/Hardhat.
2. Provide reusable Solidity functionality not covered by OpenZeppelin and,
   where appropriate, provide respective Go bindings with round-trip testing.
