// Copyright 2023 Chris Wheeler

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
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

// ConfirmReplaceExistingFile asks the user to confirm whether they want one of their
// existing files to be replaced. This function does NOT check whether a file exists.
func ConfirmReplaceExistingFile(fileName string) error {
	fmt.Printf("File %q already exists. Replace? (y/n) ", fileName)
	var choice string
	_, err := fmt.Scan(&choice)
	if err != nil {
		return err
	}
	choice = strings.ToLower(choice)
	if choice != "y" && choice != "n" {
		return fmt.Errorf("Invalid choice. Please choose y or n")
	} else if choice == "n" {
		return fmt.Errorf("Canceled")
	}

	return nil
}

func ScanStdin() ([]byte, error) {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("stdin scan error: %s", err)
	}
	return []byte(strings.Join(lines, "\n")), nil
}
