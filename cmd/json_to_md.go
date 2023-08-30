package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// jsonToMdFile converts JSON bytes to markdown, prints the markdown to a file or
// stdout, and returns the destination's name. If the destination name is "-", output
// goes to stdout. If the destination's name is empty, a file is created with a unique
// name based on the given JSON. Only an empty destination name will be changed from
// what is given before being returned.
func jsonToMdFile(jsonBytes []byte, destName, tmplName, tmplStr string, statusRanges [][]int, confirmReplaceExistingFile bool) (string, error) {
	collection, err := parseCollection(jsonBytes)
	if err != nil {
		return "", fmt.Errorf("parseCollection: %s", err)
	}
	filterResponsesByStatus(collection, statusRanges)
	if v, err := getVersion(collection.Routes); err == nil {
		collection.Info.Name += " " + v
	}

	destFile, destName, err := getDestFile(destName, collection.Info.Name, confirmReplaceExistingFile)
	if err != nil {
		return "", err
	}
	if destName != "-" {
		// destFile is not os.Stdout
		defer destFile.Close()
	}

	if err = executeTemplate(destFile, collection, tmplName, tmplStr); err != nil {
		return "", err
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

// getVersion returns the version number of a collection. If the collection's first
// route has a version number like `/v1/something`, then `v1` is returned. If no version
// number is found, an error is returned.
func getVersion(routes []Route) (string, error) {
	if len(routes) > 0 && len(routes[0].Request.Url.Path) > 0 {
		maybeVersion := routes[0].Request.Url.Path[0]
		if matched, err := regexp.Match(`v\d+`, []byte(maybeVersion)); err == nil && matched {
			return maybeVersion, nil
		}
	}
	return "", fmt.Errorf("No version number found")
}

// getDestFile gets the destination file and its name. If the given destination name is
// "-", the destination file is os.Stdout. If the given destination name is empty, a new
// file is created with a name based on the collection name and the returned name will
// be different from the given one. If the given destination name refers to an existing
// file and confirmation to replace an existing file is not given, an error is returned.
// Any returned file is open.
func getDestFile(destName, collectionName string, confirmReplaceExistingFile bool) (*os.File, string, error) {
	if destName == "-" {
		return os.Stdout, destName, nil
	}
	if len(destName) == 0 {
		fileName := FormatFileName(collectionName)
		if len(fileName) == 0 {
			fileName = "collection"
		}
		destName = CreateUniqueFileName(fileName, ".md")
	} else if FileExists(destName) && !confirmReplaceExistingFile {
		return nil, "", fmt.Errorf("File %q already exists. Run the command again with the --replace flag to confirm replacing it.", destName)
	}
	destFile, err := os.Create(destName)
	if err != nil {
		return nil, "", fmt.Errorf("os.Create: %s", err)
	}
	return destFile, destName, nil
}

// executeTemplate uses a template and FuncMap to convert the collection to markdown and
// saves to the given destination file. The destination file is not closed.
func executeTemplate(destFile *os.File, collection *Collection, tmplName, tmplStr string) error {
	tmpl, err := template.New(tmplName).Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("Template parsing error: %s", err)
	}

	return tmpl.Execute(destFile, collection)
}
