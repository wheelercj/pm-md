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

package main

import (
	"fmt"
	"reflect"
	"testing"
)

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
	if !FileExists("LICENSE") {
		t.Error("FileExists(\"LICENSE\") = false, want true")
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
		{"LICENSE", "", "LICENSE(1)"},
		{"README", ".md", "README(1).md"},
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
		{"README", "md"},
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
