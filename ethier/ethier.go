// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
//
// The ethier binary is a CLI tool for the @divergencetech/ethier suite of
// Solidity contracts and Go packages for Ethereum development.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Long: "The ethier binary is a CLI tool for the @divergencetech/ethier suite of Solidity contracts and Go packages for Ethereum development.",
}
