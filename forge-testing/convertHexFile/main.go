// The convertHexFile binary convert a file containing a hex string into
// a file containing the content of the hex string.
package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {
	path := os.Args[1]
	if err := convertHexFile(path, 4*1024); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func convertHexFile(path string, bufSize int) error {
	fIn, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("os.Open(%q): %v", path, err)
	}
	defer fIn.Close()

	// Fast-forward if we have a 0x prefix.
	prefix, err := io.ReadAll(io.NewSectionReader(fIn, 0, 2))
	if err != nil {
		return fmt.Errorf("io.ReadAll(io.NewSectionReader([%q], 0, 2)): %v", path, err)
	}
	if string(prefix) == "0x" {
		if at, err := fIn.Seek(2, io.SeekStart); err != nil || at != 2 {
			return fmt.Errorf("%T.Seek(0, Start): (%d, %v); want (2, nil)", fIn, at, err)
		}
	}

	// Although we can technically write to the beginning of the file because
	// reads are twice as fast as writes, that risks corrupting a file if
	// there's an error in the conversion.
	fOut, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf(`os.CreateTemp("", ""): %v`, err)
	}

	dec := hex.NewDecoder(fIn)
	buf := make([]byte, bufSize)
	for {
		n, err := dec.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%T.Read(…): %v", dec, err)
		}

		if _, err := fOut.Write(buf[:n]); err != nil {
			return fmt.Errorf("%T.Write(…): %v", fOut, err)
		}
	}
	if err := fOut.Close(); err != nil {
		return fmt.Errorf("%T.Close(): %v", fOut, err)
	}

	return os.Rename(fOut.Name(), fIn.Name())
}
