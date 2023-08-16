package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

//go:embed collection.tmpl
var tmplStr string

func main() {
	flag.Parse()
	jsonFilePath := flag.Arg(0)
	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		panic(err)
	}
	slog.Debug("file size", "characters", len(fileContent))

	var collection Collection
	if err := json.Unmarshal(fileContent, &collection); err != nil {
		panic(err)
	}
	if collection.Info.Schema != "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" {
		log.Fatal("When exporting from Postman, export as Collection v2.1.0")
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

// If the collection's first route has a version number like `/v1/something`, then `v1`
// is returned. If no version number is found, an error is returned.
func getVersion(routes Routes) (string, error) {
	if len(routes) > 0 && len(routes[0].Request.Url.Path) > 0 {
		maybeVersion := routes[0].Request.Url.Path[0]
		if matched, err := regexp.Match(`v\d+`, []byte(maybeVersion)); err == nil && matched {
			return maybeVersion, nil
		}
	}
	return "", fmt.Errorf("No version number found")
}
