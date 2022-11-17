// The convertHexFile binary convert a file containing a hex string into
// a file containing the content of the hex string.
package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	path := os.Args[1]
	if err := convertHexFile(path); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func convertHexFile(path string) (retErr error) {
	fin, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("os.Open(%q): %w", path, err)
	}

	fout, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("os.CreateTemp(): %w", err)
	}

	r := bufio.NewReader(fin)
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%T.Read(): %w", r, err)
		}

		buf = buf[:n]

		data, err := hex.DecodeString(strings.TrimPrefix(string(buf), "0x"))
		if err != nil {
			return fmt.Errorf("hex.DecodeString([data]): %w", err)
		}
		_, err = fout.Write(data)
		if err != nil {
			return fmt.Errorf("%T.Write(%T): %w", fout, buf, err)
		}

	}

	if err := fin.Close(); err != nil {
		return fmt.Errorf("%T.Close(): %w", fin, err)
	}

	if err := fout.Close(); err != nil {
		return fmt.Errorf("%T.Close(): %w", fout, err)
	}

	if err := os.Rename(fout.Name(), fin.Name()); err != nil {
		return fmt.Errorf("os.Rename(%q,%q): %w", fout.Name(), fin.Name(), err)
	}

	return nil
}
