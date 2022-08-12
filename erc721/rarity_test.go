package erc721

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// md returns a new Metadata with each pair of strings used as the "key"
// (TraitType) and Value of a new Attribute.
func md(t *testing.T, attrKVPairs ...string) *Metadata {
	t.Helper()
	if len(attrKVPairs)%2 != 0 {
		t.Fatalf("bad test setup; odd number of attribute KV pairs %q", attrKVPairs)
	}

	var m Metadata
	for i := 0; i < len(attrKVPairs); i += 2 {
		a := &Attribute{
			TraitType: attrKVPairs[i],
			Value:     attrKVPairs[i+1],
		}
		m.Attributes = append(m.Attributes, a)
	}
	return &m
}

func TestRarity(t *testing.T) {
	id := func(i int) TokenID {
		return *TokenIDFromInt(i)
	}

	tests := []struct {
		name     string
		metadata []*Metadata
		bucket   func(interface{}) string
		want     *Rarity
	}{
		{
			name: "single item has zero information",
			metadata: []*Metadata{
				md(t, "key", "val", "foo", "bar", "hello", "world"),
			},
			want: &Rarity{
				Entropy: 0,
				Scores: map[TokenID]float64{
					id(0): 0,
				},
			},
		},
		{
			name: "identical items have zero information because 100% probability",
			metadata: []*Metadata{
				md(t, "key", "val", "foo", "bar", "hello", "world"),
				md(t, "key", "val", "foo", "bar", "hello", "world"),
				md(t, "key", "val", "foo", "bar", "hello", "world"),
			},
			want: &Rarity{
				Entropy: 0,
				Scores: map[TokenID]float64{
					id(0): 0,
					id(1): 0,
					id(2): 0,
				},
			},
		},
		{
			name: "uniform probability for one trait across 2 items == 1 bit",
			metadata: []*Metadata{
				md(t, "key", "0"),
				md(t, "key", "1"),
			},
			want: &Rarity{
				Entropy: 1,
				Scores: map[TokenID]float64{
					id(0): 1,
					id(1): 1,
				},
			},
		},
		{
			name: "uniform probability for two traits across 2 items == 2 bits",
			metadata: []*Metadata{
				md(t, "key", "0", "foo", "bar"),
				md(t, "key", "1", "foo", "baz"),
			},
			want: &Rarity{
				Entropy: 2,
				Scores: map[TokenID]float64{
					id(0): 1,
					id(1): 1,
				},
			},
		},
		{
			name: "uniform probability for one trait across 4 items == 2 bits",
			metadata: []*Metadata{
				md(t, "key", "0"),
				md(t, "key", "1"),
				md(t, "key", "2"),
				md(t, "key", "3"),
			},
			want: &Rarity{
				Entropy: 2,
				Scores: map[TokenID]float64{
					id(0): 1,
					id(1): 1,
					id(2): 1,
					id(3): 1,
				},
			},
		},
		{
			name: "two identical items and two different items == 1 + 0.5 bits",
			metadata: []*Metadata{
				md(t, "key", "0"),
				md(t, "key", "0"),
				md(t, "key", "2"),
				md(t, "key", "3"),
			},
			want: &Rarity{
				Entropy: 1.5,
				Scores: map[TokenID]float64{
					id(0): 1. / 1.5,
					id(1): 1. / 1.5,
					id(2): 2. / 1.5,
					id(3): 2. / 1.5,
				},
			},
		},
		{
			name: "absence of a trait treated as a null value",
			metadata: []*Metadata{
				md(t),
				md(t),
				md(t, "key", "2"),
				md(t, "key", "3"),
			},
			want: &Rarity{
				Entropy: 1.5,
				Scores: map[TokenID]float64{
					id(0): 1. / 1.5,
					id(1): 1. / 1.5,
					id(2): 2. / 1.5,
					id(3): 2. / 1.5,
				},
			},
		},
		{
			name: "bucket fuction used for non-string attribute",
			bucket: func(i interface{}) string {
				return fmt.Sprintf("%T:%v", i, i)
			},
			metadata: []*Metadata{
				{Attributes: []*Attribute{{TraitType: "key", Value: uint16(42)}}},
				md(t, "key", "uint16:42"), // effectively identical to above
				md(t, "key", "val"),
				md(t, "key", "val"),
			},
			want: &Rarity{
				Entropy: 1,
				Scores: map[TokenID]float64{
					id(0): 1,
					id(1): 1,
					id(2): 1,
					id(3): 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("%T generated from %+v", Collection{}, tt.metadata)
			got := CollectionFromMetadata(tt.metadata).Rarity(tt.bucket)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%T.Rarity() diff (-want +got):\n%s", Collection{}, diff)
			}
		})
	}
}

func TestMD(t *testing.T) {
	tests := []struct {
		attributes []string
		want       *Metadata
	}{
		{
			attributes: nil,
			want:       &Metadata{},
		},
		{
			attributes: []string{"key", "val"},
			want: &Metadata{
				Attributes: []*Attribute{
					{TraitType: "key", Value: "val"},
				},
			},
		},
		{
			attributes: []string{"key", "val", "foo", "bar"},
			want: &Metadata{
				Attributes: []*Attribute{
					{TraitType: "key", Value: "val"},
					{TraitType: "foo", Value: "bar"},
				},
			},
		},
	}

	for _, tt := range tests {
		if diff := cmp.Diff(tt.want, md(t, tt.attributes...)); diff != "" {
			t.Errorf("Incorrect test helper md(%q) diff (-want +got):\n%s", tt.attributes, diff)
		}
	}
}
