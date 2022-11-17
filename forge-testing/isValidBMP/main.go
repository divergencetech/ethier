// The isValidBMP binary checks if a given file contains a valid BMP image.

package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/golang/glog"
	"golang.org/x/image/bmp"
)

func main() {
	path := os.Args[1]
	valid, err := isValidBMP(path)
	if err != nil {
		glog.Exit(err)
	}

	if valid {
		fmt.Println("0")
	} else {
		fmt.Println("1")
	}
}

func isValidBMP(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("os.Open(%q): %w", path, err)
	}
	defer file.Close()

	_, err = bmp.Decode(file)
	if err != nil {
		return false, fmt.Errorf("bmp.Decode(%T): %w", file, err)
	}

	return true, nil
}
