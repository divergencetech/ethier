// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "raffle --n=[num_to_select] --entropy=[uncontrollable_seed_entropy]",
		Short: "Randomly selects n entrants from a given list of addresses in a deterministic but fair manner",
		Long: `The command binary randomly selects n entrants from a given list in a
deterministic but fair manner. Fairness is defined as an inability on the
part of the party running the binary to influence its outcome with
non-negligible probability.
		
Two sources of entropy are used: the list of entrants itself, and an
arbitrary hex seed which MUST be out of the control of the party running the
binary (e.g. specify a future Ethereum block number and then use its hex hash
without prefix). Each source is hashed with keccak256 and the resulting 8x
64-bit blocks are XORd to create an int64 seed for math/rand to shuffle the
entrants.
		
The list of addresses MUST be commited to before the entropy becomes known
otherwise the party running the binary can influence its outcome by adding
dummy addresses.
		
Assuming math/rand.Shuffle is adequately uniform in its selection of a
permutation then all fairness is based on the seed entropy.`,
	}

	// TODO(aschlosberg): investigate the idiomatic way of accessing cobra
	// flags, presumably via the Command passed to run.
	r := &raffle{
		n:       cmd.Flags().Int("n", 0, "Number of elements to select from stdin"),
		entropy: cmd.Flags().BytesHex("entropy", nil, "An arbitrary hexidecimal seed to influence the selection; this MUST be set to something out of your control, such as a future block's hash"),
	}
	cmd.RunE = r.run

	rootCmd.AddCommand(cmd)
}

type raffle struct {
	n       *int
	entropy *[]byte
}

func (r *raffle) run(c *cobra.Command, args []string) error {
	winners, err := choose(os.Stdin, *r.n, *r.entropy)
	if err != nil {
		return fmt.Errorf("choosing: %v", err)
	}
	for _, w := range winners {
		fmt.Println(w)
	}

	return nil
}

// choose reads newline-delimeted hexadecimal addresses from r and returns n of
// them at random, deterministically seeded.
func choose(from io.Reader, n int, entropy []byte) ([]common.Address, error) {
	addrs, err := readAddresses(from)
	if err != nil {
		return nil, fmt.Errorf("read addresses: %v", err)
	}
	if n > len(addrs) {
		return nil, fmt.Errorf("selecting %d from only %d entrant(s)", n, len(addrs))
	}

	s := seed(addrs, entropy)
	log.Printf("Seed: %d", s)
	rng := rand.New(rand.NewSource(s))

	rng.Shuffle(len(addrs), func(i, j int) {
		addrs[i], addrs[j] = addrs[j], addrs[i]
	})
	return addrs[:n], nil
}

// readAddresses reads newline-delimeted hexadecimal addresses from r, parses
// them, sorts according to bytes.Compare(), and returns the slice.
func readAddresses(r io.Reader) ([]common.Address, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll(): %v", err)
	}
	raw := strings.Split(strings.TrimSpace(string(buf)), "\n")

	var addrs []common.Address
	for _, r := range raw {
		r = strings.TrimSpace(r)
		if !common.IsHexAddress(r) {
			return nil, fmt.Errorf("invalid address %q", r)
		}
		addrs = append(addrs, common.HexToAddress(r))
	}
	sort.Slice(addrs, func(i, j int) bool {
		return bytes.Compare(addrs[i].Bytes(), addrs[j].Bytes()) == -1
	})

	return addrs, nil
}

// seed separately hashes (a) the bytes of all Addresses and (b) the entropy,
// and "folds" (with xor) the resulting output into a 64-bit value to be used as
// a seed to math/rand.Shuffle.
func seed(addrs []common.Address, entropy []byte) int64 {
	var addrBytes [][]byte
	for _, a := range addrs {
		addrBytes = append(addrBytes, a.Bytes())
	}

	// We need an int64 to seed math/rand so collapse the 4 words from each
	// entropy source with xor as a standard means of combining streams of
	// randomness.
	var result [8]byte
	for lbl, src := range map[string][]byte{
		"entropy":  crypto.Keccak256(entropy),
		"entrants": crypto.Keccak256(addrBytes...),
	} {
		log.Printf("Keccak256 of %s: %#x", lbl, src)

		if n := len(src); n != 32 {
			log.Fatalf("Entropy source from %s of length %d; must be 32", lbl, n)
		}
		var buf [32]byte
		copy(buf[:], src)

		for word := 0; word < 32; word += 8 {
			for i, b := range buf[word : word+8] {
				result[i] ^= b
			}
		}
	}

	return intFromBytes(&result)
}

// intFromBytes casts an 8-byte array as an int64.
func intFromBytes(x *[8]byte) int64 {
	return *(*int64)(unsafe.Pointer(x))
}
