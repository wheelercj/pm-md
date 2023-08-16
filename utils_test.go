package main

import (
	"fmt"
	"testing"
)

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
			ans, err := CreateUniqueFileName(test.a, test.b)
			if ans != test.want || err != nil {
				t.Errorf(
					"CreateUniqueFileName(%q, %q) = (%q, %q), want (%q, nil)",
					test.a, test.b, ans, err, test.want,
				)
			}
		})
	}
}

func TestCreateUniqueFileNameError(t *testing.T) {
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
			if ans, err := CreateUniqueFileName(test.a, test.b); err == nil {
				t.Errorf(
					"CreateUniqueFileName(%q, %q) = (%q, nil), want (\"\", non-nil error)",
					test.a, test.b, ans,
				)
			}
		})
	}
}
