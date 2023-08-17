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
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"
)

//go:embed collection.tmpl
var tmplStr string

func main() {
	instructions := "usage: pm-md json_file\n\nYou can get the JSON file from Postman by exporting a collection as a v2.1.0 collection."
	if len(os.Args) == 1 {
		fmt.Println(instructions)
		os.Exit(0)
	}
	jsonFilePath := os.Args[1]
	if !strings.HasSuffix(jsonFilePath, ".json") && !strings.HasSuffix(jsonFilePath, ".JSON") {
		fmt.Println(instructions)
		os.Exit(0)
	}

	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	// fmt.Printf("file size: %v characters", len(fileContent))

	var collection Collection
	if err := json.Unmarshal(fileContent, &collection); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if collection.Info.Schema != "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" {
		fmt.Println("Error: unknown JSON schema. When exporting from Postman, export as Collection v2.1.0")
		os.Exit(1)
	}

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
