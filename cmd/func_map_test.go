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
	"testing"
)

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
