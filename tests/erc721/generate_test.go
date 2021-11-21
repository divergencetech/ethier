package erc721_test

//go:generate sh -c "solc TestableERC721Common.sol --base-path ../../ --include-path ../../node_modules --combined-json abi,bin | abigen --combined-json /dev/stdin --pkg erc721_test --out generated_test.go"
