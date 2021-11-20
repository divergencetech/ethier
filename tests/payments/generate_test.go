package payments_test

//go:generate sh -c "solc TestableDutchAuction.sol --base-path ../../ --include-path ../../node_modules --combined-json abi,bin | abigen --combined-json /dev/stdin --pkg payments_test --out generated_test.go"
