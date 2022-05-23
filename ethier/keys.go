package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/divergencetech/ethier/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/integration/gcpkms"
	"github.com/google/tink/go/tink"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	keys := &cobra.Command{
		Use:   "keys",
		Short: "Functionality for private keys and seed phrases of hot wallets",
		Long:  "NOTE: it is NOT secure to use this feature to back up recovery phrases for hardware wallets",
	}

	v := &keyVault{}
	manage := &cobra.Command{
		Use:   "manage",
		Short: "Manage a key-store of private keys",
		RunE:  v.initAndRun,
	}

	fs := manage.Flags()
	fs.String("db", "", "Path to file containing the key-store database; will be created if non-existent")

	// TODO(aschlosberg): add support for AWS and Hashicorp KMSs.
	fs.String("gcp_project", "", "GCP project in which encryption key(s) are managed")
	fs.String("gcp_location", "us-central1", "GCP location in which encryption key(s) are stored")
	fs.String("gcp_keyring", "", "GCP key ring containing encryption key(s)")
	fs.String("gcp_key", "", "GCP key, within the key ring, to use for envelope encryption of seed phrases")

	keys.AddCommand(manage)
	alphaCmd.AddCommand(keys)
}

// A keyVault implements functionality of the `keys manage` sub command.
type keyVault struct {
	db         *gorm.DB
	keyWrapper tink.AEAD

	seedPhrasePrompts func() (*seedPhrase, error)
}

// initAndRun initialises the keyVault with production values and returns run().
func (v *keyVault) initAndRun(cmd *cobra.Command, args []string) error {
	v.seedPhrasePrompts = seedPhrasePrompts

	dbFile, err := dbPath(cmd)
	if err != nil {
		return err
	}
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&EncryptedSeedPhrase{}); err != nil {
		return err
	}
	v.db = db

	kms, err := kmsKeyWrapper(cmd)
	if err != nil {
		return err
	}
	v.keyWrapper = kms

	return v.run()
}

// run runs the `keys manage` sub command.
func (v *keyVault) run() error {
	for {
		action := promptui.Select{
			Label: "Choose an action",
			Items: []string{
				"Decrypt existing seed phrase",
				"Add new seed phrase",
				"Exit",
			},
		}

		switch do, _, err := action.Run(); {
		case err == promptui.ErrInterrupt:
			// Exit gracefully
			return nil

		case err != nil:
			return err

		case do == 0:
			if err := v.showSeedPhrase(); err != nil && err != promptui.ErrInterrupt {
				return nil
			}

		case do == 1:
			if err := v.addSeedPhrase(); err != nil && err != promptui.ErrInterrupt {
				return err
			}

		default:
			return nil
		}
	}
}

func (v *keyVault) addSeedPhrase() error {
	phrase, err := v.seedPhrasePrompts()
	if err != nil {
		return err
	}

	ct, err := v.envelope().Encrypt([]byte(phrase.secret), phrase.addr.Bytes())
	if err != nil {
		return fmt.Errorf("Encrypt(): %v", err)
	}
	phrase.secret = ""

	return v.db.Create(&EncryptedSeedPhrase{
		Label:        phrase.label,
		AddressBytes: phrase.addr.Bytes(),
		Ciphertext:   ct,
	}).Error
}

func (v *keyVault) showSeedPhrase() error {
	var phrases []*EncryptedSeedPhrase
	if err := v.db.Model(&EncryptedSeedPhrase{}).Find(&phrases).Error; err != nil {
		return err
	}
	sort.Slice(phrases, func(i, j int) bool {
		return phrases[i].Label < phrases[j].Label
	})

	var items []string
	for _, p := range phrases {
		lbl := p.Label
		if lbl == "" {
			lbl = "?"
		}
		items = append(items, fmt.Sprintf("[%s] %v", p.Label, p.Address()))
	}

	s := &promptui.Select{
		Label: "Choose an address once nobody is looking at your monitor",
		Items: items,
	}
	show, _, err := s.Run()
	if err != nil {
		return err
	}

	p := phrases[show]
	pt, err := v.envelope().Decrypt(p.Ciphertext, p.AddressBytes)
	if err != nil {
		return err
	}
	// Displaying the plaintext as a Select means it gives the user
	// control over when to hide it.
	(&promptui.Select{
		Label: string(pt),
		Items: []string{"Done"},
	}).Run()
	return nil
}

// seedPhrase carries a SECRET seed (or recovery) phrase from which private keys
// are derived to sign Ethereum transactions.
type seedPhrase struct {
	// label is a human-readable value describing the seed phrase.
	label string
	// addr is the primary derived address from the HD path m/44'/60'/0'/0/ and
	// the 0 account (i.e. what MetaMask derives).
	addr common.Address
	// secret is the actual value, named as such to communicate that it's a
	// sensitive value.
	secret string
}

// EncryptedSeedPhrase is a GORM model for storing encrypted seedphrases in a
// sqlite database.
type EncryptedSeedPhrase struct {
	AddressBytes []byte `gorm:"primarykey"`
	Label        string
	Ciphertext   []byte
}

// Address parses and returns the raw address bytes.
func (p *EncryptedSeedPhrase) Address() common.Address {
	return common.BytesToAddress(p.AddressBytes)
}

// seedPhrasePrompts runs a series of Prompts to accept a seed phrase from the
// user, validate it, and confirm the primary derived address before returning
// the phrase.
func seedPhrasePrompts() (*seedPhrase, error) {
	show, err := boolPrompt("Show seed phrase while typing? Is your screen protected from everyone else's view")
	if err != nil {
		return nil, err
	}

	var signer *eth.Signer
	phraseP := &promptui.Prompt{
		Label: "Seed phrase",
		Validate: func(phrase string) error {
			s, err := eth.DefaultHDPathPrefix.SignerFromSeedPhrase(phrase, "", 0)
			if err != nil {
				return err
			}

			signer = s
			return nil
		},
	}
	if show {
		phraseP.HideEntered = true
	} else {
		phraseP.Mask = '*'
	}

	labelP := &promptui.Prompt{
		Label: "(Optional) Nickname for this address",
	}

	for {
		phrase, err := phraseP.Run()
		if err != nil {
			return nil, err
		}

		correct, err := boolPrompt("Confirm address %v", signer.Address())
		if err != nil {
			return nil, err
		}
		if !correct {
			continue
		}

		lbl, err := labelP.Run()
		if err != nil {
			return nil, err
		}

		return &seedPhrase{
			label:  lbl,
			secret: phrase,
			addr:   signer.Address(),
		}, nil
	}
}

// actualFile returns the absolute path to f after any symlinks are evaluated.
func dbPath(cmd *cobra.Command) (string, error) {
	f, err := getStringFlag(cmd, "db")
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(f)
	if err != nil {
		return "", fmt.Errorf("determine absolute path to %q: %v", f, err)
	}

	eval, err := filepath.EvalSymlinks(abs)
	switch {
	case err == nil:
		return eval, nil
	case os.IsNotExist(err):
		return abs, nil
	default:
		return "", fmt.Errorf("evaluate any symlinks in %q: %v", abs, err)
	}
}

// kmsKeyWrapper returns an AEAD backed by a cloud key-management service, for
// use in envelope encryption.
func kmsKeyWrapper(cmd *cobra.Command) (tink.AEAD, error) {
	proj, err := getStringFlag(cmd, "gcp_project")
	if err != nil {
		return nil, err
	}

	loc, err := getStringFlag(cmd, "gcp_location")
	if err != nil {
		return nil, err
	}

	ring, err := getStringFlag(cmd, "gcp_keyring")
	if err != nil {
		return nil, err
	}

	key, err := getStringFlag(cmd, "gcp_key")
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf(`gcp-kms://projects/%s/locations/%s/keyRings/%s/cryptoKeys/`, proj, loc, ring)

	kms, err := gcpkms.NewClient(prefix)
	if err != nil {
		return nil, fmt.Errorf("gcpkms.NewClient(%q): %v", prefix, err)
	}

	keyURI := prefix + key
	wrap, err := kms.GetAEAD(keyURI)
	if err != nil {
		return nil, fmt.Errorf("%T.GetAEAD(%q): %v", kms, keyURI, err)
	}
	return wrap, nil
}

// envelope returns a new AEAD for envelope encryption, where the wrapping key
// is provided by v.keyWrapper which, in practice, is returned by kmsKeyWrapper.
func (v *keyVault) envelope() tink.AEAD {
	return aead.NewKMSEnvelopeAEAD2(aead.AES256GCMKeyTemplate(), v.keyWrapper)
}
