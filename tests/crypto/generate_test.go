package crypto_test

//go:generate sh -c "solc TestableSignatureChecker.sol --base-path ../../ --include-path ../../node_modules --combined-json abi,bin | abigen --combined-json /dev/stdin --pkg crypto_test --out generated_test.go"
