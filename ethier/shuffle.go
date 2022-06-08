package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func init() {
	const short = "Reads lines from stdin and shuffles them in a verifiable manner; useful for allow-list selection or metadata shuffling."

	cmd := &cobra.Command{
		Use:   "shuffle",
		Short: short,
		Long: short + `

By committing to the input data and an entropy source out of one's control, shuffle is both transparent and deterministic so its results can be verified by a third party.

The source of entropy can either be verifiably random, or sourced via a secondary commitment such as the hash of an Ethereum block in the future.`,
		RunE: shuffle,
	}

	cmd.Flags().BytesHexP("entropy", "e", nil, "Hexadecimal source of entropy to control shuffling")
	cmd.Flags().IntP("number", "n", 0, "Output first n values; 0 = all")

	rootCmd.AddCommand(cmd)
}

// shuffle implements the `ethier shuffle` command.
func shuffle(cmd *cobra.Command, args []string) error {
	selectN, err := cmd.Flags().GetInt("number")
	if err != nil {
		return err
	}

	seed, err := externalEntropy(cmd)
	if err != nil {
		return err
	}

	lines, err := sortedNonEmpty(os.Stdin)
	if err != nil {
		return err
	}
	seed.hashAndFold(lines...)

	seed.rand().Shuffle(len(lines), func(i, j int) {
		lines[i], lines[j] = lines[j], lines[i]
	})

	k := len(lines)
	if selectN == 0 || selectN > k {
		selectN = k
	}
	log.Printf("Selecting %d of %d", selectN, k)
	fmt.Printf("%s\n", bytes.Join(lines[:selectN], []byte("\n")))

	return nil
}

// externalEntropy returns a new entropy collector, initially seeded with the
// Command's --entropy flag.
func externalEntropy(cmd *cobra.Command) (*entropy, error) {
	buf, err := cmd.Flags().GetBytesHex("entropy")
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, errors.New("--entropy flag not specified")
	}

	log.Printf("External entropy: %#x", buf)

	e := new(entropy)
	e.hashAndFold(buf)
	return e, nil
}

// sortedNonEmpty reads all of r, splits the data by \n, trims surrounding
// whitespace of each line, removes empty lines, and returns all remaining
// values sorted with bytes.Compare().
func sortedNonEmpty(r io.Reader) ([][]byte, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read input: %v", err)
	}

	var lines [][]byte
	for _, l := range bytes.Split(buf, []byte("\n")) {
		l = bytes.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		lines = append(lines, l)
	}
	sort.Slice(lines, func(i, j int) bool {
		return bytes.Compare(lines[i], lines[j]) == -1
	})

	return lines, nil
}

// entropy collects sources of entropy by xor-folding 8-byte words into the
// existing value; it can then be used to create a seeded rand.Rand.
type entropy uint64

// hashAndFold hashes data with Keccak256 and folds it into the existing entropy
// value.
func (e *entropy) hashAndFold(data ...[]byte) {
	e.fold(crypto.Keccak256(data...))
}

// fold pads buf to a multiple of 8 bytes, treats each 8-byte word as a
// big-endian uint64, and sets e to the xor of all words and of the existing
// entropy.
func (e *entropy) fold(buf []byte) {
	buf = append(buf, make([]byte, 8-len(buf)%8)...)
	for i, n := 0, len(buf); i < n; i += 8 {
		*e ^= entropy(binary.BigEndian.Uint64(buf[i : i+8]))
	}
}

// rand returns a new rand.Rand with e as it's Source seed.
func (e *entropy) rand() *rand.Rand {
	return rand.New(rand.NewSource(int64(*e)))
}
