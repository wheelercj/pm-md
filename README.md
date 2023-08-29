# Postman to Markdown

[![Go Reference](https://pkg.go.dev/badge/github.com/wheelercj/pm-md.svg)](https://pkg.go.dev/github.com/wheelercj/pm-md)

Convert a Postman collection to markdown documentation.

[sample result](samples/calendar-API-v1.md)

**Template-driven**. You can easily customize the output by editing the template. See the "custom templates" section below for more details.

## download

3 choices for how to download:

* [download a zipped executable file](https://github.com/wheelercj/pm-md/releases), unzip it, and run the app with `./pm-md --help`
* `go install github.com/wheelercj/pm-md@latest` and then `pm-md --help`
* install from source following the instructions below

### install from source

These steps require [Go](https://go.dev/) to be installed.

1. Choose one of the source code download options [here](https://github.com/wheelercj/pm-md/releases) and unzip the folder (or use `git clone`).
2. Open a terminal in the new folder.
3. Run `go build` to create an executable file.
4. Run `go install` to install the executable file. If you get an error message, you may need to [edit your PATH environment variable](https://go.dev/doc/tutorial/compile-install).
5. Run the app with `pm-md --help`.

## command examples

* `pm-md collection.json documentation.md` reads collection.json and saves markdown to documentation.md.
* `pm-md collection.json` reads collection.json and saves markdown to a new file with a unique name based on the collection's name. This will NEVER replace an existing file.
* `pm-md collection.json --statuses=200` does the same as the previous example but does not include any sample responses except those with a status code of 200.
* `pm-md collection.json --statuses=200-299,400-499` does not include any sample responses except those with a status code within the ranges 200-299 and 400-499 (inclusive).
* `pm-md collection.json -` reads collection.json and returns markdown to stdout.
* `pm-md - -` receives JSON from stdin and returns markdown to stdout, such as with `cat collection.json | pm-md - -`.
* `pm-md - out.md` receives JSON from stdin and saves markdown to out.md.

### custom templates

Easily customize the output by editing the template.

* `pm-md --get-template` creates a new file of the default template as an easier starting point for customization.
* `pm-md api.json --template=custom.tmpl` reads api.json and saves text formatted using the custom template file custom.tmpl to a new file.

In a template, you can use the variables and functions defined in [types.go](cmd/types.go). Sometimes it's helpful to look at the JSON exported from Postman to understand the variables. These template docs might also be helpful:

  * [How To Use Templates in Go — DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-use-templates-in-go#step-4-writing-a-template)
  * [the template package — Go's standard library](https://pkg.go.dev/text/template)

Here's how to add a table of contents with links that work in GitHub:

```
{{ range .Routes }}
* {{ formatHeaderLink .Name }}
{{- end }}
```

## tips

The result looks best when there is an example saved for each endpoint (after clicking "Send" in Postman, a "Save as Example" button appears).

Here's how to export a collection from Postman (choose the v2.1.0 export option):

![gif showing how to export from postman](https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExYzFnb2JicjN6czk3dTJqcjg5Zm1yMjVtOXZ4cGVzd2d6YjFuYm5tdyZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/1xp8s0yXApAYtF1c1q/giphy.gif)

## developer resources

* [intro to Go](https://wheelercj.github.io/notes/pages/20221122173910.html)
* [JSON and Go — The Go Blog](https://go.dev/blog/json)
* [the embed package — Go's standard library](https://pkg.go.dev/embed)
* [How To Use Struct Tags in Go — DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go)
* [spf13/cobra](https://github.com/spf13/cobra)
* [GoReleaser](https://goreleaser.com/)
* [how to create a custom terminal command](https://wheelercj.github.io/notes/pages/20220320181252.html)
* [Command Line Interface Guidelines](https://clig.dev/) by Prasad et al.
