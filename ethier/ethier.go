// The ethier binary is a CLI tool for the @divergencetech/ethier suite of
// Solidity contracts and Go packages for Ethereum development.
package main

import "github.com/spf13/cobra"

func main() {
	rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Long: "The ethier binary is a CLI tool for the @divergencetech/ethier suite of Solidity contracts and Go packages for Ethereum development.",
}
