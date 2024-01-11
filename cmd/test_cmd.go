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
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test [api.json custom.tmpl expected.md]",
	Short: "Test your custom template with expected output",
	Args:  testArgsFunc,
	RunE:  testRunFunc,
}

// testArgsFunc does some input validation on the `test` subcommand's args and flags.
func testArgsFunc(cmd *cobra.Command, args []string) error {
	if len(CustomTmplPath) > 0 {
		return fmt.Errorf("with the test subcommand, choose a custom template without using the flag")
	}
	if err := cobra.ExactArgs(3)(cmd, args); err != nil {
		return err
	}
	if !strings.HasSuffix(strings.ToLower(args[0]), ".json") {
		return fmt.Errorf("%q must end with \".json\"", args[0])
	}
	if !strings.HasSuffix(strings.ToLower(args[1]), ".tmpl") {
		return fmt.Errorf("%q must end with \".tmpl\"", args[1])
	}
	return nil
}

// testRunFunc parses the `test` subcommand's args and flags, and asserts the given JSON
// and template result in the given plaintext.
func testRunFunc(cmd *cobra.Command, args []string) error {
	jsonPath := args[0]
	tmplPath := args[1]
	wantPath := args[2]

	statusRanges, err := parseStatusRanges(Statuses)
	if err != nil {
		return err
	}

	err = AssertGenerateNoDiff(jsonPath, tmplPath, wantPath, statusRanges)
	if err == nil {
		fmt.Fprintf(os.Stderr, "Perfect match!")
	} else {
		fmt.Fprint(os.Stderr, fmt.Sprint(err))
	}

	return nil
}

// loadTmpl loads a template's name and the template itself into strings. If the given
// template path is empty, the default template is used.
func loadTmpl(tmplPath string) (tmplName string, tmplStr string, err error) {
	if len(tmplPath) > 0 {
		tmplBytes, err := os.ReadFile(tmplPath)
		if err != nil {
			return "", "", err
		}
		tmplStr = string(tmplBytes)
		tmplName = path.Base(strings.ReplaceAll(tmplPath, "\\", "/"))
	} else {
		tmplStr = defaultTmplStr
		tmplName = defaultTmplName
	}

	return tmplName, tmplStr, nil
}
