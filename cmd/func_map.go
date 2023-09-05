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
	"fmt"
	"html/template"
	"slices"
	"strings"
)

var funcMap = template.FuncMap{
	"formatHeaderLink": formatHeaderLink,
	"add": func(a, b int) int {
		return a + b
	},
	"join": func(elems []any, sep string) string {
		strElems := make([]string, len(elems))
		for i, e := range elems {
			strElems[i] = fmt.Sprint(e)
		}
		return strings.Join(strElems, sep)
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

var headerPathCache = make([]string, 0, 10)

// formatHeaderLink formats a markdown header body as a markdown link to the header
// compatible with GitHub's markdown rendering. When GitHub and this function find
// duplicate headers, they append `-1` to the header link for the second occurence, `-2`
// for the third, and so on.
func formatHeaderLink(headerBody string) string {
	headerPath := formatHeaderPath(headerBody)
	uniqueHeaderPath := headerPath
	for i := 1; slices.Contains(headerPathCache, uniqueHeaderPath); i++ {
		uniqueHeaderPath = fmt.Sprintf("%s-%d", headerPath, i)
	}
	headerPathCache = append(headerPathCache, uniqueHeaderPath)
	return fmt.Sprintf("[%s](%s)", headerBody, uniqueHeaderPath)
}

// formatHeaderPath formats a markdown header body as a relative link path compatible
// with GitHub's markdown rendering. Letters are lowercased, leading spaces are removed,
// remaining spaces are replaced with dashes, special characters except dashes and
// underscores are removed, and one `#` will be prepended. Current limitation:
// formatHeaderPath ignores all emoji whereas GitHub removes some emoji.
func formatHeaderPath(headerBody string) string {
	headerBody = strings.ReplaceAll(
		strings.TrimLeft(strings.ToLower(headerBody), " "), " ", "-",
	)
	toRemove := "=+!@#$%^&*()|\\'\";:/?.,<>[]{}`~"
	result := make([]rune, 0, len(headerBody)/2)
	result = append(result, '#')
	for _, ch := range headerBody {
		if !strings.Contains(toRemove, string(ch)) {
			result = append(result, ch)
		}
	}

	return string(result)
}
