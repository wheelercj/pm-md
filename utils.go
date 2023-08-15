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

// If an existing file has the given file name and extension, parentheses around a
// number are appended to the file name to make it unique. Otherwise, the given file
// name and extension are concatenated and returned unchanged. The extension must start
// with a period.
func CreateUniqueFileName(fileName, extension string) string {
	if !strings.HasPrefix(extension, ".") {
		panic("extension must start with a period")
	}
	uniqueFileName := fileName + extension
	for i := 1; FileExists(uniqueFileName); i++ {
		uniqueFileName = fileName + "(" + fmt.Sprint(i) + ")" + extension
	}
	return uniqueFileName
}
