// Package erc721 provides functionality associated with ERC721 NFTs.
package erc721

import (
	"fmt"
)

// A Collection is a set of Metadata, each associated with a single token ID.
type Collection map[TokenID]*Metadata

// CollectionFromMetadata assumes that each token's ID is its index in the
// received Metadata slice, returning the respective Collection map.
func CollectionFromMetadata(md []*Metadata) Collection {
	c := make(Collection)
	for i, m := range md {
		c[*TokenIDFromInt(i)] = m
	}
	return c
}

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
