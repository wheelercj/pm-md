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

//go:embed collection_default.tmpl
var defaultTmplStr string
var defaultTmplName = "collection_default.tmpl"

// jsonToMdFile converts JSON bytes to markdown, prints the markdown to a file or
// stdout, and returns the destination's name. If the destination name is "-", output
// goes to stdout. If the destination's name is empty, a file is created with a unique
// name based on the given JSON. Only an empty destination name will be changed from
// what is given before being returned.
func jsonToMdFile(jsonBytes []byte, destName string, statusRanges [][]int, showResponseNames bool) (string, error) {
	collection, err := parseCollection(jsonBytes)
	if err != nil {
		return "", fmt.Errorf("parseCollection: %s", err)
	}
	filterResponsesByStatus(collection, statusRanges)
	if !showResponseNames {
		clearResponseNames(collection)
	}
	if v, err := getVersion(collection.Routes); err == nil {
		collection.Info.Name += " " + v
	}

	var destFile *os.File
	if len(destName) == 0 {
		destName = CreateUniqueFileName(FormatFileName(collection.Info.Name), ".md")
	} else if destName == "-" {
		destFile = os.Stdout
	} else if FileExists(destName) {
		if err := ConfirmReplaceExistingFile(destName); err != nil {
			return "", fmt.Errorf("ConfirmReplaceExistingFile: %s", err)
		}
	}

	if destFile == nil {
		destFile, err = os.Create(destName)
		if err != nil {
			return "", fmt.Errorf("os.Create: %s", err)
		}
		defer destFile.Close()
	}

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

	tmpl, err := template.New(defaultTmplName).Funcs(funcMap).Parse(defaultTmplStr)
	if err != nil {
		return "", fmt.Errorf("*Template.Parse: %s", err)
	}
	err = tmpl.Execute(destFile, collection)
	if err != nil {
		return "", fmt.Errorf("tmpl.Execute: %s", err)
	}

	return destName, nil
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
// start and end of the range. Examples: "200", "200-299", "200-299,400-499", "200-200".
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
func filterResponsesByStatus(collection *Collection, statusRanges [][]int) {
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

// clearResponseNames changes each response name to an empty string. This is helpful
// when response names are not wanted in the output.
func clearResponseNames(collection *Collection) {
	for i := range collection.Routes {
		for j := range collection.Routes[i].Responses {
			collection.Routes[i].Responses[j].Name = ""
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
