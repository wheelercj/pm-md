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
	"reflect"
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

func TestFileExists(t *testing.T) {
	if !FileExists("../LICENSE") {
		t.Error("FileExists(\"../LICENSE\") = false, want true")
	}
}

func TestFileDoesNotExist(t *testing.T) {
	if FileExists("nonexistent-file") {
		t.Error("FileExists(\"nonexistent-file\") = true, want false")
	}
}

func TestCreateUniqueFileName(t *testing.T) {
	tests := []struct {
		name, ext, want string
	}{
		{"../LICENSE", "", "../LICENSE(1)"},
		{"../README", ".md", "../README(1).md"},
		{"nonexistent-file", ".txt", "nonexistent-file.txt"},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%q,%q", test.name, test.ext)
		t.Run(testName, func(t *testing.T) {
			ans := CreateUniqueFileName(test.name, test.ext)
			if ans != test.want {
				t.Errorf(
					"CreateUniqueFileName(%q, %q) = %q, want %q",
					test.name, test.ext, ans, test.want,
				)
			}
		})
	}
}

func TestCreateUniqueFileNamePanic(t *testing.T) {
	tests := []struct {
		name, ext string
	}{
		{"../README", "md"},
		{"nonexistent-file", "."},
		{"nonexistent-file", "a"},
		{"", ""},
		{"", ".md"},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%q,%q", test.name, test.ext)
		t.Run(testName, func(t *testing.T) {
			assertPanic(t, CreateUniqueFileName, test.name, test.ext)
		})
	}
}

func TestFormatFileName(t *testing.T) {
	tests := []struct {
		name, input, want string
	}{
		{"spaces", "file name with spaces", "file-name-with-spaces"},
		{"special characters", "lots-of-#<>$+%&/\\*|{}!?`'\"=:@-special-characters", "lots-of-----------------------special-characters"},
		{"invalid start and end", ".  invalid-start-and-end__--", "invalid-start-and-end"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ans := FormatFileName(test.input)
			if ans != test.want {
				t.Errorf("FormatFileName(%q) = %q, want %q", test.input, ans, test.want)
			}
		})
	}
}

func TestFormatHeaderLink(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"", "[](#)"},
		{"create account", "[create account](#create-account)"},
		{"sample request body", "[sample request body](#sample-request-body)"},
		{"sample request body", "[sample request body](#sample-request-body-1)"},
		{"sample request body", "[sample request body](#sample-request-body-2)"},
	}

	for i, test := range tests {
		name := fmt.Sprintf("(%d) %q", i, test.input)
		t.Run(name, func(t *testing.T) {
			ans := formatHeaderLink(test.input)
			if ans != test.want {
				t.Errorf("formatHeaderLink(%q) = %q, want %q", test.input, ans, test.want)
			}
		})
	}
}

func TestFormatHeaderPath(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"with space", "#with-space"},
		{" with spaces ", "#with-spaces-"},
		{"special@#$%^&+*()=/\\|'\":;!?.>,<[]{}`~characters", "#specialcharacters"},
		{"dash-and_underscore", "#dash-and_underscore"},
		{"UPPERCASE", "#uppercase"},
		{"123", "#123"},
		{"课客果国", "#课客果国"},
		{"", "#"},
	}

	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			ans := formatHeaderPath(test.input)
			if ans != test.want {
				t.Errorf("formatHeaderPath(%q) = %q, want %q", test.input, ans, test.want)
			}
		})
	}
}
