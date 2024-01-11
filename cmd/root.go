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
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed default.tmpl
var defaultTmplStr string

//go:embed minimal.tmpl
var minimalTmplStr string

const defaultTmplName = "default.tmpl"
// const minimalTmplName = "minimal.tmpl"

const short = "Convert a Postman collection to markdown documentation"
const jsonHelp = "You can get a JSON file from Postman by exporting a collection as a v2.1 collection"
const github = "More help available here: github.com/wheelercj/pm2md"
const version = "v0.0.11 (you can check for updates here: https://github.com/wheelercj/pm2md/releases)"
const example = `  pm2md collection.json
  pm2md collection.json output.md
  pm2md collection.json --template=custom.tmpl
  pm2md test collection.json custom.tmpl expected.md`

var Statuses string
var CustomTmplPath string
var GetDefault bool
var GetMinimal bool
var ConfirmReplaceExistingFile bool

var rootCmd = &cobra.Command{
	Use:     "pm2md [postman_export.json [output.md]]",
	Short:   short,
	Long:    fmt.Sprintf("%s\n\n%s.\n%s", short, jsonHelp, github),
	Example: example,
	Version: version,
	Args:    argsFunc,
	RunE:    runFunc,
}

// argsFunc does some input validation on the command args and flags.
func argsFunc(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && (GetDefault || GetMinimal) {
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
}

// runFunc parses command args and flags, generates plaintext, and saves the result to a
// file or prints to stdout.
func runFunc(cmd *cobra.Command, args []string) error {
	destPath, destFile, collection, statusRanges, err := parseInput(cmd, args)
	if err != nil {
		return err
	}
	if destFile != os.Stdout {
		defer destFile.Close()
	}

	err = generateText(
		collection,
		destFile,
		CustomTmplPath,
		statusRanges,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else if destPath != "-" {
		fmt.Fprintf(os.Stderr, "Created %q\n", destPath)
	}
	return nil
}

// parseInput parses command args and flags, opens the destination file, and returns all
// of these results.
func parseInput(cmd *cobra.Command, args []string) (string, *os.File, map[string]any, [][]int, error) {
	if GetDefault {
		fileName := exportText("default", ".tmpl", defaultTmplStr)
		fmt.Fprintf(os.Stderr, "Created %q\n", fileName)
		if len(args) == 0 {
			os.Exit(0)
		}
	}
	if GetMinimal {
		fileName := exportText("minimal", ".tmpl", minimalTmplStr)
		fmt.Fprintf(os.Stderr, "Created %q\n", fileName)
		if len(args) == 0 {
			os.Exit(0)
		}
	}

	jsonPath := args[0]
	var destPath string
	if len(args) == 2 {
		destPath = args[1]
	}

	statusRanges, err := parseStatusRanges(Statuses)
	if err != nil {
		return "", nil, nil, nil, err
	}

	var jsonBytes []byte
	if jsonPath == "-" {
		jsonBytes, err = ScanStdin()
	} else {
		jsonBytes, err = os.ReadFile(jsonPath)
	}
	if err != nil {
		return "", nil, nil, nil, err
	}
	collection, err := parseCollection(jsonBytes)
	if err != nil {
		return "", nil, nil, nil, err
	}

	collectionName := collection["info"].(map[string]any)["name"].(string)
	destFile, destPath, err := openDestFile(destPath, collectionName, ConfirmReplaceExistingFile)
	if err != nil {
		return "", nil, nil, nil, err
	}

	return destPath, destFile, collection, statusRanges, nil
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
	rootCmd.AddCommand(testCmd)

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
		&GetDefault,
		"get-default",
		"d",
		false,
		"Creates a file of the default template for customization",
	)
	rootCmd.Flags().BoolVarP(
		&GetMinimal,
		"get-minimal",
		"m",
		false,
		"Creates a file of a minimal template for customization",
	)
	rootCmd.Flags().BoolVar(
		&ConfirmReplaceExistingFile,
		"replace",
		false,
		"Confirm whether to replace a chosen existing output file",
	)
	rootCmd.Flags().MarkHidden("replace")
}

// openDestFile gets the destination file and its path. If the given destination path is
// "-", the destination file is os.Stdout. If the given destination path is empty, a new
// file is created with a path based on the collection name and the returned path will
// be different from the given one. If the given destination path refers to an existing
// file and confirmation to replace an existing file is not given, an error is returned.
// Any returned file is open.
func openDestFile(destPath, collectionName string, confirmReplaceExistingFile bool) (*os.File, string, error) {
	if destPath == "-" {
		return os.Stdout, destPath, nil
	}
	if len(destPath) == 0 {
		fileName := FormatFileName(collectionName)
		if len(fileName) == 0 {
			fileName = "collection"
		}
		destPath = CreateUniqueFileName(fileName, ".md")
	} else if FileExists(destPath) && !confirmReplaceExistingFile {
		return nil, "", fmt.Errorf("file %q already exists. Run the command again with the --replace flag to confirm replacing it", destPath)
	}
	destFile, err := os.Create(destPath)
	if err != nil {
		return nil, "", fmt.Errorf("os.Create: %s", err)
	}
	return destFile, destPath, nil
}
