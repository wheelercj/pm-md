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

// assertJsonToMdFileNoDiff converts JSON to markdown and asserts the resulting markdown
// is the same as a given example. If the given custom template path is empty, the
// default template is used. If the given output path is empty, a new file is created
// with a unique name based on the JSON. The wanted output path is the path to an
// existing file containing the wanted output.
func assertJsonToMdFileNoDiff(t *testing.T, inputJsonFilePath, customTmplPath, outputPath, wantOutputPath string) {
	// Skip this test if unique file name creation isn't working correctly.
	TestCreateUniqueFileName(t)
	TestCreateUniqueFileNamePanic(t)
	if t.Failed() {
		return
	}
	if outputPath == "-" {
		t.Error("This test cannot use stdout")
		return
	}

	jsonBytes, err := os.ReadFile(inputJsonFilePath)
	if err != nil {
		t.Errorf("Failed to open %q", inputJsonFilePath)
		return
	}
	tmplName, tmplStr, err := loadTmpl(customTmplPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	outputPath, err = jsonToMdFile(
		jsonBytes,
		outputPath,
		tmplName,
		tmplStr,
		nil,
		false,
	)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(outputPath)
	ansBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Errorf("Failed to open %q", outputPath)
		return
	}
	wantBytes, err := os.ReadFile(wantOutputPath)
	if err != nil {
		t.Errorf("Failed to open %q", wantOutputPath)
		return
	}
	ans := strings.ReplaceAll(string(ansBytes), "\r\n", "\n")
	want := strings.ReplaceAll(string(wantBytes), "\r\n", "\n")

	assertNoDiff(t, ans, want, "\n")
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

func TestParseStatusRanges(t *testing.T) {
	tests := []struct {
		input string
		want  [][]int
	}{
		{"", nil},
		{"200", [][]int{{200, 200}}},
		{"200-299", [][]int{{200, 299}}},
		{"200-299,400-499", [][]int{{200, 299}, {400, 499}}},
		{"200-200", [][]int{{200, 200}}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			ans, err := parseStatusRanges(test.input)
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(ans, test.want) {
				t.Errorf("parseStatusRanges(%q) = %v, want %v", test.input, ans, test.want)
				return
			}
		})
	}
}

func TestParseStatusRangesWithInvalidInput(t *testing.T) {
	inputs := []string{"200-299-300", "a-299", "200-b", "200-", "-299", "-", "a"}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if statusRanges, err := parseStatusRanges(input); err == nil {
				t.Errorf("parseStatusRanges(%q) = (%v, nil), want non-nil error", input, statusRanges)
			}
		})
	}
}

func TestParseEmptyCollection(t *testing.T) {
	collection, err := parseCollection([]byte(""))
	if err == nil {
		t.Errorf("parseCollection([]byte(\"\")) = (%v, %v), want (nil, error)", collection, err)
	}
}

func TestJsonToMdFile(t *testing.T) {
	inputFilePath := "../samples/calendar-API.postman_collection.json"
	wantOutputPath := "../samples/calendar-API-v1.md"
	assertJsonToMdFileNoDiff(t, inputFilePath, "", "", wantOutputPath)
}

func TestJsonToMdFileWithCustomOutputFileName(t *testing.T) {
	inputFilePath := "../samples/calendar-API.postman_collection.json"
	customOutputFileName := "custom-file-name-for-testing.md"
	wantOutputPath := "../samples/calendar-API-v1.md"
	assertJsonToMdFileNoDiff(t, inputFilePath, "", customOutputFileName, wantOutputPath)
}

func TestJsonToMdFileWithCustomTemplate(t *testing.T) {
	inputFilePath := "../samples/calendar-API.postman_collection.json"
	customTmplPath := "../samples/custom.tmpl"
	wantOutputPath := "../samples/calendar-API-v1-from-custom-templ.md"
	assertJsonToMdFileNoDiff(t, inputFilePath, customTmplPath, "", wantOutputPath)
}

func TestInvalidJsonToMdFile(t *testing.T) {
	// Skip this test if unique file name creation isn't working correctly.
	TestCreateUniqueFileName(t)
	TestCreateUniqueFileNamePanic(t)
	if t.Failed() {
		return
	}

	invalidJson := []byte(`
		{
			"info": {
				"_postman_id": "23799766-64ba-4c7c-aaa9-0d880964db54",
				"name": "calendar API",
				"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
				"_exporter_id": "23363106"
			},
	`)
	tmplName, tmplStr, err := loadTmpl("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	destName, err := jsonToMdFile(invalidJson, "-", tmplName, tmplStr, nil, false)
	if err == nil {
		t.Error("Error expected")
		if destName != "-" {
			t.Errorf("The destination name should not have changed from \"-\" to %q", destName)
		}
	}
}

func TestJsonToMdFileExistingFileErr(t *testing.T) {
	inputFilePath := "../samples/calendar-API.postman_collection.json"
	jsonBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		t.Errorf("Failed to open %s", inputFilePath)
		return
	}
	tmplName, tmplStr, err := loadTmpl("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	destName, err := jsonToMdFile(jsonBytes, "../LICENSE", tmplName, tmplStr, nil, false)
	if err == nil {
		t.Errorf("jsonToMdFile(jsonBytes, \"../LICENSE\", \"\", nil, false) = (%q, nil), want non-nil error", destName)
	}
}

func TestParseCollectionWithOldSchema(t *testing.T) {
	inputFilePath := "../samples/calendar-API.postman_collection.json"
	jsonBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		t.Errorf("Failed to open %s", inputFilePath)
		return
	}
	jsonStr := string(jsonBytes)

	v210Url := "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	v200Url := "https://schema.getpostman.com/json/collection/v2.0.0/collection.json"
	if !strings.Contains(jsonStr, v210Url) {
		t.Error("The given JSON doesn't contain the expected URL")
		return
	}
	jsonStr = strings.Replace(jsonStr, v210Url, v200Url, 1)

	if collection, err := parseCollection([]byte(jsonStr)); err == nil {
		t.Errorf("want (nil, error), got a nil error and a non-nil collection: %v", collection)
	}
}

func getCollection(t *testing.T) (map[string]any, error) {
	inputFilePath := "../samples/calendar-API.postman_collection.json"
	jsonBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s", inputFilePath)
	}

	collection, err := parseCollection(jsonBytes)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func TestFilterResponses(t *testing.T) {
	collection, err := getCollection(t)
	if err != nil {
		t.Error(err)
		return
	}

	filterResponsesByStatus(collection, [][]int{{200, 200}})
	for _, endpointAny := range collection["item"].([]any) {
		endpoint := endpointAny.(map[string]any)
		for _, responseAny := range endpoint["response"].([]any) {
			response := responseAny.(map[string]any)
			code := int(response["code"].(float64))
			if code != 200 {
				t.Errorf("want 200, got %d", code)
				return
			}
		}
	}
}

func TestGetDestFileStdout(t *testing.T) {
	destFile, destName, err := getDestFile("-", "", false)
	if destFile != os.Stdout || destName != "-" || err != nil {
		t.Errorf("getDestFile(\"-\", \"\") = (%p, %q, %q), want (%p, \"-\", nil)", destFile, destName, err, os.Stdout)
	}
}

func TestGetDestFileExistingFileErr(t *testing.T) {
	destFile, destName, err := getDestFile("../LICENSE", "", false)
	if err == nil {
		t.Errorf("getDestFile(\"../LICENSE\", \"\", false) = (%p, %q, nil), want non-nil error", destFile, destName)
		if destName != "-" {
			destFile.Close()
		}
	}
}

func TestGetDestFile(t *testing.T) {
	tests := []struct {
		originalDestName, collectionName, wantName string
	}{
		{"", "web API", "web-API.md"},
		{"my-API.md", "a collection name", "my-API.md"},
	}

	for _, test := range tests {
		t.Run(test.collectionName, func(t *testing.T) {
			destFile, destName, err := getDestFile(test.originalDestName, test.collectionName, false)
			if err != nil {
				t.Errorf(
					"getDestFile(%q, %q) = (%p, %q, %v), want nil error",
					test.originalDestName, test.collectionName, destFile, destName, err,
				)
				return
			}
			if destFile == os.Stdout {
				t.Errorf(
					"getDestFile(%q, %q) = (os.Stdout, %q, nil), want non-std file",
					test.originalDestName, test.collectionName, destName,
				)
				return
			}
			if destFile == os.Stdin {
				t.Errorf(
					"getDestFile(%q, %q) = (os.Stdin, %q, nil), want non-std file",
					test.originalDestName, test.collectionName, destName,
				)
				return
			}
			if destFile == os.Stderr {
				t.Errorf(
					"getDestFile(%q, %q) = (os.Stderr, %q, nil), want non-std file",
					test.originalDestName, test.collectionName, destName,
				)
				return
			}
			destFile.Close()
			defer os.Remove(destName)
			if destName != test.wantName {
				t.Errorf(
					"getDestFile(%q, %q) = (%p, %q, nil), want (%p, %q, nil)",
					test.originalDestName, test.collectionName, destFile, destName, destFile, test.wantName,
				)
			}
		})
	}
}

func TestGetDestFileWithEmptyNames(t *testing.T) {
	wantDestName := "collection.md"
	destFile, destName, err := getDestFile("", "", false)
	if err != nil || destName != wantDestName || destFile == nil {
		t.Errorf("getDestFile(\"\", \"\") = (%p, %q, %v), want (non-nil *os.File, %q, nil)", destFile, destName, err, wantDestName)
	}
	if destFile == os.Stdout {
		t.Error("getDestFile(\"\", \"\") returned os.Stdout, want non-std file pointer")
	} else if destFile == os.Stdin {
		t.Error("getDestFile(\"\", \"\") returned os.Stdin, want non-std file pointer")
	} else if destFile == os.Stderr {
		t.Error("getDestFile(\"\", \"\") returned os.Stderr, want non-std file pointer")
	} else if err == nil {
		destFile.Close()
		os.Remove(destName)
	}
}

func TestGetDestFileNameReplaceError(t *testing.T) {
	destFile, destName, err := getDestFile("samples/calendar-API-v1.md", "", false)
	if err == nil {
		t.Errorf("getDestFile targeting an existing file returned nil error, want non-nil error")
		t.Errorf("getDestFile(<existing file>, \"\") = (%p, %q, nil), want (nil, \"\", <non-nil error>)", destFile, destName)
		if destName != "-" {
			destFile.Close()
		}
	}
}

func TestExecuteTemplateWithInvalidTemplate(t *testing.T) {
	err := executeTemplate(nil, nil, "api v1", "# {{ .Name ")
	if err == nil {
		t.Errorf("executeTemplate(nil, nil, \"api v1\", \"# {{ .Name \") = nil, want non-nil error")
	}
}
