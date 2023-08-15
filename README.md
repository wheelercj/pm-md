# Postman to Markdown

Convert a Postman collection to markdown documentation.

[sample result](samples/calendar%20API%20v1.md)

If you install from source, the resulting markdown file's format is easy to customize by editing `collection.tmpl` using the types defined in `types.go` and the `template.FuncMap` defined in `main.go`. See the links under "developer resources" below for more details about templates.

The results look best when there is an example saved for each endpoint. After clicking "Send" in Postman, you can click "Save as Example" to save an example.

## install

### Windows

1. [Click here to download](https://github.com/wheelercj/pm-md/releases/download/v0.0.2/pm-md.zip).
2. Unzip the file.
3. In Postman, export a collection as a v2.1.0 collection.
4. Run the app in a terminal with `pm-md "json file path here"`. The terminal's working directory must be where pm-md.exe is.

If you will use this app often, you might want to [create a custom terminal command](https://wheelercj.github.io/notes/pages/20220320181252.html).

### Mac, Linux, and Windows (install from source)

These steps require [Go](https://go.dev/) to be installed.

1. In Postman, export a collection as a v2.1.0 collection.
2. Choose one of the source code download options [here](https://github.com/wheelercj/pm-md/releases) (or `git clone`).
3. Unzip the folder.
4. Open a terminal in the new folder.
5. Run `go build` to create an executable file.
6. Run `go install` to install the executable file. If you get an error message, you may need to [edit your PATH environment variable](https://go.dev/doc/tutorial/compile-install).
7. Run the app with `pm-md "json file path here"`. If you installed the executable file, this should work in any directory.

## developer resources

Here are some resources that were helpful when creating this app.

* [the template package -- Go's standard library](https://pkg.go.dev/text/template)
* [How To Use Templates in Go -- DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-use-templates-in-go)
* [JSON and Go -- The Go Blog](https://go.dev/blog/json)
* [the embed package -- Go's standard library](https://pkg.go.dev/embed)
