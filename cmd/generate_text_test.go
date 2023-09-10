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
	"os"
	"reflect"
	"strings"
	"testing"
)

// assertGenerateNoDiff asserts the given JSON and template result in expected
// plaintext. If the given template path is empty, the default template is used.
// wantPath is the path to an existing file containing the wanted output.
func assertGenerateNoDiff(t *testing.T, jsonPath, tmplPath, wantPath string) {
	// Skip this test if unique file name creation isn't working correctly.
	TestCreateUniqueFileName(t)
	TestCreateUniqueFileNamePanic(t)
	if t.Failed() {
		return
	}

	err := AssertGenerateNoDiff(jsonPath, tmplPath, wantPath, nil)
	if err != nil {
		t.Error(err)
	}
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

func TestGenerateText(t *testing.T) {
	inputPath := "../samples/calendar-API.postman_collection.json"
	wantOutputPath := "../samples/calendar-API-v1.md"
	assertGenerateNoDiff(t, inputPath, "", wantOutputPath)
}

func TestGenerateTextWithCustomTemplate(t *testing.T) {
	inputPath := "../samples/minimal-calendar-API.postman_collection.json"
	customTmplPath := "../samples/custom.tmpl"
	wantOutputPath := "../samples/custom-calendar-API-v1.md"
	assertGenerateNoDiff(t, inputPath, customTmplPath, wantOutputPath)
}

func TestGenerateTextWithMinimalTemplate(t *testing.T) {
	inputPath := "../samples/minimal-calendar-API.postman_collection.json"
	customTmplPath := "minimal.tmpl"
	wantOutputPath := "../samples/minimal-calendar-API-v1.md"
	assertGenerateNoDiff(t, inputPath, customTmplPath, wantOutputPath)
}

func TestParseCollectionWithInvalidJson(t *testing.T) {
	invalidJson := []byte(`
		{
			"info": {
				"_postman_id": "23799766-64ba-4c7c-aaa9-0d880964db54",
				"name": "calendar API",
				"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
				"_exporter_id": "23363106"
			},
	`)

	_, err := parseCollection(invalidJson)
	if err == nil {
		t.Error("Error expected")
	}
}

func TestParseCollectionWithOldSchema(t *testing.T) {
	inputPath := "../samples/calendar-API.postman_collection.json"
	jsonBytes, err := os.ReadFile(inputPath)
	if err != nil {
		t.Error(err)
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

// getCollection loads JSON from the file at the given path and converts the JSON to a
// map.
func getCollection(t *testing.T, jsonPath string) (map[string]any, error) {
	jsonBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	collection, err := parseCollection(jsonBytes)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// assertAllStatuses200 asserts that every "response" object in the given items has a
// status code of 200.
func assertAllStatuses200(t *testing.T, items []any) {
	for _, itemAny := range items {
		item := itemAny.(map[string]any)
		if subItemsAny, ok := item["item"]; ok { // if item is a folder
			assertAllStatuses200(t, subItemsAny.([]any))
		} else { // if item is an endpoint
			for _, responseAny := range item["response"].([]any) {
				response := responseAny.(map[string]any)
				code := int(response["code"].(float64))
				if code != 200 {
					t.Errorf("want 200, got %d", code)
				}
			}
		}
	}
}

// assertLevels asserts that each "item" and "response" object has a "level" integer
// property, and that it has the expected value.
func assertLevels(t *testing.T, items []any, wantLevel int) {
	for _, itemAny := range items {
		item := itemAny.(map[string]any)
		if ansLevel, ok := item["level"]; !ok {
			t.Errorf("Item %q at level %d has no \"level\" property", item["name"], wantLevel)
		} else if ansLevel != wantLevel {
			t.Errorf("Item %q has level %d, want level %d", item["name"], ansLevel, wantLevel)
		}
		if subItemsAny, ok := item["item"]; ok { // if item is a folder
			assertLevels(t, subItemsAny.([]any), wantLevel+1)
		} else { // if item is an endpoint
			for _, responseAny := range item["response"].([]any) {
				response := responseAny.(map[string]any)
				if ansLevel, ok := response["level"]; !ok {
					t.Errorf("Endpoint %q at level %d has no \"level\" property", item["name"], wantLevel)
				} else if ansLevel != wantLevel {
					t.Errorf("Endpoint %q has level %d, want level %d", item["name"], ansLevel, wantLevel)
				}
			}
		}
	}
}

func TestFilterResponses(t *testing.T) {
	jsonPath := "../samples/calendar-API.postman_collection.json"
	collection, err := getCollection(t, jsonPath)
	if err != nil {
		t.Error(err)
		return
	}

	filterResponsesByStatus(collection, [][]int{{200, 200}})
	items := collection["item"].([]any)
	assertAllStatuses200(t, items)
}

func TestFilterResponsesWithFolders(t *testing.T) {
	jsonPath := "../samples/calendar-API.postman_collection.json"
	collection, err := getCollection(t, jsonPath)
	if err != nil {
		t.Error(err)
		return
	}

	filterResponsesByStatus(collection, [][]int{{200, 200}})
	items := collection["item"].([]any)
	assertAllStatuses200(t, items)
}

func TestAddLevelProperty(t *testing.T) {
	jsonPath := "../samples/calendar-API.postman_collection.json"
	collection, err := getCollection(t, jsonPath)
	if err != nil {
		t.Error(err)
		return
	}

	addLevelProperty(collection)
	items := collection["item"].([]any)
	assertLevels(t, items, 1)
}

func TestGetDestFileStdout(t *testing.T) {
	destFile, destName, err := openDestFile("-", "", false)
	if destFile != os.Stdout || destName != "-" || err != nil {
		t.Errorf("openDestFile(\"-\", \"\") = (%p, %q, %q), want (%p, \"-\", nil)", destFile, destName, err, os.Stdout)
	}
}

func TestGetDestFileExistingFileErr(t *testing.T) {
	destFile, destName, err := openDestFile("../LICENSE", "", false)
	if err == nil {
		t.Errorf("openDestFile(\"../LICENSE\", \"\", false) = (%p, %q, nil), want non-nil error", destFile, destName)
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
			destFile, destName, err := openDestFile(test.originalDestName, test.collectionName, false)
			if err != nil {
				t.Errorf(
					"openDestFile(%q, %q) = (%p, %q, %v), want nil error",
					test.originalDestName, test.collectionName, destFile, destName, err,
				)
				return
			}
			if destFile == os.Stdout {
				t.Errorf(
					"openDestFile(%q, %q) = (os.Stdout, %q, nil), want non-std file",
					test.originalDestName, test.collectionName, destName,
				)
				return
			}
			if destFile == os.Stdin {
				t.Errorf(
					"openDestFile(%q, %q) = (os.Stdin, %q, nil), want non-std file",
					test.originalDestName, test.collectionName, destName,
				)
				return
			}
			if destFile == os.Stderr {
				t.Errorf(
					"openDestFile(%q, %q) = (os.Stderr, %q, nil), want non-std file",
					test.originalDestName, test.collectionName, destName,
				)
				return
			}
			destFile.Close()
			defer os.Remove(destName)
			if destName != test.wantName {
				t.Errorf(
					"openDestFile(%q, %q) = (%p, %q, nil), want (%p, %q, nil)",
					test.originalDestName, test.collectionName, destFile, destName, destFile, test.wantName,
				)
			}
		})
	}
}

func TestGetDestFileWithEmptyNames(t *testing.T) {
	wantDestName := "collection.md"
	destFile, destName, err := openDestFile("", "", false)
	if err != nil || destName != wantDestName || destFile == nil {
		t.Errorf("openDestFile(\"\", \"\") = (%p, %q, %v), want (non-nil *os.File, %q, nil)", destFile, destName, err, wantDestName)
	}
	if destFile == os.Stdout {
		t.Error("openDestFile(\"\", \"\") returned os.Stdout, want non-std file pointer")
	} else if destFile == os.Stdin {
		t.Error("openDestFile(\"\", \"\") returned os.Stdin, want non-std file pointer")
	} else if destFile == os.Stderr {
		t.Error("openDestFile(\"\", \"\") returned os.Stderr, want non-std file pointer")
	} else if err == nil {
		destFile.Close()
		os.Remove(destName)
	}
}

func TestGetDestFileNameReplaceError(t *testing.T) {
	destFile, destName, err := openDestFile("samples/calendar-API-v1.md", "", false)
	if err == nil {
		t.Errorf("openDestFile targeting an existing file returned nil error, want non-nil error")
		t.Errorf("openDestFile(<existing file>, \"\") = (%p, %q, nil), want (nil, \"\", <non-nil error>)", destFile, destName)
		if destName != "-" {
			destFile.Close()
		}
	}
}

func TestExecuteTmplWithInvalidTemplate(t *testing.T) {
	err := executeTmpl(nil, nil, "api v1", "# {{ .Name ")
	if err == nil {
		t.Errorf("executeTmpl(nil, nil, \"api v1\", \"# {{ .Name \") = nil, want non-nil error")
	}
}
