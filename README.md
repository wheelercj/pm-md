# Postman to Markdown

Convert a Postman collection to markdown documentation.

[sample result](samples/calendar%20API%20v1.md)

If you install from source, the resulting markdown file's format is easy to customize by editing `collection.tmpl` using the types defined in `types.go` and the `template.FuncMap` defined in `main.go`. See the links under "developer resources" below for more details about templates.

The results look best when there is an example saved for each endpoint. After clicking "Send" in Postman, you can click "Save as Example" to save an example. Examples are inserted into the markdown file without any escaping, so **only use this app on JSON files that you trust**.

## install

### Windows

1. [Click here to download](https://github.com/wheelercj/postman-to-markdown/releases/download/v0.0.1/pm-md.zip).
2. Unzip the file.
3. In Postman, export a collection as a v2.1.0 collection.
4. Run the app in a terminal with `pm-md "json file path here"`. The terminal's working directory must be where pm-md.exe is.

If you will use this app often, you might want to [create a custom terminal command](https://wheelercj.github.io/notes/pages/20220320181252.html).

### Mac, Linux, and Windows (install from source)

These steps require [Git](https://git-scm.com/) and [Go](https://go.dev/) to be installed.

1. In Postman, export a collection as a v2.1.0 collection.
2. In a terminal, run `git clone https://github.com/wheelercj/postman-to-markdown.git` where you want this app's folder to appear.
3. `cd` into the new folder.
4. Use `go build` to create an executable file.
5. Run the app with `./postman-to-markdown "json file path here"`. You can rename the executable file from `postman-to-markdown` to something shorter to make this easier if you want to.

If you will use this app often, you might want to [install the app](https://go.dev/doc/tutorial/compile-install).

## developer resources

Here are some resources that were helpful when creating this app.

* [JSON and Go -- The Go Blog](https://go.dev/blog/json)
* [the template package -- Go's standard library](https://pkg.go.dev/text/template)
* [How To Use Templates in Go -- DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-use-templates-in-go)
