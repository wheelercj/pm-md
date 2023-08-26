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

var ShowResponseNames bool
var Statuses string

var rootCmd = &cobra.Command{
	Use:     "pm-md postman_export.json",
	Short:   short,
	Long:    fmt.Sprintf("%s\n\n%s", short, jsonHelp),
	Example: "pm-md collection.json\npm-md collection.json --statuses=200\npm-md collection.json --statuses=200-299,400-499",
	Version: "v0.0.6 (you can check for updates here: https://github.com/wheelercj/pm-md/releases)",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		if !strings.HasSuffix(strings.ToLower(args[0]), ".json") {
			return fmt.Errorf("%q does not end with .json", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		jsonFilePath := args[0]
		// fmt.Println("json file path:", jsonFilePath)
		// fmt.Printf("statuses: %q\n", Statuses)
		// fmt.Printf("show response names: %q\n", ShowResponseNames)

		statusRanges, err := parseStatusRanges(Statuses)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		jsonBytes, err := os.ReadFile(jsonFilePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if mdFileName, err := jsonToMdFile(jsonBytes, statusRanges, ShowResponseNames); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			fmt.Fprintln(os.Stderr, "Created", mdFileName)
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
	rootCmd.Flags().BoolVarP(
		&ShowResponseNames,
		"show-response-names",
		"n",
		false,
		"Include the names of sample responses in the output",
	)
}
