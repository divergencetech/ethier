// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
//
// The raffle binary randomly selects n entrants from a given list in a
// deterministic but fair manner. Fairness is defined as an inability on the
// part of the party running the binary to influence its outcome with
// non-negligible probability.
//
// Two sources of entropy are used: the list of entrants itself, and an
// arbitrary hex seed which MUST be out of the control of the party running the
// binary (e.g. specify a future Ethereum block number and then use its hex hash
// without prefix). Each source is hashed with keccak256 and the resulting 8x
// 64-bit blocks are XORd to create an int64 seed for math/rand to shuffle the
// entrants.
//
// The list of addresses MUST be commited to before the entropy becomes known
// otherwise the party running the binary can influence its outcome by adding
// dummy addresses.
//
// Assuming math/rand.Shuffle is adequately uniform in its selection of a
// permutation then all fairness is based on the seed entropy.
//
// Usage, where the file "entrants" is new-line delimited:
// $ < entrants raffle --n=[num_to_select] --entropy=[uncontrollable_seed_entropy]
package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strings"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
)

func main() {
	n := flag.Int("n", 0, "Number of elements to select from stdin")
	rawEntropy := flag.String("entropy", "", "An arbitrary hexidecimal seed to influence the selection; this MUST be set to something out of your control, such as a future block's hash")
	flag.Parse()

	ent, err := hex.DecodeString(*rawEntropy)
	if err != nil {
		glog.Exitf("Decode hex entropy %q: %v", *rawEntropy, err)
	}
	if len(ent) == 0 {
		glog.Exit("Empty entropy buffer")
	}
	glog.Infof("Entropy: %#x", ent)

	winners, err := choose(os.Stdin, *n, ent)
	if err != nil {
		glog.Exitf("Choosing: %v", err)
	}
	for _, w := range winners {
		fmt.Println(w)
	}
}

// choose reads newline-delimeted hexadecimal addresses from r and returns n of
// them at random, deterministically seeded.
func choose(r io.Reader, n int, entropy []byte) ([]common.Address, error) {
	addrs, err := readAddresses(r)
	if err != nil {
		return nil, fmt.Errorf("read addresses: %v", err)
	}
	if n > len(addrs) {
		return nil, fmt.Errorf("selecting %d from only %d entrant(s)", n, len(addrs))
	}

	s := seed(addrs, entropy)
	glog.Infof("Seed: %d", s)
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
		glog.Infof("Keccak256 of %s: %#x", lbl, src)

		if n := len(src); n != 32 {
			glog.Fatalf("Entropy source from %s of length %d; must be 32", lbl, n)
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
