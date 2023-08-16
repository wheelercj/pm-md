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
// appended to the file name to make it unique and the extension is concatenated.
// Otherwise, the given file name and extension are concatenated and returned unchanged.
// The extension must be empty or be a period followed by one or more characters.
func CreateUniqueFileName(fileName, extension string) (string, error) {
	if extension == "." || (len(extension) > 0 && !strings.HasPrefix(extension, ".")) {
		return "", fmt.Errorf("Extension must be empty or be a period followed by one or more characters")
	}
	uniqueFileName := fileName + extension
	for i := 1; FileExists(uniqueFileName); i++ {
		uniqueFileName = fileName + "(" + fmt.Sprint(i) + ")" + extension
	}
	return uniqueFileName, nil
}
