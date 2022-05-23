// The ethier binary is a CLI tool for the @divergencetech/ethier suite of
// Solidity contracts and Go packages for Ethereum development.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd.AddCommand(alphaCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Long: "The ethier binary is a CLI tool for the @divergencetech/ethier suite of Solidity contracts and Go packages for Ethereum development.",
}

var alphaCmd = &cobra.Command{
	Use:   "alpha",
	Short: "Commands currently in alpha stage, considered highly experimental.",
}
