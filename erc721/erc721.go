// Package erc721 provides functionality associated with ERC721 NFTs.
package erc721

import (
	"encoding/json"
	"fmt"
	"io"
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

// MarshalJSONTo marshals the Metadata to JSON and writes it to the Writer,
// returning the number of bytes written and any error that may occur.
func (md *Metadata) MarshalJSONTo(w io.Writer) (int, error) {
	buf, err := json.Marshal(md)
	if err != nil {
		return 0, fmt.Errorf("json.Marshal(%T): %v", md, err)
	}
	return w.Write(buf)
}

// An Attribute is a single attribute in Metadata.
type Attribute struct {
	TraitType   string             `json:"trait_type,omitempty"`
	Value       interface{}        `json:"value"`
	DisplayType OpenSeaDisplayType `json:"display_type,omitempty"`
}

// An OpenSeaDisplayType is an OpenSea-specific metadata concept to control how
// their UI treats numerical values.
//
// See https://docs.opensea.io/docs/metadata-standards for details.
type OpenSeaDisplayType int

// Allowable OpenSeaDisplayType values.
const (
	DisplayDefault OpenSeaDisplayType = iota
	DisplayNumber
	DisplayBoostNumber
	DisplayBoostPercentage
	DisplayDate
	endDisplayTypes // used for bounds checking
)

// String returns the display type as a string.
func (t OpenSeaDisplayType) String() string {
	switch t {
	case DisplayDefault:
		return ""
	case DisplayNumber:
		return "number"
	case DisplayBoostNumber:
		return "boost_number"
	case DisplayBoostPercentage:
		return "boost_percentage"
	case DisplayDate:
		return "date"
	default:
		return fmt.Sprintf("%T(%d)", t, t)
	}
}

// MarshalJSON returns the display type as JSON.
func (t OpenSeaDisplayType) MarshalJSON() ([]byte, error) {
	if t == DisplayDefault {
		return nil, nil
	}
	if t >= endDisplayTypes {
		return nil, fmt.Errorf("unsupported %T = %d", t, t)
	}
	return []byte(fmt.Sprintf("%q", t.String())), nil
}

// UnmarshalJSON parses the JSON buffer into the display type.
func (t *OpenSeaDisplayType) UnmarshalJSON(buf []byte) error {
	var s string
	if err := json.Unmarshal(buf, &s); err != nil {
		return err
	}

	switch s {
	case "":
		*t = DisplayDefault
	case "number":
		*t = DisplayNumber
	case "boost_number":
		*t = DisplayBoostNumber
	case "boost_percentage":
		*t = DisplayBoostPercentage
	case "date":
		*t = DisplayDate
	default:
		return fmt.Errorf("unsupported %T = %q", *t, s)
	}
	return nil
}

// String returns a human-readable string of TraitType:Value.
func (a *Attribute) String() string {
	return fmt.Sprintf("%s:%v", a.TraitType, a.Value)
}
