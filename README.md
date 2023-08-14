# Postman to Markdown

Convert a Postman collection to markdown documentation.

* [sample result](samples/calendar%20API%20v1.md)
* [sample source JSON](samples/calendar%20API.postman_collection.json)

## usage

1. In Postman, export the collection as a v2.1.0 collection.
2. In a terminal, run `git clone https://github.com/wheelercj/postman-to-markdown.git` where you want this program's folder to appear and `cd` into the new folder.
3. Use `go run main.go "file path here"`.

If you will use this often, you may want to [create a custom terminal command](https://wheelercj.github.io/notes/pages/20220320181252.html).

## developer resources

Here are some resources that were helpful when creating this program.

* [JSON and Go -- The Go Blog](https://go.dev/blog/json)
* [the template package -- Go's standard library](https://pkg.go.dev/text/template)
