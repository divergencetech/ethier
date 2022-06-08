package main

import (
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
	signCmd := &cobra.Command{
		Use:   "sign",
		Short: "Signs messages from stdin using an ECDSA signer.",
	}

	rootCmd.AddCommand(signCmd)

	signAddrCmd := &cobra.Command{
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

// sign generates a new signer and signs a given message
func signAddresses(_ *cobra.Command, args []string) error {
	signer, err := eth.NewSigner(256)
	if err != nil {
		return fmt.Errorf("generate new signer: %v", err)
	}

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("read address on stdin: %v", err)
	}

	addresses := strings.Split(strings.TrimSpace(string(buf)), "\n")

	log.Printf("Signer: %v\n\n", signer)

	var signedAddresses []SignedAddress

	for _, address := range addresses {
		addr := common.HexToAddress(address)
		sig, err := signer.PersonalSignAddress(addr)
		if err != nil {
			return fmt.Errorf("signing address %v: %v", address, err)
		}
		signedAddresses = append(signedAddresses, SignedAddress{
			Address:   address,
			Signature: fmt.Sprintf("%#x", sig),
		})
	}

	buf, err = json.MarshalIndent(signedAddresses, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding json: %v", err)
	}

	fmt.Println(string(buf))

	return nil
}
