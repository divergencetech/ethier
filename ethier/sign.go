package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/divergencetech/ethier/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	var signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Signs messages from stdin using an ECDSA signer.",
	}

	rootCmd.AddCommand(signCmd)
	signCmd.PersistentFlags().String("mnemonic", "", "Path to file containing the mnemonic, password and account number (newline separated) for key derivation")

	var signAddrCmd = &cobra.Command{
		Use:   "addresses",
		Short: "Signs addresses from stdin using an ECDSA signer.",
		RunE:  signAddresses,
	}

	signCmd.AddCommand(signAddrCmd)
}

type SignedAddress struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

func getSignerFromMnemonic(filepath string) (*eth.Signer, error) {
	dat, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Read mnemonic: %v", err)
	}

	lines := strings.Split(string(dat), "\n")

	mnemonic := lines[0]
	var password string = ""
	var account uint = 0

	if len(lines) > 1 {
		password = lines[1]
	}

	if len(lines) > 2 {
		if lines[2] != "" {
			acc, err := strconv.Atoi(lines[2])
			if err != nil {
				return nil, fmt.Errorf("Converting derivation account: %v", err)
			}
			account = uint(acc)
		}
	}

	signer, err := eth.DefaultHDPathPrefix.SignerFromSeedPhrase(
		mnemonic, password, account,
	)

	if err != nil {
		return nil, fmt.Errorf("Generate signer: %v", err)
	}
	return signer, nil
}

// sign addresses generates a new signer or derives it from a given mnemonic
// to sign a list of addresses
func signAddresses(cmd *cobra.Command, args []string) (retErr error) {
	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("signing: %w", retErr)
		}
	}()

	mnemonicPath, err := cmd.Flags().GetString("mnemonic")
	if err != nil {
		log.Fatalf("Get mnemonics flag: %v", err)
	}

	var signer *eth.Signer

	if mnemonicPath != "" {
		signer, err = getSignerFromMnemonic(mnemonicPath)
		if err != nil {
			log.Fatalf("Generate signer from mnemonic: %v", err)
		}
	} else {
		signer, err = eth.NewSigner(256)
		if err != nil {
			log.Fatalf("Generate new signer: %v", err)
		}
	}

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Read stdin: %v", err)
	}

	addresses := strings.Split(strings.TrimSpace(string(buf)), "\n")

	log.Printf("Signer: %v\n\n", signer)

	var signedAddresses []SignedAddress

	for _, address := range addresses {
		addr := common.HexToAddress(address)
		sig, err := signer.SignAddress(addr)
		if err != nil {
			log.Fatalf("Signing address %v: %v", address, err)
		}
		signedAddresses = append(signedAddresses, SignedAddress{
			Address:   address,
			Signature: "0x" + hex.EncodeToString(sig),
		})
	}

	json_, err := json.MarshalIndent(signedAddresses, "", "  ")
	if err != nil {
		log.Fatalf("Encoding json: %v", err)
	}

	fmt.Println(string(json_))

	return nil
}
