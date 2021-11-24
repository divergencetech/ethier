package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

// sign generates a new signer (if none is provided) and signs a given message
// TODO given signers
func signAddresses(_ *cobra.Command, args []string) (retErr error) {
	// pwd, err := os.Getwd()
	// if err != nil {
	// 	return fmt.Errorf("os.Getwd(): %v", err)
	// }

	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("signing: %w", retErr) // TODO What's %w
		}
	}()

	signer, err := eth.NewSigner(256)
	if err != nil {
		log.Fatalf("Generate signer: %v", err)
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
