// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gen",
		Short: "Compiles Solidity contracts to generate Go ABI bindings with go:generate",
		RunE:  gen,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("no source files provided")
			}
			for _, a := range args {
				if !strings.HasSuffix(a, ".sol") {
					return fmt.Errorf("non-Solidity file %q", a)
				}
			}
			return nil
		},
	})
}

// gen runs `solc | abigen` on the Solidity source files passed as the args.
// TODO: support wildcard / glob matching of files.
func gen(_ *cobra.Command, args []string) (retErr error) {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd(): %v", err)
	}
	// The Go package for abigen.
	pkg := filepath.Base(pwd)
	log.Printf("Generating package %q: %s", pkg, args)

	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("generating %q: %w", pkg, retErr)
		}
	}()

	// solc requires a base-path within which absolute includes are found. We
	// define this as the base path of the Go module.
	basePath := pwd
	for ; ; basePath = filepath.Join(basePath, "..") {
		if _, err := os.Stat(filepath.Join(basePath, "go.mod")); !errors.Is(err, os.ErrNotExist) {
			break
		}
	}

	args = append(
		args,
		"--base-path", basePath,
		"--include-path", filepath.Join(basePath, "node_modules"),
		"--combined-json", "abi,bin",
	)
	solc := exec.Command("solc", args...)
	solc.Stderr = os.Stderr

	// TODO: use bind.Bind() directly, instead of piping to abigen, which
	// requires that it's installed and within PATH. Blocked by
	// https://github.com/ethereum/go-ethereum/issues/23939 for which we've
	// submitted a fix.
	abigen := exec.Command(
		"abigen",
		"--combined-json", "/dev/stdin",
		"--pkg", pkg,
		"--out", "generated.go",
	)
	abigen.Stderr = os.Stderr

	r, w := io.Pipe()
	solc.Stdout = w
	abigen.Stdin = r

	if err := solc.Start(); err != nil {
		return fmt.Errorf("start `solc`: %v", err)
	}
	if err := abigen.Start(); err != nil {
		return fmt.Errorf("start `abigen`: %v", err)
	}
	if err := solc.Wait(); err != nil {
		w.Close()
		return fmt.Errorf("`solc` returned: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close write-half of pipe from solc to abigen: %v", err)
	}
	if err := abigen.Wait(); err != nil {
		return fmt.Errorf("`abigen` returned: %v", err)
	}
	return r.Close()
}
