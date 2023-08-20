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
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseStatusRanges(t *testing.T) {
	tests := []struct {
		input string
		want  [][]int
	}{
		{"", nil},
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
	inputs := []string{"200-299-300", "a-299", "200-b", "200-", "-299", "200", "-"}
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
	// Skip this test if unique file name creation isn't working.
	TestCreateUniqueFileName(t)
	TestCreateUniqueFileNamePanic(t)
	if t.Failed() {
		return
	}

	inputFilePath := "samples/calendar API.postman_collection.json"
	wantFilePath := "samples/calendar API v1.md"
	jsonBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		t.Errorf("Failed to open %s", inputFilePath)
		return
	}
	mdFileName, err := jsonToMdFile(jsonBytes, nil)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(mdFileName)
	ansBytes, err := os.ReadFile(mdFileName)
	if err != nil {
		t.Errorf("Failed to open %s", mdFileName)
		return
	}
	wantBytes, err := os.ReadFile(wantFilePath)
	if err != nil {
		t.Errorf("Failed to open %s", wantFilePath)
		return
	}
	ans := strings.ReplaceAll(string(ansBytes), "\r\n", "\n")
	want := strings.ReplaceAll(string(wantBytes), "\r\n", "\n")

	assertNoDiff(t, ans, want, "\n")
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
	mdFileName, err := jsonToMdFile(invalidJson, nil)
	if err == nil {
		t.Error("Error expected")
		os.Remove(mdFileName)
	}
}

func TestParseCollectionWithOldSchema(t *testing.T) {
	inputFilePath := "samples/calendar API.postman_collection.json"
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

func TestFilterResponses(t *testing.T) {
	inputFilePath := "samples/calendar API.postman_collection.json"
	jsonBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		t.Errorf("Failed to open %s", inputFilePath)
		return
	}

	collection, err := parseCollection(jsonBytes)
	if err != nil {
		t.Error(err)
		return
	}

	filterResponses(collection, [][]int{{200, 200}})
	for _, route := range collection.Routes {
		for _, response := range route.Responses {
			if response.Code != 200 {
				t.Errorf("want 200, got %d", response.Code)
				return
			}
		}
	}
}
