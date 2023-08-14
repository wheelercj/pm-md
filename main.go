package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	flag.Parse()
	jsonFilePath := flag.Arg(0)
	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("file size: %d characters\n", len(fileContent))

	var collection Collection
	if err := json.Unmarshal(fileContent, &collection); err != nil {
		panic(err)
	}
	if collection.Info.Schema != "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" {
		log.Fatal("When exporting from Postman, export as Collection v2.1.0")
	}

	routes := collection.Routes
	versionedCollectionName := getVersionedCollectionName(collection.Info.Name, routes)
	mdFileName := CreateUniqueFileName(versionedCollectionName, ".md")
	mdFile, err := os.Create(mdFileName)
	if err != nil {
		panic(err)
	}
	defer mdFile.Close()

	mdFile.Write(toMarkdown(versionedCollectionName, routes))
}

// If the collection's first route has a version number like `/v1/something`, then ` v1`
// is appended to the collection's name. Otherwise, the collection's name is returned
// unchanged.
func getVersionedCollectionName(baseCollectionName string, routes Routes) string {
	if len(routes) > 0 && len(routes[0].Request.Url.Path) > 0 {
		maybeVersion := routes[0].Request.Url.Path[0]
		if matched, err := regexp.Match(`v\d+`, []byte(maybeVersion)); err == nil && matched {
			return baseCollectionName + " " + maybeVersion
		}
	}
	return baseCollectionName
}

func toMarkdown(collectionName string, routes Routes) []byte {
	var mdBuffer bytes.Buffer
	mdBuffer.WriteString("# " + collectionName + "\n\n")
	for _, route := range routes {
		mdBuffer.WriteString("----------------------------------------\n\n## ")
		mdBuffer.WriteString(route.Name + "\n\n")
		url := "/" + strings.Join(route.Request.Url.Path, "/")
		mdBuffer.WriteString(route.Request.Method + " `" + url + "`\n\n")
		for _, response := range route.Responses {
			mdBuffer.WriteString("### sample response (status: ")
			mdBuffer.WriteString(fmt.Sprint(response.Code) + " ")
			mdBuffer.WriteString(response.Status + ")\n\n```")
			mdBuffer.WriteString(response.Language + "\n" + response.Body)
			mdBuffer.WriteString("\n```\n\n")
		}
	}
	return mdBuffer.Bytes()
}
