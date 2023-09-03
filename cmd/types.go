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
	"html/template"
	"strings"
)

var funcMap = template.FuncMap{
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
	"formatHeaderLink": formatHeaderLink,
}

// As each endpoint in an API could have multiple requests in Postman, the endpoints
// slice below may have "duplicates".

type Collection struct {
	Info struct {
		PostmanId  string `json:"_postman_id"`
		Name       string `json:"name"`
		Schema     string `json:"schema"`
		ExporterId string `json:"_exporter_id"`
	} `json:"info"`
	Endpoints []Endpoint `json:"item"`
	Events    []struct {
		Listen string `json:"listen"`
		Script struct {
			Type string   `json:"type"`
			Exec []string `json:"exec"`
		} `json:"script"`
	} `json:"event"`
	Variables []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"variable"`
}

type Endpoint struct {
	Name                    string `json:"name"`
	ProtocolProfileBehavior struct {
		DisableBodyPruning bool `json:"disableBodyPruning"`
	} `json:"protocolProfileBehavior"`
	Request   Request    `json:"request"`
	Responses []Response `json:"response"`
}

type Request struct {
	Method  string `json:"method"`
	Headers []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"header"`
	Body struct {
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

type Response struct {
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
}
