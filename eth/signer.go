package eth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	hdwallet "github.com/divergencetech/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/tink/go/prf"
)

// A Signer abstracts signing of arbitrary messages by wrapping an ECDSA private
// key and, optionally, its associated BIP39 mnemonic.
type Signer struct {
	key      *ecdsa.PrivateKey
	mnemonic string
}

// NewSigner is equivalent to
// DefaultHDPathPrefix.SignerFromSeedPhrase(NewMnemonic(), "", 0).
func NewSigner(bitSize int) (*Signer, error) {
	m, err := NewMnemonic(bitSize)
	if err != nil {
		return nil, err
	}
	return DefaultHDPathPrefix.SignerFromSeedPhrase(m, "", 0)
}

// NewMnemonic is a convenience wrapper around go-bip39 entropy and mnemonic
// creation.
func NewMnemonic(bitSize int) (string, error) {
	buf, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", fmt.Errorf("generate entropy: %v", err)
	}
	m, err := bip39.NewMnemonic(buf)
	if err != nil {
		return "", fmt.Errorf("convert entropy to mnemonic: %v", err)
	}
	return m, nil
}

// An HDPathPrefix is a prefix for use in deriving private keys from BIP39
// mnemonics. It is appended with the account number. Values MUST include a
// trailing slash.
type HDPathPrefix string

// DefaultHDPathPrefix is the default format for derived accounts when using
// SignerFromSeedPhrase().
const DefaultHDPathPrefix = HDPathPrefix("m/44'/60'/0'/0/")

// SignerFromSeedPhrase confirms that the mnemonic is valid under BIP39 and then
// uses it to derive a private key (see HDPathF)
func (hdp HDPathPrefix) SignerFromSeedPhrase(mnemonic, password string, account uint) (*Signer, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("create seed from mnemoic: %v", err)
	}
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return nil, fmt.Errorf("create wallet from seed: %v", err)
	}

	path, err := hdwallet.ParseDerivationPath(fmt.Sprintf("%s%d", hdp, account))
	if err != nil {
		return nil, fmt.Errorf("parse derivation path: %v", err)
	}
	acc, err := wallet.Derive(path, false)
	if err != nil {
		return nil, fmt.Errorf("derive account: %v", err)
	}

	key, err := wallet.PrivateKey(acc)
	if err != nil {
		return nil, fmt.Errorf("obtain private key: %v", err)
	}
	return &Signer{key, mnemonic}, nil
}

// SignerFromPRF deterministically derives a private key from the pseudo-random
// function and the input bytes. By definition, the output of a PRF is
// indistinguishable from a random function.
//
// The input parameter allows for different sets of HD wallets to be derived
// from the same underlying PRF key. Only the PRF key need be secret.
//
// SignerFromPRF can be thought of as a method for securely creating new
// mnemonic seed phrases from a single underlying key and different input
// parameters. Although the resulting mnemonic is accessible, SignerFromPRF is
// intended for use in an automated environment, which is why it relies on
// Google Tink.
func (hdp HDPathPrefix) SignerFromPRF(src prf.PRF, input []byte, account uint) (*Signer, error) {
	entropy, err := src.ComputePRF(input, 32)
	if err != nil {
		return nil, fmt.Errorf("compute entropy from PRF: %v", err)
	}
	mn, err := hdwallet.NewMnemonicFromEntropy(entropy)
	if err != nil {
		return nil, fmt.Errorf("derive mnemonic from entropy: %v", err)
	}
	return hdp.SignerFromSeedPhrase(mn, "", account)
}

// SignerFromPRFSet returns hdp.SifnerFromPRF() using the set's primary PRF.
// This is simply a convenience function as the prf package doesn't accomodate
// direct creation of a prf.PRF.
func (hdp HDPathPrefix) SignerFromPRFSet(set *prf.Set, input []byte, account uint) (*Signer, error) {
	return hdp.SignerFromPRF(set.PRFs[set.PrimaryID], input, account)
}

// String returns s.Address() as a string.
func (s *Signer) String() string {
	return s.Address().String()
}

// Mnemonic returns the mnemonic used to derive the Signer's private key. USE
// WITH CAUTION.
func (s *Signer) Mnemonic() string {
	return s.mnemonic
}

// Address returns the Signer's public key converted to an Ethereum address.
func (s *Signer) Address() common.Address {
	return crypto.PubkeyToAddress(s.key.PublicKey)
}

// CompactSignature converts a signature with the final byte, the y parity
// (always 0 or 1), carried in the highest bit of the s parameter, as per
// EIP-2098. Using compact signatures reduces gas by removing a word from
// calldata, and is compatible with OpenZeppelin's ECDSA.recover() helper.
func CompactSignature(rsv []byte) ([]byte, error) {
	// Convert the 65-byte signature returned by Sign() into a 64-byte
	// compressed version, as described in
	// https://eips.ethereum.org/EIPS/eip-2098.
	if n := len(rsv); n != 65 {
		return nil, fmt.Errorf("signature length %d; expecting 65", n)
	}
	v := rsv[64]
	if v != 0 && v != 1 {
		return nil, fmt.Errorf("signature V = %d; expecting 0 or 1", v)
	}
	rsv[32] |= v << 7
	return rsv[:64], nil
}

// AppendRandomNonce appends random 32 bytes to the buffer, commonly used in
// signature nonces.
func appendRandomNonce(buf []byte) ([]byte, [32]byte, error) {
	var nonce [32]byte
	if n, err := rand.Read(nonce[:]); n != 32 || err != nil {
		return nil, nonce, fmt.Errorf("read 32 random bytes: got %d bytes with err %v", n, err)
	}

	return append(buf, nonce[:]...), nonce, nil
}

// WithPersonalMessagePrefix converts a given message to conform to the signed data
// standard according to EIP-191.
func WithPersonalMessagePrefix(message []byte) []byte {
	prefix := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message)))
	return append(prefix, message...)
}

type signOpts struct {
	raw, compact, personal, withNonce bool
}

// sign signs a given buffer depending on the chosen options:
// withNonce = true, appends a nonce to the message
// compact = true, returns a compactified version of the signature according to
// EIP-2098.
// personal = true, adds a prefix to the message to conform to the EIP-191
// personal message standard.
// raw = false, the message is hashed before signing
func (s *Signer) sign(buf []byte, opts signOpts) ([]byte, *[32]byte, error) {
	var nonce *[32]byte
	var err error

	if opts.withNonce {
		var n [32]byte
		buf, n, err = appendRandomNonce(buf)
		nonce = &n
		if err != nil {
			return nil, nil, err
		}
	}

	if opts.personal {
		buf = WithPersonalMessagePrefix(buf)
	}

	if !opts.raw {
		buf = crypto.Keccak256(buf)
	}

	sig, err := crypto.Sign(buf, s.key)
	if err != nil {
		return nil, nil, err
	}

	if !opts.compact {
		return sig, nonce, nil
	}

	sig, err = CompactSignature(sig)
	if err != nil {
		return nil, nil, err
	}

	return sig, nonce, nil
}

// RawSign returns an ECDSA signature of buf. USE WITH CAUTION as signed data
// SHOULD be hashed first to avoid chosen-plaintext attacks. Prefer
// Signer.Sign().
func (s *Signer) RawSign(buf []byte) ([]byte, error) {
	sig, _, err := s.sign(buf, signOpts{
		raw:       true,
		compact:   false,
		personal:  false,
		withNonce: false,
	})
	return sig, err
}

// Sign returns an ECDSA signature of keccak256(buf).
func (s *Signer) Sign(buf []byte) ([]byte, error) {
	sig, _, err := s.sign(buf, signOpts{
		raw:       false,
		compact:   false,
		personal:  false,
		withNonce: false,
	})
	return sig, err
}

// PersonalSign returns an EIP-191 conform personal ECDSA signature of buf
// Convenience wrapper for s.CompactSign(WithPersonalMessagePrefix(buf))
func (s *Signer) PersonalSign(buf []byte) ([]byte, error) {
	sig, _, err := s.sign(buf, signOpts{
		raw:       false,
		compact:   true,
		personal:  true,
		withNonce: false,
	})
	return sig, err
}

// PersonalSignWithNonce generates a 32-byte nonce with crypto/rand and returns
// s.PersonalSign(append(buf, nonce)).
func (s *Signer) PersonalSignWithNonce(buf []byte) ([]byte, [32]byte, error) {
	sig, nonce, err := s.sign(buf, signOpts{
		raw:       false,
		compact:   true,
		personal:  true,
		withNonce: true,
	})
	if err != nil {
		return nil, [32]byte{}, err
	}
	return sig, *nonce, err

}

// SignAddress is a convenience wrapper for s.PersonalSign(addr.Bytes()).
func (s *Signer) PersonalSignAddress(addr common.Address) ([]byte, error) {
	return s.PersonalSign(addr.Bytes())
}

// TransactorWithChainID returns bind.NewKeyedTransactorWithChainID(<key>,
// chainID) where <key> is the Signer's private key.
func (s *Signer) TransactorWithChainID(chainID *big.Int) (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(s.key, chainID)
}
