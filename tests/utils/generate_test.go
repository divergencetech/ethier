package utils_test

//go:generate sh -c "solc ../../contracts/utils/OwnerPausable.sol --base-path ../../ --include-path ../../node_modules --combined-json abi,bin | abigen --combined-json /dev/stdin --pkg utils_test --out generated_test.go"
