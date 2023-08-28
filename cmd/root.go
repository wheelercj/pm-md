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
	"strings"

	"github.com/spf13/cobra"
)

const short = "Convert a Postman collection to markdown documentation"
const jsonHelp = "You can get a JSON file from Postman by exporting a collection as a v2.1.0 collection"
const github = "More help available here: github.com/wheelercj/pm-md"
const version = "v0.0.6 (you can check for updates here: https://github.com/wheelercj/pm-md/releases)"
const example = `pm-md collection.json
pm-md collection.json documentation.md
pm-md collection.json -
pm-md collection.json --statuses=200-299,400-499`

var Statuses string
var CustomTmplPath string
var ShowResponseNames bool
var GetTemplate bool
var ConfirmReplaceExistingFile bool

var rootCmd = &cobra.Command{
	Use:     "pm-md [postman_export.json [output.md]]",
	Short:   short,
	Long:    fmt.Sprintf("%s\n\n%s.\n%s", short, jsonHelp, github),
	Example: example,
	Version: version,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && GetTemplate {
			return nil
		}
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		if err := cobra.MaximumNArgs(2)(cmd, args); err != nil {
			return err
		}
		if args[0] != "-" && !strings.HasSuffix(strings.ToLower(args[0]), ".json") {
			return fmt.Errorf("%q must be \"-\" or end with \".json\"", args[0])
		}
		if len(CustomTmplPath) > 0 && !strings.HasSuffix(CustomTmplPath, ".tmpl") {
			return fmt.Errorf("%q must end with \".tmpl\"", CustomTmplPath)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if GetTemplate {
			exportDefaultTemplate()
			if len(args) == 0 {
				os.Exit(0)
			}
		}
		jsonFilePath := args[0]
		var destName string
		if len(args) == 2 {
			destName = args[1]
		}
		// fmt.Printf("json file path: %q\n", jsonFilePath)
		// fmt.Printf("output destination: %q\n", destName)
		// fmt.Printf("statuses: %q\n", Statuses)
		// fmt.Println("show response names:", ShowResponseNames)
		// fmt.Println("get template:", GetTemplate)
		// fmt.Printf("custom template: %q\n", CustomTmplPath)

		statusRanges, err := parseStatusRanges(Statuses)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var jsonBytes []byte
		if jsonFilePath == "-" {
			jsonBytes, err = ScanStdin()
		} else {
			jsonBytes, err = os.ReadFile(jsonFilePath)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		destName, err = jsonToMdFile(
			jsonBytes,
			destName,
			CustomTmplPath,
			statusRanges,
			ShowResponseNames,
			ConfirmReplaceExistingFile,
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else if destName != "-" {
			fmt.Fprintf(os.Stderr, "Created %q\n", destName)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(
		&Statuses,
		"statuses",
		"s",
		"",
		"Include only the sample responses with status codes in given range(s)",
	)
	rootCmd.Flags().StringVarP(
		&CustomTmplPath,
		"template",
		"t",
		"",
		"Use a custom template for the output",
	)
	rootCmd.Flags().BoolVarP(
		&ShowResponseNames,
		"show-response-names",
		"n",
		false,
		"Include the names of sample responses in the output",
	)
	rootCmd.Flags().BoolVarP(
		&GetTemplate,
		"get-template",
		"g",
		false,
		"Creates a file of the default template for customization",
	)
	rootCmd.Flags().BoolVar(
		&ConfirmReplaceExistingFile,
		"replace",
		false,
		"Confirm whether to replace a chosen existing output file",
	)
	rootCmd.Flags().MarkHidden("replace")
}
