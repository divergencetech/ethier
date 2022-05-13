package main

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSortedNonEmpty(t *testing.T) {
	asBytes := func(strs ...string) [][]byte {
		var b [][]byte
		for _, s := range strs {
			b = append(b, []byte(s))
		}
		return b
	}

	tests := []struct {
		in   string
		want [][]byte
	}{
		{
			in: `foo
			bar
			baz`,
			want: asBytes("bar", "baz", "foo"),
		},
		{
			in: `
			
			foo
			bar
			baz
			
			`,
			want: asBytes("bar", "baz", "foo"),
		},
		{
			in: `
			
baz
		
			bar
			qux

			foo		`,
			want: asBytes("bar", "baz", "foo", "qux"),
		},
	}

	for _, tt := range tests {
		got, err := sortedNonEmpty(strings.NewReader(tt.in))
		if err != nil {
			t.Errorf("sortedNonEmpty(%q) error %v", tt.in, err)
			continue
		}

		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("sortedNonEmpty(%q) diff (-want +got):\n%s", tt.in, diff)
		}
	}
}

func TestFold(t *testing.T) {
	tests := []struct {
		name  string
		start entropy
		fold  [][]byte
		want  entropy
	}{
		{
			name:  "big-endian padding 7 bytes",
			start: 0,
			fold:  [][]byte{{1}},
			want:  1 << (7 * 8),
		},
		{
			name:  "big-endian padding 5 bytes",
			start: 0,
			fold:  [][]byte{{1, 0, 1}},
			want:  1<<(7*8) + 1<<(5*8),
		},
		{
			name:  "no padding",
			start: 0,
			fold:  [][]byte{{1, 0, 1, 0, 0, 0, 0, 0}},
			want:  1<<(7*8) + 1<<(5*8),
		},
		{
			name:  "xor words",
			start: 0,
			fold: [][]byte{
				{1, 0, 3, 0, 5, 0, 1, 1},
				{1, 1, 1, 0, 7, 1, 0, 1},
			},
			want: 0 + // unnecessary but makes others easier to read
				0<<(7*8) +
				1<<(6*8) +
				2<<(5*8) +
				0<<(4*8) +
				2<<(3*8) +
				1<<(2*8) +
				1<<(1*8) +
				0<<(0*8), // no-op shift for completeness
		},
		{
			name:  "xor with prior",
			start: 42,
			fold:  [][]byte{{1, 0, 0, 0, 0, 0, 0, 43}},
			want:  1<<(7*8) + 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.start
			for _, f := range tt.fold {
				got.fold(f)
			}
			if got != tt.want {
				t.Errorf("%T(%d) after folding %v; got %d; want %d", got, tt.start, tt.fold, got, tt.want)
			}
		})
	}
}

func TestEntropyAsSeed(t *testing.T) {
	truth := rand.New(rand.NewSource(42))
	e := entropy(42)

	r := e.rand()
	for i := 0; i < 10; i++ {
		got := r.Int63()
		want := truth.Int63()
		if got != want {
			t.Errorf("Call %d to %T.Int63() got %d; want %d", i, r, got, want)
		}
	}
}
