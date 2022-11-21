package main

import (
	"bytes"
	"os"
	"testing"
)

func TestConvertFile(t *testing.T) {
	tests := []struct {
		fileContent string
		bufferSize  int
		want        []byte
	}{
		{"0102abcd", 4, []byte{0x01, 0x02, 0xab, 0xcd}},
		{"0x0102abcd", 32, []byte{0x01, 0x02, 0xab, 0xcd}},
		{"abcdef", 32, []byte{0xab, 0xcd, 0xef}},
		{"0xabcdef", 2, []byte{0xab, 0xcd, 0xef}},
	}
	for _, tt := range tests {
		t.Run(tt.fileContent, func(t *testing.T) {
			// Write test file
			file, err := os.CreateTemp("", "")
			if err != nil {
				t.Fatalf(`os.CreateTemp("", "") error %v`, err)
			}
			defer os.Remove(file.Name())
			defer file.Close()

			if _, err = file.WriteString(tt.fileContent); err != nil {
				t.Fatalf("%T.WriteString(%q) error %v", file, tt.fileContent, err)
			}
			if err = convertHexFile(file.Name(), tt.bufferSize); err != nil {
				t.Fatalf("convertHexFile([tmp file], %d) error %v", tt.bufferSize, err)
			}

			// Test new content
			got, err := os.ReadFile(file.Name())
			if err != nil {
				t.Fatalf("os.ReadFile([tmp file]) error %v", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("convertHexFile() wrote %+x; want %+x", got, tt.want)
			}
		})
	}
}
