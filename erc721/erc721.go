// Package erc721 provides functionality associated with ERC721 NFTs.
package erc721

import (
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
)

// A TokenID is a uint256 Solidty tokenId.
type TokenID [32]byte

// TokenIDFromUint256 returns the 32-byte buffer underlying u, typed as a
// TokenID pointer.
func TokenIDFromUint256(u *uint256.Int) *TokenID {
	t := TokenID(u.Bytes32())
	return &t
}

// TokenIDFromBig is a convenience wrapper for
// TokenIDFromUint256(uint256.FromBig(b)), returning an error if b overflows.
func TokenIDFromBig(b *big.Int) (*TokenID, error) {
	u, overflow := uint256.FromBig(b)
	if overflow {
		return nil, fmt.Errorf("%T %v overflows uint256", b, b)
	}
	return TokenIDFromUint256(u), nil
}

// TokenIDFromHex is a convenience wrapper for
// TokenIDFromUint256(uint256.FromHex(b)), returning an error if FromHex() does.
func TokenIDFromHex(h string) (*TokenID, error) {
	u, err := uint256.FromHex(h)
	if err != nil {
		return nil, err
	}
	return TokenIDFromUint256(u), nil
}

// A Collection is a set of Metadata, each associated with a single uint256 (32
// bytes).
type Collection map[TokenID]*Metadata

// Metadata carries a parsed JSON payload from ERC721 metadata, compatible with
// OpenSea.
type Metadata struct {
	Name         string       `json:"name"`
	Description  string       `json:"description,omitempty"`
	Image        string       `json:"image"`
	AnimationURL string       `json:"animation_url,omitempty"`
	ExternalURL  string       `json:"external_url,omitempty"`
	Attributes   []*Attribute `json:"attributes,omitempty"`
}

// An Attribute is a single attribute in Metadata.
type Attribute struct {
	TraitType   string      `json:"trait_type,omitempty"`
	Value       interface{} `json:"value"`
	DisplayType string      `json:"display_type,omitempty"`
}

// String returns a human-readable string of TraitType:Value.
func (a *Attribute) String() string {
	return fmt.Sprintf("%s:%v", a.TraitType, a.Value)
}
