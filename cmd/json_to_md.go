package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

//go:embed collection.tmpl
var tmplStr string

// jsonToMdFile converts JSON bytes into markdown, saves the markdown into a file, and
// returns the new markdown file's name. The new file is guaranteed to not replace an
// existing file. The file's name is based on the contents of the given JSON.
func jsonToMdFile(jsonBytes []byte, statusRanges [][]int) (mdFileName string, err error) {
	collection, err := parseCollection(jsonBytes)
	if err != nil {
		return "", err
	}
	filterResponses(collection, statusRanges)
	if v, err := getVersion(collection.Routes); err == nil {
		collection.Info.Name += " " + v
	}
	mdFileName = CreateUniqueFileName(collection.Info.Name, ".md")
	mdFile, err := os.Create(mdFileName)
	if err != nil {
		return "", err
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
		// "assumeSafeHtml": func(s string) template.HTML {
		// 	// This prevents HTML escaping. Never run this with untrusted input.
		// 	return template.HTML(s)
		// },
	}

	tmplFileName := "collection.tmpl"
	tmpl, err := template.New(tmplFileName).Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(mdFile, collection)
	if err != nil {
		return "", err
	}

	return mdFileName, nil
}

// parseCollection converts a collection from a slice of bytes of JSON to a Collection
// instance.
func parseCollection(jsonBytes []byte) (*Collection, error) {
	var collection Collection
	if err := json.Unmarshal(jsonBytes, &collection); err != nil {
		return nil, err
	}
	if collection.Info.Schema != "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" {
		return nil, fmt.Errorf("Unknown JSON schema. When exporting from Postman, export as Collection v2.1.0")
	}

	return &collection, nil
}

// parseStatusRanges converts a string of status ranges to a slice of slices of
// integers. The slice may be nil, but any inner slices each have two elements: the
// start and end of the range. Example ranges: "200-299", "200-299,400-499", "200-200".
func parseStatusRanges(statusesStr string) ([][]int, error) {
	if len(statusesStr) == 0 {
		return nil, nil
	}
	statusRangeStrs := strings.Split(statusesStr, ",")
	statusRanges := make([][]int, len(statusRangeStrs))
	for i, statusRangeStr := range statusRangeStrs {
		startAndEnd := strings.Split(statusRangeStr, "-")
		if len(startAndEnd) != 2 {
			return nil, fmt.Errorf("Invalid status range format. There should be one dash (-) per range.")
		}
		start, err := strconv.Atoi(startAndEnd[0])
		if err != nil {
			return nil, fmt.Errorf("Invalid status range format. Expected an integer, got %q", startAndEnd[0])
		}
		end, err := strconv.Atoi(startAndEnd[1])
		if err != nil {
			return nil, fmt.Errorf("Invalid status range format. Expected an integer, got %q", startAndEnd[1])
		}
		statusRanges[i] = make([]int, 2)
		statusRanges[i][0] = start
		statusRanges[i][1] = end
	}

	return statusRanges, nil
}

// filterResponses removes all sample responses with status codes outside the given
// range(s). If no status ranges are given, the collection remains unchanged.
func filterResponses(collection *Collection, statusRanges [][]int) {
	if statusRanges == nil || len(statusRanges) == 0 {
		return
	}
	for i, route := range collection.Routes {
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
			collection.Routes[i] = route
		}
	}
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
