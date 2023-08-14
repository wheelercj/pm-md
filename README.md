# Postman to Markdown

Convert a Postman collection to markdown documentation.

* [sample result](samples/calendar%20API%20v1.md)
* [sample source JSON](samples/calendar%20API.postman_collection.json)

## install

### Windows

1. [Click here to download](https://github.com/wheelercj/postman-to-markdown/releases/download/v0.0.1/pm-md.zip).
2. Unzip the file.
3. In Postman, export a collection as a v2.1.0 collection.
4. Run the app in a terminal with `pm-md.exe "json file path here"`.

If you will use this app often, you might want to [create a custom terminal command](https://wheelercj.github.io/notes/pages/20220320181252.html).

### Mac, Linux, and Windows

1. In a terminal, run `git clone https://github.com/wheelercj/postman-to-markdown.git` where you want this app's folder to appear and `cd` into the new folder.
2. In Postman, export a collection as a v2.1.0 collection.
3. In the terminal, use `go run main.go "json file path here"`.

If you will use this app often, you might want to [install the app](https://go.dev/doc/tutorial/compile-install).

## developer resources

Here are some resources that were helpful when creating this app.

* [JSON and Go -- The Go Blog](https://go.dev/blog/json)
* [the template package -- Go's standard library](https://pkg.go.dev/text/template)
