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

// Postman's JSON schema: https://schema.postman.com/json/collection/v2.1.0/collection.json

// As each endpoint in an API could have multiple requests in Postman, the `Endpoints`
// slice in the `Collection` struct may have multiple instances of some actual
// endpoints (possibly with different parameters and/or responses).

// Some assumptions are made about the JSON from Postman:
// * No folders are used.
// * Each request is an object, not a string.
// * Each request body is an object, not null.
// * Each response body is a string, not null.
// * Each header is an object, not a string nor null.
// * Each auth is an object, not null.

type Collection struct {
	Info                    `json:"info"`           // required
	Endpoints               []Endpoint              `json:"item"` // required
	Events                  []Event                 `json:"event"`
	Variables               []Variable              `json:"variable"`
	Auth                    Auth                    `json:"auth"`
	ProtocolProfileBehavior ProtocolProfileBehavior `json:"protocolProfileBehavior"`
}

type Info struct {
	Name        string `json:"name"`   // required
	Schema      string `json:"schema"` // required
	PostmanId   string `json:"_postman_id"`
	Description any    `json:"description"`
	Version     any    `json:"version"`
	ExporterId  string `json:"_exporter_id"`
}

type Endpoint struct {
	Request                 Request                 `json:"request"` // required
	Name                    string                  `json:"name"`
	Id                      string                  `json:"id"`
	Description             any                     `json:"description"`
	Responses               []Response              `json:"response"`
	Variables               []Variable              `json:"variable"`
	Events                  []Event                 `json:"event"`
	ProtocolProfileBehavior ProtocolProfileBehavior `json:"protocolProfileBehavior"`
}

type Request struct {
	Method      string      `json:"method"`
	Headers     []Header    `json:"header"`
	Body        RequestBody `json:"body"`
	Url         Url         `json:"url"`
	Description any         `json:"description"`
	Auth        Auth        `json:"auth"`
	Proxy       Proxy       `json:"proxy"`
	Certificate Certificate `json:"certificate"`
}

type RequestBody struct {
	Mode       string                `json:"mode"`
	Raw        string                `json:"raw"`
	UrlEncoded []UrlEncodedParameter `json:"urlencoded"`
	FormData   []FormParameter       `json:"formdata"`
	File       struct {
		Src     any    `json:"src"`
		Content string `json:"content"`
	} `json:"file"`
	GraphQL any `json:"graphql"`
	Options struct {
		Raw struct {
			Language string `json:"language"`
		} `json:"raw"`
	} `json:"options"`
	Disabled bool `json:"disabled"`
}

type Response struct {
	Name            string   `json:"name"`
	Id              string   `json:"id"`
	OriginalRequest Request  `json:"originalRequest"`
	Status          string   `json:"status"`
	Code            int      `json:"code"`
	Language        string   `json:"_postman_previewlanguage"`
	Headers         []Header `json:"header"`
	Cookies         []Cookie `json:"cookie"`
	Body            string   `json:"body"`
	ResponseTime    any      `json:"responseTime"`
	Timings         any      `json:"timings"`
}

type Url struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type Header struct {
	Key         string `json:"key"`   // required
	Value       string `json:"value"` // required
	Description any    `json:"description"`
	Disabled    bool   `json:"disabled"`
}

type UrlEncodedParameter struct {
	Key         string `json:"key"` // required
	Value       string `json:"value"`
	Description any    `json:"description"`
	Disabled    bool   `json:"disabled"`
}

type FormParameter struct {
	Key         string `json:"key"` // required
	Value       string `json:"value"`
	Description any    `json:"description"`
	Disabled    bool   `json:"disabled"`
	Src         any    `json:"src"`
	Type        string `json:"type"`
	ContentType string `json:"contentType"`
}

type Cookie struct {
	Domain     string `json:"domain"` // required
	Path       string `json:"path"`   // required
	Name       string `json:"name"`
	Expires    any    `json:"expires"`
	MaxAge     string `json:"maxAge"`
	Value      string `json:"value"`
	Secure     bool   `json:"secure"`
	Session    bool   `json:"session"`
	HostOnly   bool   `json:"hostOnly"`
	HttpOnly   bool   `json:"httpOnly"`
	Extensions []any  `json:"extensions"`
}

type Auth struct {
	Type     string          `json:"type"` // required
	ApiKey   []AuthAttribute `json:"apikey"`
	AwsV4    []AuthAttribute `json:"awsv4"`
	Basic    []AuthAttribute `json:"basic"`
	Bearer   []AuthAttribute `json:"bearer"`
	Digest   []AuthAttribute `json:"digest"`
	EdgeGrid []AuthAttribute `json:"edgegrid"`
	Hawk     []AuthAttribute `json:"hawk"`
	Ntlm     []AuthAttribute `json:"ntlm"`
	OAuth1   []AuthAttribute `json:"oauth1"`
	OAuth2   []AuthAttribute `json:"oauth2"`
}

type AuthAttribute struct {
	Key   string `json:"key"` // required
	Value any    `json:"value"`
	Type  string `json:"type"`
}

type Event struct {
	Listen string `json:"listen"`
	Script struct {
		Type string   `json:"type"`
		Exec []string `json:"exec"`
	} `json:"script"`
}

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type ProtocolProfileBehavior struct {
	DisableBodyPruning bool `json:"disableBodyPruning"`
}

type Proxy struct {
	Match    string `json:"match"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Tunnel   bool   `json:"tunnel"`
	Disabled bool   `json:""`
}

type Certificate struct {
	Name    string   `json:"name"`
	Matches []string `json:"matches"`
	Key     struct {
		Src any `json:"src"`
	} `json:"key"`
	Cert struct {
		Src any `json:"src"`
	} `json:"cert"`
	Passphrase string `json:"passphrase"`
}
