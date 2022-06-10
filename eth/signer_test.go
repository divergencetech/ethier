package eth_test

import (
	"bytes"
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/prf"
	"github.com/google/tink/go/tink"
	"github.com/h-fam/errdiff"

	// These tests require ethtest.SimulatedBackend but that would result in a
	// cyclical dependency. As this is limited to these tests and not the
	// package itself, we simply move them from package eth to package eth_test
	// and dot-import to avoid having to qualify the full package name. This
	// MUST NOT be considered precedent outside of tests and SHOULD be avoided
	// where possible.
	. "github.com/divergencetech/ethier/eth"
)

const (
	testOnlyPRFKey0 = `{"encryptedKeyset":"CPC2o+0NEokBCn0KMXR5cGUuZ29vZ2xlYXBpcy5jb20vZ29vZ2xlLmNyeXB0by50aW5rLkhtYWNQcmZLZXkSRhICCAQaQCNsZSUAD/nuSe3o/PKoljLf1Yg4BSqLuMGcEpF9w9ISXUYR1lu1H96n6q5Ixzzuuy6oEZW873QmdL9xNhiGaskYARABGPC2o+0NIAM=","keysetInfo":{"primaryKeyId":3718830960,"keyInfo":[{"typeUrl":"type.googleapis.com/google.crypto.tink.HmacPrfKey","status":"ENABLED","keyId":3718830960,"outputPrefixType":"RAW"}]}}`
	testOnlyPRFKey1 = `{"encryptedKeyset":"CJqX8oMLEokBCn0KMXR5cGUuZ29vZ2xlYXBpcy5jb20vZ29vZ2xlLmNyeXB0by50aW5rLkhtYWNQcmZLZXkSRhICCAQaQGpncr+9+ZVltjXAgkW7eQ9zQCPNz7Wxcjg0AGoFSu9SR3B2iSH/efkKCD5fRhEyznkfrA15Y6Fr3J7yllqB+98YARABGJqX8oMLIAM=","keysetInfo":{"primaryKeyId":2960952218,"keyInfo":[{"typeUrl":"type.googleapis.com/google.crypto.tink.HmacPrfKey","status":"ENABLED","keyId":2960952218,"outputPrefixType":"RAW"}]}}`
)

func TestDeterministicSigner(t *testing.T) {
	// This is not a test of correctness but rather a means of locking in
	// deterministic outputs to avoid regression. It also acts as a
	// demonstration of the functionality as different parameters are changed.
	//
	// See the Tink documentation (https://github.com/google/tink) for
	// information regarding key management.

	type address struct {
		input        []byte
		account      uint
		want         common.Address
		wantMnemonic string
	}

	// This is used to demonstrate that the mnemonic is determined by the
	// key:input combination, after which the approach functions identically to
	// standard HD Wallets. See comment on SignerFromPRF() method re mnemonics.
	const key0NilInputMnemonic = "hybrid fox hover between identify only taste this cliff denial main buffalo slide start dirt diary version thumb remain aim hybrid uncle grit pony"

	tests := []struct {
		jsonKeySet string
		addrs      []address
	}{
		{
			jsonKeySet: testOnlyPRFKey0,
			addrs: []address{
				{
					input:        nil,
					account:      0,
					want:         common.HexToAddress("0x71e059FA4594b69200541A189010188eDFFbC34D"),
					wantMnemonic: key0NilInputMnemonic,
				},
				{
					input:        nil,
					account:      1,
					want:         common.HexToAddress("0x8000B0045d0Ce1265d74FFCF60d4311c565C140B"),
					wantMnemonic: key0NilInputMnemonic,
				},
				{
					input:        []byte("hello"),
					account:      0,
					want:         common.HexToAddress("0xaA358Da5C7f1D5EbA1d2a54ea03EB8961aDE39a1"),
					wantMnemonic: "note outdoor column left narrow fresh curious orphan jelly similar cross bulb curious sting blush provide regret crisp harbor canvas bamboo secret swing uphold",
				},
			},
		},
		{
			jsonKeySet: testOnlyPRFKey1,
			addrs: []address{
				{
					input:        nil,
					account:      0,
					want:         common.HexToAddress("0x1DAb98Dc62bB2ED1089Eeb4aDfC7011342Cd4764"),
					wantMnemonic: "snow stereo deal adapt hybrid ritual nest normal budget lawn cactus sting curious wise soap bus pudding order annual siege pave surface occur saddle",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			r := keyset.NewJSONReader(strings.NewReader(tt.jsonKeySet))
			kh, err := keyset.Read(r, nonSecureAEADOnlyForTesting{})
			if err != nil {
				t.Fatalf("keyset.Read(%s): %v", tt.jsonKeySet, err)
			}
			set, err := prf.NewPRFSet(kh)
			if err != nil {
				t.Fatalf("prf.NewPRFSet() error %v", err)
			}

			for _, addr := range tt.addrs {
				sig, err := DefaultHDPathPrefix.SignerFromPRFSet(set, addr.input, addr.account)
				if err != nil {
					t.Fatalf("SignerFromPRFSet() error %v", err)
				}

				if got := sig.Address(); !bytes.Equal(got.Bytes(), addr.want.Bytes()) {
					t.Errorf("%s.SignerFromPRFSet([set], %q, %d) got address %v; want %v", DefaultHDPathPrefix, addr.input, addr.account, got, addr.want)
				}

				if got := sig.Mnemonic(); got != addr.wantMnemonic {
					t.Errorf("%s.SignerFromPRFSet([set], %q, %d) got mnemonic %q; want %q", DefaultHDPathPrefix, addr.input, addr.account, got, addr.wantMnemonic)
				}
			}
		})
	}
}

type nonSecureAEADOnlyForTesting struct{}

var _ tink.AEAD = nonSecureAEADOnlyForTesting{}

func (nonSecureAEADOnlyForTesting) Encrypt(plaintext, _ []byte) ([]byte, error) {
	return plaintext, nil
}

func (nonSecureAEADOnlyForTesting) Decrypt(ciphertext, _ []byte) ([]byte, error) {
	return ciphertext, nil
}

func TestTransactorWithChainID(t *testing.T) {
	ctx := context.Background()

	// Test that Signer.TransactorWithChainID works by using the resulting TransactOpts
	// to send funds on the SimulatedBackend.

	sim := ethtest.NewSimulatedBackendTB(t, 1)
	signer, err := NewSigner(256)
	if err != nil {
		t.Fatalf("NewSigner(256) error %v", err)
	}
	t.Logf("Faucet: %v", sim.Addr(0))
	t.Logf("Signer under test: %v", signer.Address())

	wantBalance := func(ctx context.Context, t *testing.T, desc string, addr common.Address, want *big.Int) {
		t.Helper()
		t.Run(desc, func(t *testing.T) {
			got := sim.BalanceOf(ctx, t, addr)
			if got.Cmp(want) != 0 {
				t.Errorf("sim.BalanceOf(%v) got %d; want %d", addr, got, want)
			}
		})
	}
	wantBalance(ctx, t, "faucet at beginning", sim.Addr(0), Ether(100)) // default of SimulatedBackend
	wantBalance(ctx, t, "signer at beginning", signer.Address(), big.NewInt(0))

	gasPrice := EtherFraction(1, 1e9)
	const gasLimit = 21000
	txFee := new(big.Int).Mul(gasPrice, big.NewInt(gasLimit))

	sendEth := func(t *testing.T, opts *bind.TransactOpts, to common.Address, value *big.Int, errDiffAgainst interface{}) {
		t.Helper()
		unsigned := types.NewTransaction(0, to, value, gasLimit, gasPrice, nil)
		tx, err := opts.Signer(opts.From, unsigned)
		if err != nil {
			t.Fatalf("%T.Signer(%+v) error %v", opts, unsigned, err)
		}
		if diff := errdiff.Check(sim.SendTransaction(ctx, tx), errDiffAgainst); diff != "" {
			t.Fatalf("%T.SendTransaction() %s", sim, diff)
		}
	}

	sendEth(t, sim.Acc(0), signer.Address(), Ether(42), nil)
	wantBalance(ctx, t, "faucet after sending 42", sim.Addr(0), new(big.Int).Sub(Ether(100-42), txFee))
	wantBalance(ctx, t, "signer after receiving 42", signer.Address(), Ether(42))

	t.Run("correct chain ID", func(t *testing.T) {
		chainID := sim.Blockchain().Config().ChainID
		opts, err := signer.TransactorWithChainID(chainID)
		if err != nil {
			t.Fatalf("%T.TransactorWithChainID(%d) error %v", signer, chainID, err)
		}
		sendEth(t, opts, sim.Addr(0), Ether(21), nil)
		wantBalance(ctx, t, "faucet after sending 42 and receiving 21", sim.Addr(0), new(big.Int).Sub(Ether(100-42+21), txFee))
		wantBalance(ctx, t, "signer after receiving 42 and sending 21", signer.Address(), new(big.Int).Sub(Ether(42-21), txFee))
	})

	t.Run("incorrect chain ID", func(t *testing.T) {
		chainID := new(big.Int).Add(sim.Blockchain().Config().ChainID, big.NewInt(1))
		opts, err := signer.TransactorWithChainID(chainID)
		if err != nil {
			t.Fatalf("%T.TransactorWithChainID(%d) error %v", signer, chainID, err)
		}
		sendEth(t, opts, sim.Addr(0), Ether(1), "invalid chain id")
	})
}
