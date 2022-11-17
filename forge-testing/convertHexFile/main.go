// The convertHexFile binary convert a file containing a hex string into
// a file containing the content of the hex string.

package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/golang/glog"
)

func main() {
	path := os.Args[1]
	if err := convertHexFile(path); err != nil {
		glog.Exit(err)
	}
}

func convertHexFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("os.Open(%q): %w", path, err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("io.ReadAll(%T): %w", f, err)
	}
	f.Close()

	data, err := hex.DecodeString(strings.TrimPrefix(string(b), "0x"))
	if err != nil {
		return fmt.Errorf("hex.DecodeString([data]): %w", err)
	}

	f, err = os.Create(path)
	if err != nil {
		return fmt.Errorf("os.Create(%q): %w", path, err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("%T.Write(%T): %w", f, data, err)
	}

	return nil
}
