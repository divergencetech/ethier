package random_test

//go:generate sh -c "solc TestablePRNG.sol --base-path ../../ --include-path ../../node_modules --combined-json abi,bin | abigen --combined-json /dev/stdin --pkg random_test --out generated_test.go"
