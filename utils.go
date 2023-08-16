package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Checks if a given file or folder exists on the device.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// If an existing file has the given name and extension, parentheses around a number are
// appended to the file name to make it unique. Otherwise, the given file name remains
// unchanged. The extension is then concatenated. The extension must be empty or be a
// period followed by one or more characters.
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
