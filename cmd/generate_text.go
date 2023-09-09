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
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

// generateText converts a collection to plaintext and saves it into the given open file
// without closing the file. `Seek(0, 0)` is then called on the file so the file pointer
// is at the beginning of the file unless an error occurs. If the given template path is
// empty, the default template is used. If any status ranges are given, responses with
// statuses outside those ranges are removed from the collection. A `level` integer
// property is added to each "item" and each "response" object within the collection.
// The level starts at 1 for the outermost item object and increases by 1 for each level
// of item nesting.
func generateText(collection map[string]any, openAnsFile *os.File, tmplPath string, statusRanges [][]int) error {
	filterResponsesByStatus(collection, statusRanges)
	addLevelProperty(collection)

	tmplName, tmplStr, err := loadTmpl(tmplPath)
	if err != nil {
		return err
	}

	return executeTmpl(collection, openAnsFile, tmplName, tmplStr)
}

// parseCollection converts a collection from a slice of bytes of JSON to a map.
func parseCollection(jsonBytes []byte) (map[string]any, error) {
	var collection map[string]any
	if err := json.Unmarshal(jsonBytes, &collection); err != nil {
		return nil, err
	}
	if collection["info"].(map[string]any)["schema"] != "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" {
		return nil, fmt.Errorf("Unknown JSON schema. When exporting from Postman, export as Collection v2.1.0")
	}

	return collection, nil
}

// parseStatusRanges converts a string of status ranges to a slice of slices of
// integers. The slice may be nil, but any inner slices each have two elements: the
// start and end of the range. Example inputs: "200", "200-299", "200-299,400-499",
// "200-200".
func parseStatusRanges(statusesStr string) ([][]int, error) {
	if len(statusesStr) == 0 {
		return nil, nil
	}
	statusRangeStrs := strings.Split(statusesStr, ",")
	statusRanges := make([][]int, len(statusRangeStrs))
	for i, statusRangeStr := range statusRangeStrs {
		startAndEnd := strings.Split(statusRangeStr, "-")
		if len(startAndEnd) > 2 {
			return nil, fmt.Errorf("Invalid status format. There should be zero or one dashes in %s", statusRangeStr)
		}
		start, err := strconv.Atoi(startAndEnd[0])
		if err != nil {
			return nil, fmt.Errorf("Invalid status range format. Expected an integer, got %q", startAndEnd[0])
		}
		end := start
		if len(startAndEnd) > 1 {
			end, err = strconv.Atoi(startAndEnd[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid status range format. Expected an integer, got %q", startAndEnd[1])
			}
		}
		statusRanges[i] = make([]int, 2)
		statusRanges[i][0] = start
		statusRanges[i][1] = end
	}

	return statusRanges, nil
}

// filterResponsesByStatus removes all sample responses with status codes outside the
// given range(s). If no status ranges are given, the collection remains unchanged.
func filterResponsesByStatus(collection map[string]any, statusRanges [][]int) {
	if statusRanges == nil || len(statusRanges) == 0 {
		return
	}
	items := collection["item"].([]any)
	_filterResponsesByStatus(items, statusRanges)
}

func _filterResponsesByStatus(items []any, statusRanges [][]int) {
	for _, itemAny := range items {
		item := itemAny.(map[string]any)
		if subItemsAny, ok := item["item"]; ok { // if item is a folder
			_filterResponsesByStatus(subItemsAny.([]any), statusRanges)
		} else { // if item is an endpoint
			responses := item["response"].([]any)
			for j := len(responses) - 1; j >= 0; j-- {
				response := responses[j].(map[string]any)
				inRange := false
				for _, statusRange := range statusRanges {
					code := int(response["code"].(float64))
					if code >= statusRange[0] && code <= statusRange[1] {
						inRange = true
						break
					}
				}
				if !inRange {
					responses = slices.Delete(responses, j, j+1)
					item["response"] = responses
				}
			}
		}
	}
}

// addLevelProperty adds a "level" property within each "item" and each "response"
// object. The level starts at 1 for the outermost item object and increases by 1 for
// each level of item nesting.
func addLevelProperty(collection map[string]any) {
	items := collection["item"].([]any)
	_addLevelProperty(items, 1)
}

func _addLevelProperty(items []any, level int) {
	for _, itemAny := range items {
		item := itemAny.(map[string]any)
		item["level"] = level
		if subItemsAny, ok := item["item"]; ok { // if item is a folder
			_addLevelProperty(subItemsAny.([]any), level+1)
		} else { // if item is an endpoint
			responses := item["response"].([]any)
			for _, responseAny := range responses {
				response := responseAny.(map[string]any)
				response["level"] = level
			}
		}
	}
}

// executeTmpl uses a template and FuncMap to convert the collection to plaintext
// and saves to the given open destination file without closing it. `Seek(0, 0)` is then
// called on the file so the file pointer is at the beginning of the file unless an
// error occurs.
func executeTmpl(collection map[string]any, openAnsFile *os.File, tmplName, tmplStr string) error {
	tmpl, err := template.New(tmplName).Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("Template parsing error: %s", err)
	}

	err = tmpl.Execute(openAnsFile, collection)
	if err != nil {
		return err
	}

	openAnsFile.Seek(0, 0)
	return nil
}
