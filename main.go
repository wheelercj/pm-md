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
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

//go:embed collection.tmpl
var tmplStr string

const version string = "v0.0.4"

func main() {
	os.Args[0] = strings.Split(path.Base(strings.Replace(os.Args[0], "\\", "/", -1)), ".")[0]
	help := fmt.Sprintf(
		"usage: %s json_file\n\nYou can get the JSON file from Postman by exporting a"+
			" collection as a v2.1.0 collection.\n\n%s %s  You can check for updates here:"+
			" https://github.com/wheelercj/pm-md/releases", os.Args[0], os.Args[0], version)

	versionFlag := flag.Bool("version", false, "See this app's version")
	statusesFlag := flag.String(
		"statuses",
		"",
		"Include only the sample responses within the given range(s) of status codes."+
			" Example ranges: \"200-299\" or \"200-299,400-499\"")

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(help)
		os.Exit(0)
	}

	statusRanges := parseStatusRanges(*statusesFlag)

	jsonFilePath := args[0]
	if !strings.HasSuffix(jsonFilePath, ".json") && !strings.HasSuffix(jsonFilePath, ".JSON") {
		fmt.Println(help)
		os.Exit(0)
	}

	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	var collection Collection
	if err := json.Unmarshal(fileContent, &collection); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if collection.Info.Schema != "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" {
		fmt.Println("Error: unknown JSON schema. When exporting from Postman, export as Collection v2.1.0")
		os.Exit(1)
	}
	collection = filterResponses(collection, statusRanges)

	routes := collection.Routes
	if v, err := getVersion(routes); err == nil {
		collection.Info.Name += " " + v
	}
	mdFileName := CreateUniqueFileName(collection.Info.Name, ".md")
	mdFile, err := os.Create(mdFileName)
	if err != nil {
		panic(err)
	}
	defer mdFile.Close()

	funcMap := template.FuncMap{
		"join": func(elems []string, sep string) string {
			return strings.Join(elems, sep)
		},
		"allowJsonOrPlaintext": func(s string) any {
			if json.Valid([]byte(s)) {
				return template.HTML(s)
			}
			return s
		},
		"assumeSafeHtml": func(s string) template.HTML {
			// This prevents HTML escaping. Never run this with untrusted input.
			return template.HTML(s)
		},
	}

	tmplFileName := "collection.tmpl"
	tmpl, err := template.New(tmplFileName).Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(mdFile, collection)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created", mdFileName)
}

// parseStatusRanges converts a string of status ranges to a slice of slices of
// integers. The slice may be empty, but any inner slices each have two elements: the
// start and end of the range.
func parseStatusRanges(statusesStr string) [][]int {
	statusRangeStrs := strings.Split(statusesStr, ",")
	statusRanges := make([][]int, len(statusRangeStrs))
	for i, statusRangeStr := range statusRangeStrs {
		startAndEnd := strings.Split(statusRangeStr, "-")
		if len(startAndEnd) != 2 {
			fmt.Println("Error: invalid status range format. There should be one dash (-) per range.")
			os.Exit(1)
		}
		start, err := strconv.Atoi(startAndEnd[0])
		if err != nil {
			fmt.Println("Error: invalid status range format. Expected an integer, got", startAndEnd[0])
			os.Exit(1)
		}
		end, err := strconv.Atoi(startAndEnd[1])
		if err != nil {
			fmt.Println("Error: invalid status range format. Expected an integer, got", startAndEnd[1])
			os.Exit(1)
		}
		statusRanges[i] = make([]int, 2)
		statusRanges[i][0] = start
		statusRanges[i][1] = end
	}
	return statusRanges
}

// filterResponses removes all sample responses with status codes outside the given
// range(s). If no status ranges are given, the collection remains unchanged.
func filterResponses(collection Collection, statusRanges [][]int) Collection {
	if len(statusRanges) == 0 {
		return collection
	}
	for _, route := range collection.Routes {
		for j := len(route.Responses) - 1; j >= 0; j-- {
			response := route.Responses[j]
			inRange := false
			for _, statusRange := range statusRanges {
				if response.Code >= statusRange[0] && response.Code <= statusRange[1] {
					inRange = true
					break
				}
			}
			if !inRange {
				route.Responses = slices.Delete(route.Responses, j, j+1)
			}
		}
	}
	return collection
}

// getVersion returns the version number of a collection. If the collection's first
// route has a version number like `/v1/something`, then `v1` is returned. If no version
// number is found, an error is returned.
func getVersion(routes Routes) (string, error) {
	if len(routes) > 0 && len(routes[0].Request.Url.Path) > 0 {
		maybeVersion := routes[0].Request.Url.Path[0]
		if matched, err := regexp.Match(`v\d+`, []byte(maybeVersion)); err == nil && matched {
			return maybeVersion, nil
		}
	}
	return "", fmt.Errorf("No version number found")
}
