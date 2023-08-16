# Postman to Markdown

Convert a Postman collection to markdown documentation.

[sample result](samples/calendar%20API%20v1.md)

The result looks best when there is an example saved for each endpoint (after clicking "Send" in Postman, a "Save as Example" button appears).

## download

3 choices for how to download:

* [download a zipped executable file](https://github.com/wheelercj/pm-md/releases), unzip it, and run the app with `./pm-md "json file path here"`
* `go install github.com/wheelercj/pm-md@latest` and then `pm-md "json file path here"`
* install from source following the instructions below

You can get the JSON file from Postman by exporting a collection as a v2.1.0 collection.

### install from source

These steps require [Go](https://go.dev/) to be installed.

1. In Postman, export a collection as a v2.1.0 collection.
2. Choose one of the source code download options [here](https://github.com/wheelercj/pm-md/releases) and unzip the folder (or use `git clone`).
3. Open a terminal in the new folder.
4. Run `go build` to create an executable file.
5. Run `go install` to install the executable file. If you get an error message, you may need to [edit your PATH environment variable](https://go.dev/doc/tutorial/compile-install).
6. Run the app with `pm-md "json file path here"`. If you installed the executable file, this should work in any directory.

If you install from source, the resulting markdown file's format can be customized by editing `collection.tmpl` using the types defined in `types.go` and the `template.FuncMap` defined in `main.go`. See the links under "developer resources" below for more details about templates. Use `go build` and `go install` after editing.

## developer resources

Here are some resources that were helpful when creating this app.

* [intro to Go](https://wheelercj.github.io/notes/pages/20221122173910.html)
* [the template package — Go's standard library](https://pkg.go.dev/text/template)
* [How To Use Templates in Go — DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-use-templates-in-go)
* [JSON and Go — The Go Blog](https://go.dev/blog/json)
* [the embed package — Go's standard library](https://pkg.go.dev/embed)
* [How To Use Struct Tags in Go — DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go)
* [GoReleaser](https://goreleaser.com/)
* [how to create a custom terminal command](https://wheelercj.github.io/notes/pages/20220320181252.html)
