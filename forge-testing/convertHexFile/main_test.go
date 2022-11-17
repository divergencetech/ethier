package main

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestConvertFile(t *testing.T) {
	tests := []struct {
		FileContent string
		Want        []byte
	}{
		{"0x0102abcd", []byte{0x01, 0x02, 0xab, 0xcd}},
		{"abcdef", []byte{0xab, 0xcd, 0xef}},
	}
	for _, tt := range tests {
		// Write test file

		file, err := os.CreateTemp("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(file.Name())
		defer file.Close()

		_, err = file.WriteString(tt.FileContent)
		if err != nil {
			t.Fatal(err)
		}

		err = convertHexFile(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		// Test new content

		file, err = os.Open(file.Name())
		if err != nil {
			t.Fatal(err)
		}

		got, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, tt.Want) {
			t.Errorf("Incorrect file content after conversion, want %x, got %x", tt.Want, got)
		}
	}
}
