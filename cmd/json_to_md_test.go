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
	"strconv"
	"strings"
	"testing"
)

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
	assertJsonToMdFileNoDiff(t, inputFilePath, "", "", wantOutputPath, false)
}

func TestJsonToMdFileWithResponseNames(t *testing.T) {
	inputFilePath := "../samples/calendar-API.postman_collection.json"
	wantOutputPath := "../samples/calendar-API-v1-with-response-names.md"
	assertJsonToMdFileNoDiff(t, inputFilePath, "", "", wantOutputPath, true)
}

func TestJsonToMdFileWithCustomOutputFileName(t *testing.T) {
	inputFilePath := "../samples/calendar-API.postman_collection.json"
	customOutputFileName := "custom-file-name-for-testing.md"
	wantOutputPath := "../samples/calendar-API-v1-with-response-names.md"
	assertJsonToMdFileNoDiff(t, inputFilePath, "", customOutputFileName, wantOutputPath, true)
}

func TestJsonToMdFileWithCustomTemplate(t *testing.T) {
	inputFilePath := "../samples/calendar-API.postman_collection.json"
	customTmplPath := "../samples/custom.tmpl"
	wantOutputPath := "../samples/calendar-API-v1-from-custom-templ.md"
	assertJsonToMdFileNoDiff(t, inputFilePath, customTmplPath, "", wantOutputPath, true)
}

func TestInvalidJsonToMdFile(t *testing.T) {
	// Skip this test if unique file name creation isn't working.
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
	destName, err := jsonToMdFile(invalidJson, "-", "", nil, false)
	if err == nil {
		t.Error("Error expected")
		if destName != "-" {
			t.Errorf("The destination name should not have changed from \"-\" to %q", destName)
		}
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
		t.Errorf("want (nil, error), got a nil error and a non-nil collection: %v", *collection)
	}
}

func getCollection(t *testing.T) (*Collection, error) {
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
	for _, route := range collection.Routes {
		for _, response := range route.Responses {
			if response.Code != 200 {
				t.Errorf("want 200, got %d", response.Code)
				return
			}
		}
	}
}

func TestGetVersionWithoutVersionedRoutes(t *testing.T) {
	collection, err := getCollection(t)
	if err != nil {
		t.Error(err)
		return
	}
	if len(collection.Routes) == 0 {
		t.Error("No routes to test")
		return
	}

	for i, route := range collection.Routes {
		if len(route.Request.Url.Path) == 0 {
			t.Errorf("Request missing path: %v", route.Request)
			return
		}

		// Delete any version number from the route.
		maybeVersion := route.Request.Url.Path[0]
		if strings.HasPrefix(maybeVersion, "v") {
			maybeNumber := strings.TrimPrefix(maybeVersion, "v")
			if _, err := strconv.Atoi(maybeNumber); err == nil {
				collection.Routes[i].Request.Url.Path = route.Request.Url.Path[1:]
			}
		}
	}

	version, err := getVersion(collection.Routes)
	if err == nil {
		t.Errorf("getVersion returned (%q, nil), want non-nil error", version)
	}
}

func TestGetDestFileStdout(t *testing.T) {
	destFile, destName, err := getDestFile("-", "")
	if destFile != os.Stdout || destName != "-" || err != nil {
		t.Errorf("getDestFile(\"-\", \"\") = (%p, %q, %q), want (%p, \"-\", nil)", destFile, destName, err, os.Stdout)
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
			destFile, destName, err := getDestFile(test.originalDestName, test.collectionName)
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
