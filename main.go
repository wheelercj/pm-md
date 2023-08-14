package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// Checks if a given file or folder exists on the device.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// If an existing file has the given file name and extension, parentheses around a
// number are appended to the file name to make it unique. Otherwise, the given file
// name and extension are concatenated and returned unchanged. The extension must start
// with a period.
func CreateUniqueFileName(fileName, extension string) string {
	if !strings.HasPrefix(extension, ".") {
		panic("extension must start with a period")
	}
	uniqueFileName := fileName + extension
	for i := 1; i < 100 && FileExists(uniqueFileName); i++ {
		uniqueFileName = fileName + "(" + fmt.Sprint(i) + ")" + extension
	}
	return uniqueFileName
}

type Request struct {
	Method string `json:"method"`
	Header []any  `json:"header"`
	Body   struct {
		Mode    string `json:"mode"`
		Raw     string `json:"raw"`
		Options struct {
			Raw struct {
				Language string `json:"language"`
			} `json:"raw"`
		} `json:"options"`
	} `json:"body"`
	Url struct {
		Raw  string   `json:"raw"`
		Host []string `json:"host"`
		Path []string `json:"path"`
	} `json:"url"`
}

type Routes []struct {
	Name                    string `json:"name"`
	ProtocolProfileBehavior struct {
		DisableBodyPruning bool `json:"disableBodyPruning"`
	} `json:"protocolProfileBehavior"`
	Request   Request `json:"request"`
	Responses []struct {
		Name            string  `json:"name"`
		OriginalRequest Request `json:"originalRequest"`
		Status          string  `json:"status"`
		Code            int     `json:"code"`
		Language        string  `json:"_postman_previewlanguage"`
		Headers         []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"header"`
		Cookies []any  `json:"cookie"`
		Body    string `json:"body"`
	} `json:"response"`
}

type Collection struct {
	Info struct {
		PostmanId  string `json:"_postman_id"`
		Name       string `json:"name"`
		Schema     string `json:"schema"`
		ExporterId string `json:"_exporter_id"`
	} `json:"info"`
	Routes Routes `json:"item"`
	Events []struct {
		Listen string `json:"listen"`
		Script struct {
			Type string `json:"type"`
			Exec []string
		} `json:"script"`
	} `json:"event"`
	Variables []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"variable"`
}

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
