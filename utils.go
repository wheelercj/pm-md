package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// FileExists checks if a given file or folder exists on the device.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// CreateUniqueFileName returns the given file name and extension (concatenated) if no
// file with them exists. Otherwise, parentheses around a number are inserted before the
// extension to make it unique. The extension must be empty or be a period followed by
// one or more characters.
func CreateUniqueFileName(fileName, extension string) string {
	if extension == "." || (len(extension) > 0 && !strings.HasPrefix(extension, ".")) {
		panic("Extension must be empty or be a period followed by one or more characters")
	}
	uniqueFileName := fileName + extension
	for i := 1; FileExists(uniqueFileName); i++ {
		uniqueFileName = fileName + "(" + fmt.Sprint(i) + ")" + extension
	}
	return uniqueFileName
}
