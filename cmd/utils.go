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
	"io"
	"os"
	"strings"
)

// FileExists checks if a given file or folder exists on the device.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// CreateUniqueFileName returns the given file name and extension (concatenated) if no
// file with them exists. Otherwise, a period and a number are inserted before the
// extension to make it unique. The extension must be empty or be a period followed by
// one or more characters. The function panics if the given file name is empty, if the
// extension is only ".", or if the extension is not empty but does not start with a
// period.
func CreateUniqueFileName(fileName, extension string) string {
	if len(fileName) == 0 {
		panic("The file name must not be empty")
	}
	if extension == "." || (len(extension) > 0 && !strings.HasPrefix(extension, ".")) {
		panic("Extension must be empty or be a period followed by one or more characters")
	}
	uniqueFileName := fileName + extension
	for i := 1; FileExists(uniqueFileName); i++ {
		uniqueFileName = fileName + "." + fmt.Sprint(i) + extension
	}
	return uniqueFileName
}

// FormatFileName takes a file name excluding any file extension and changes it, if
// necessary, to be compatible with all major platforms. Each invalid file name
// character is replaced with a dash, and characters that a file name cannot start or
// end with are trimmed. The invalid characters are `#<>$+%&/\\*|{}!?`'\"=: @`,
// and the invalid start or end characters are ` ._-`.
func FormatFileName(fileName string) string {
	invalidChars := "#<>$+%&/\\*|{}!?`'\"=: @"
	invalidEdgeChars := " ._-"

	result := make([]byte, len(fileName))
	for i := range fileName {
		if strings.Contains(invalidChars, string(fileName[i])) {
			result[i] = '-'
		} else {
			result[i] = fileName[i]
		}
	}

	return strings.Trim(string(result), invalidEdgeChars)
}

// ScanStdin reads input from stdin until it finds EOF or a different error, and then
// returns any input all at once. If EOF is found, the returned error is nil.
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

// exportText creates a new file with a unique name based on the given base name (no
// existing file will ever be replaced), saves the given content into it, and returns
// the new file's name. The given file extension must be empty or be a period followed
// by one or more characters.
func exportText(baseName, ext, content string) string {
	uniqueName := CreateUniqueFileName(baseName, ext)
	file, err := os.Create(uniqueName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()
	_, err = file.Write([]byte(content))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return uniqueName
}

// AssertGenerateNoDiff converts JSON to plaintext and asserts the result is the same as
// wanted text. wantPath is the path to an existing file containing the wanted output.
// If the given template path is empty, the default template is used. If any status
// ranges are given, responses with statuses outside those ranges will not be present in
// the result.
func AssertGenerateNoDiff(jsonPath, tmplPath, wantPath string, statusRanges [][]int) error {
	jsonBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	openAnsFile, err := os.CreateTemp("", "pm2md_*.md")
	if err != nil {
		return err
	}
	defer os.Remove(openAnsFile.Name())
	defer openAnsFile.Close()
	wantBytes, err := os.ReadFile(wantPath)
	if err != nil {
		return err
	}

	collection, err := parseCollection(jsonBytes)
	if err != nil {
		return err
	}

	err = generateText(
		collection,
		openAnsFile,
		tmplPath,
		statusRanges,
	)
	if err != nil {
		return err
	}
	fileInfo, err := openAnsFile.Stat()
	if err != nil {
		return err
	}
	ansBytes := make([]byte, fileInfo.Size())
	_, err = openAnsFile.Read(ansBytes)
	if err != nil && err != io.EOF {
		return err
	}

	ans := strings.ReplaceAll(string(ansBytes), "\r\n", "\n")
	want := strings.ReplaceAll(string(wantBytes), "\r\n", "\n")

	return AssertNoDiff(ans, want, "\n")
}

// AssertNoDiff compares two strings, asserting they have the same number of lines and
// the same content on each line. The strings have lines separated by linesep.
func AssertNoDiff(ans, want, linesep string) error {
	if ans == want {
		return nil
	}
	ansSlice := strings.Split(ans, linesep)
	wantSlice := strings.Split(want, linesep)
	for i := 0; i < len(ansSlice); i++ {
		if i >= len(wantSlice) {
			return fmt.Errorf(
				"Actual output longer than expected (want %d lines, got %d).\nContinues with\n  %q",
				len(wantSlice), len(ansSlice), ansSlice[i],
			)
		}
		if ansSlice[i] != wantSlice[i] {
			return fmt.Errorf(
				"Difference on line %d\nwant:\n  %q\ngot:\n  %q",
				i+1, wantSlice[i], ansSlice[i],
			)
		}
	}
	if len(ansSlice) < len(wantSlice) {
		return fmt.Errorf(
			"Actual output shorter than expected (want %d lines, got %d).\nShould continue with\n  %q",
			len(wantSlice), len(ansSlice), wantSlice[len(ansSlice)],
		)
	}

	return fmt.Errorf("The actual and expected strings don't match for an unknown reason")
}
