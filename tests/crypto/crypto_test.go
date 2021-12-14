package crypto

import (
	"log"
	"os"
	"testing"

	"github.com/h-fam/errdiff"

	"github.com/ethereum/go-ethereum/common"

	"github.com/divergencetech/ethier/eth"
	"github.com/divergencetech/ethier/ethtest"
)

// As signing is a pure function we can reuse the same set of signers across all
// tests. See deploy() for their usage.
var (
	goodSigners     []*eth.Signer
	goodSignerAddrs []common.Address
	badSigner       *eth.Signer
)

func TestMain(m *testing.M) {
	var signers []*eth.Signer
	for i := 0; i < 3; i++ {
		s, err := eth.NewSigner(256)
		if err != nil {
			log.Fatalf("eth.NewSigner(256) error %v", err)
		}
		signers = append(signers, s)
	}

	goodSigners = signers[:2]
	badSigner = signers[2]

	for _, s := range goodSigners {
		goodSignerAddrs = append(goodSignerAddrs, s.Address())
	}

	os.Exit(m.Run())
}

// deploy deploys a new TestableSignatureChecker with goodSigners as the only
// allowable signers.
func deploy(t *testing.T) (*ethtest.SimulatedBackend, *TestableSignatureChecker) {
	t.Helper()
	sim := ethtest.NewSimulatedBackendTB(t, 3)

	_, _, checker, err := DeployTestableSignatureChecker(sim.Acc(0), sim, goodSignerAddrs)
	if err != nil {
		t.Fatalf("DeployTestableSignatureChecker() error %v", err)
	}
	return sim, checker
}

// A signatureTest is a test case common to both reusable and single-use
// signatures.
type signatureTest struct {
	name                 string
	signer               *eth.Signer
	signedData, sentData []byte
	errDiffAgainst       interface{}
}

const invalidSigMsg = "SignatureChecker: Invalid signature"

func signatureTestCases() []signatureTest {
	return []signatureTest{
		{
			name:           "valid",
			signer:         goodSigners[0],
			signedData:     []byte("hello"),
			sentData:       []byte("hello"),
			errDiffAgainst: nil,
		},
		{
			name:           "valid with different data",
			signer:         goodSigners[0],
			signedData:     []byte("world"),
			sentData:       []byte("world"),
			errDiffAgainst: nil,
		},
		{
			name:           "valid with different signer",
			signer:         goodSigners[1],
			signedData:     []byte("hello"),
			sentData:       []byte("hello"),
			errDiffAgainst: nil,
		},
		{
			name:           "incorrect data",
			signer:         goodSigners[0],
			signedData:     []byte("hello"),
			sentData:       []byte("Hello"),
			errDiffAgainst: invalidSigMsg,
		},
		{
			name:           "malicious signer",
			signer:         badSigner,
			signedData:     []byte("hello"),
			sentData:       []byte("Hello"),
			errDiffAgainst: invalidSigMsg,
		},
	}
}

func TestSingleUseSignature(t *testing.T) {
	sim, checker := deploy(t)

	for _, tt := range signatureTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			sig, nonce, err := tt.signer.SignWithNonce(tt.signedData)
			if err != nil {
				t.Fatalf("%T.SignWithNonce(%v) error %v", tt.signer, tt.signedData, err)
			}

			_, got := checker.NeedsSignature(sim.Acc(0), tt.sentData, nonce, sig)
			if diff := errdiff.Check(got, tt.errDiffAgainst); diff != "" {
				t.Errorf("NeedsSignature() on first call; %s", diff)
			}

			if tt.errDiffAgainst != nil {
				return
			}

			_, gotRepeat := checker.NeedsSignature(sim.Acc(0), tt.sentData, nonce, sig)
			if diff := errdiff.Check(gotRepeat, "SignatureChecker: Message already used"); diff != "" {
				t.Errorf("NeedsSignature() on second call with same message; %s", diff)
			}
		})
	}
}

func TestReusableSignature(t *testing.T) {
	_, checker := deploy(t)

	for _, tt := range signatureTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			sig, err := tt.signer.CompactSign(tt.signedData)
			if err != nil {
				t.Fatalf("%T.CompactSign(%v) error %v", tt.signer, tt.signedData, err)
			}

			for i := 0; i < 2; i++ {
				_, got := checker.NeedsReusableSignature(nil, tt.sentData, sig)
				if diff := errdiff.Check(got, tt.errDiffAgainst); diff != "" {
					t.Errorf("NeedsReusableSignature() call #%d; %s", i+1, diff)
				}
			}
		})
	}
}

func TestAddressSignature(t *testing.T) {
	sim, checker := deploy(t)

	const (
		alice int = iota
		bob
		eve
	)

	signAddr := func(signer *eth.Signer, party int) []byte {
		t.Helper()
		addr := sim.Addr(party)

		sig, err := signer.SignAddress(addr)
		if err != nil {
			t.Fatalf("%T.SignAddress(%v) error %v", signer, addr, err)
		}
		return sig
	}

	goodSigs := map[int][][]byte{
		alice: {
			signAddr(goodSigners[0], alice),
			signAddr(goodSigners[1], alice),
		},
		bob: {
			signAddr(goodSigners[0], bob),
			signAddr(goodSigners[1], bob),
		},
	}

	t.Run("valid signatures", func(t *testing.T) {
		for party, sigs := range goodSigs {
			for _, sig := range sigs {
				if _, err := checker.NeedsSenderSignature(sim.CallFrom(party), sig); err != nil {
					t.Errorf("NeedsSenderSignature() with valid signature; got err %v; want nil err", err)
				}
			}
		}
	})

	t.Run("stolen signatures", func(t *testing.T) {
		for _, sigs := range goodSigs {
			for _, sig := range sigs {
				_, got := checker.NeedsSenderSignature(sim.CallFrom(eve), sig)
				if diff := errdiff.Substring(got, invalidSigMsg); diff != "" {
					t.Errorf("NeedsSenderSignature() with stolen signature; %s", diff)
				}
			}
		}
	})

	t.Run("malicious signatures from bad signer", func(t *testing.T) {
		sig := signAddr(badSigner, eve)
		_, got := checker.NeedsSenderSignature(sim.CallFrom(eve), sig)
		if diff := errdiff.Substring(got, invalidSigMsg); diff != "" {
			t.Errorf("NeedsSenderSignature() with valid malicious signer; %s", diff)
		}
	})
}
