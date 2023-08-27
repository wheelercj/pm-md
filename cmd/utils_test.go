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
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// assertPanic takes any function and arguments for that function, calls the given
// function with the given arguments, and asserts that the given function then panics.
func assertPanic(t *testing.T, f any, args ...any) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic expected")
		}
	}()

	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}

	reflect.ValueOf(f).Call(reflectArgs)
}

// assertNoDiff compares two strings, asserting they have the same number of lines and
// the same content on each line. The strings have lines separated by linesep.
func assertNoDiff(t *testing.T, ans, want, linesep string) {
	if ans == want {
		return
	}
	ansSlice := strings.Split(ans, linesep)
	wantSlice := strings.Split(want, linesep)
	for i := 0; i < len(ansSlice); i++ {
		if i >= len(wantSlice) {
			t.Errorf(
				"Actual output longer than expected (want %d lines, got %d).\nContinues with\n  %q",
				len(wantSlice), len(ansSlice), ansSlice[i],
			)
			return
		}
		if ansSlice[i] != wantSlice[i] {
			t.Errorf(
				"Difference on line %d\nwant:\n  %q\ngot:\n  %q",
				i+1, wantSlice[i], ansSlice[i],
			)
			return
		}
	}
	if len(ansSlice) < len(wantSlice) {
		t.Errorf(
			"Actual output shorter than expected (want %d lines, got %d).\nShould continue with\n  %q",
			len(wantSlice), len(ansSlice), wantSlice[len(ansSlice)],
		)
		return
	}
	t.Errorf("The actual and expected strings don't match for an unknown reason")
}

// assertJsonToMdFileNoDiff converts JSON to markdown and asserts the resulting markdown
// is the same as a given example.
func assertJsonToMdFileNoDiff(t *testing.T, inputJsonFilePath, mdFilePath string, showResponseNames, generateFileName bool) {
	// Skip the test if unique file name creation isn't working correctly.
	TestCreateUniqueFileName(t)
	TestCreateUniqueFileNamePanic(t)
	if t.Failed() {
		return
	}
	if mdFilePath == "-" {
		t.Error("This test cannot use stdout")
		return
	}

	jsonBytes, err := os.ReadFile(inputJsonFilePath)
	if err != nil {
		t.Errorf("Failed to open %q", inputJsonFilePath)
		return
	}
	if generateFileName {
		mdFilePath, err = jsonToMdFile(jsonBytes, "", nil, showResponseNames)
	} else {
		mdFilePath, err = jsonToMdFile(jsonBytes, mdFilePath, nil, showResponseNames)
	}
	if err != nil {
		t.Errorf("jsonToMdFile: %s", err)
		return
	}
	defer os.Remove(mdFilePath)
	ansBytes, err := os.ReadFile(mdFilePath)
	if err != nil {
		t.Errorf("Failed to open %q", mdFilePath)
		return
	}
	wantBytes, err := os.ReadFile(mdFilePath)
	if err != nil {
		t.Errorf("Failed to open %q", mdFilePath)
		return
	}
	ans := strings.ReplaceAll(string(ansBytes), "\r\n", "\n")
	want := strings.ReplaceAll(string(wantBytes), "\r\n", "\n")

	assertNoDiff(t, ans, want, "\n")
}

func TestFileExists(t *testing.T) {
	if !FileExists("../LICENSE") {
		t.Error("FileExists(\"../LICENSE\") = false, want true")
	}
}

func TestFileDoesNotExist(t *testing.T) {
	if FileExists("nonexistent file") {
		t.Error("FileExists(\"nonexistent file\") = true, want false")
	}
}

func TestCreateUniqueFileName(t *testing.T) {
	tests := []struct {
		a, b, want string
	}{
		{"../LICENSE", "", "../LICENSE(1)"},
		{"../README", ".md", "../README(1).md"},
		{"nonexistent file", ".txt", "nonexistent file.txt"},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%q,%q", test.a, test.b)
		t.Run(testName, func(t *testing.T) {
			ans := CreateUniqueFileName(test.a, test.b)
			if ans != test.want {
				t.Errorf(
					"CreateUniqueFileName(%q, %q) = %q, want %q",
					test.a, test.b, ans, test.want,
				)
			}
		})
	}
}

func TestCreateUniqueFileNamePanic(t *testing.T) {
	tests := []struct {
		a, b string
	}{
		{"../README", "md"},
		{"nonexistent file", "."},
		{"nonexistent file", "a"},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%q,%q", test.a, test.b)
		t.Run(testName, func(t *testing.T) {
			assertPanic(t, CreateUniqueFileName, test.a, test.b)
		})
	}
}
