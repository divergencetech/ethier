package erc721

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOpenSeaDisplayType(t *testing.T) {
	var _ interface {
		json.Marshaler
		json.Unmarshaler
	} = new(OpenSeaDisplayType)

	tests := []struct {
		attr   *Attribute
		asJSON string
	}{
		{
			attr:   &Attribute{},
			asJSON: `{"value":null}`,
		},
		{
			attr:   &Attribute{DisplayType: DisplayDefault},
			asJSON: `{"value":null}`,
		},
		{
			attr:   &Attribute{DisplayType: DisplayNumber},
			asJSON: `{"value":null,"display_type":"number"}`,
		},
		{
			attr:   &Attribute{DisplayType: DisplayBoostNumber},
			asJSON: `{"value":null,"display_type":"boost_number"}`,
		},
		{
			attr:   &Attribute{DisplayType: DisplayBoostPercentage},
			asJSON: `{"value":null,"display_type":"boost_percentage"}`,
		},
		{
			attr:   &Attribute{DisplayType: DisplayDate},
			asJSON: `{"value":null,"display_type":"date"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.attr.DisplayType.String(), func(t *testing.T) {
			t.Run("to JSON", func(t *testing.T) {
				got, err := json.Marshal(tt.attr)
				if err != nil {
					t.Fatalf("json.Marshal(%+v) error %v", tt.attr, err)
				}
				if diff := cmp.Diff(tt.asJSON, string(got)); diff != "" {
					t.Errorf("json.Marshal(%T %+v) diff (-want +got):\n%s", tt.attr, tt.attr, diff)
				}
			})

			t.Run("from JSON", func(t *testing.T) {
				got := new(Attribute)
				if err := json.Unmarshal([]byte(tt.asJSON), got); err != nil {
					t.Fatalf("json.Unmarshal(%q) error %v", tt.asJSON, err)
				}
				if diff := cmp.Diff(tt.attr, got); diff != "" {
					t.Errorf("json.Unmarshal(%q, %T) diff (-want +got):\n%s", tt.asJSON, got, diff)
				}
			})
		})
	}
}
