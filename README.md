## crud

[![GoDoc](https://godoc.org/github.com/jakecoffman/crud?status.svg)](https://godoc.org/github.com/jakecoffman/crud)
[![Go](https://github.com/jakecoffman/crud/actions/workflows/go.yml/badge.svg)](https://github.com/jakecoffman/crud/actions/workflows/go.yml)

A Swagger/OpenAPI builder and validation library for building HTTP/REST APIs.

Heavily inspired by [hapi](https://hapi.dev/) and the [hapi-swagger](https://github.com/glennjones/hapi-swagger) projects.

No additional dependencies besides the router you choose.

### Status

Version 1.0 is stable, version 2 will support OpenAPI v3.

### Why

Swagger is great, but up until now your options to use swagger are:

1. Write it by hand and then make your server match your spec.
2. Write it by hand and generate your server.
3. Generate it from comments in your code.

None of these options seems like a great idea.

This project takes another approach: make a specification in Go code using nice builders where possible. The swagger is generated from this spec and validation is done before your handler gets called. 

This reduces boilerplate that you have to write and gives you nice documentation too!

### Examples

- [Full Gin-Gonic Example](adapters/gin-adapter/example)
- [Full Echo Example](adapters/echo-adapter/example)
- [Full Gorilla Mux Example](adapters/gorilla-adapter/example)

### Builtin ServeMux not supported

This is disappointing, but the builtin http.ServeMux is not supported because it doesn't support routing by method, and doesn't support path params. This project is NOT a router so it will not try to reinvent these features.

### Getting started

Check the example directory under the adapters for a simple example.

Start by getting the package `go get github.com/jakecoffman/crud`

Then in your `main.go`:

1. Create a router with `NewRouter`, use an adapter from the `adapters` sub-package or write you own.
2. Add routes with `Add`.
3. Then call `Serve`.

Routes are specifications that look like this:

```go
crud.Spec{
	Method:      "PATCH",
	Path:        "/widgets/{id}",
	PreHandlers: Auth,
	Handler:     CreateHandler,
	Description: "Adds a widget",
	Tags:        []string{"Widgets"},
	Validate: crud.Validate{
		Path: crud.Object(map[string]crud.Field{
			"id": crud.Number().Required().Description("ID of the widget"),
        }),
		Body: crud.Object(map[string]crud.Field{
			"owner": crud.String().Required().Example("Bob").Description("Widget owner's name"),
			"quantity": crud.Integer().Min(1).Default(1).Description("The amount requested")
		}),
	},
}
```

This will add a route `/widgets/:id` that responds to the PATCH method. It generates swagger and serves it at the root of the web application. It validates that the ID in the path is a number, so you don't have to. It also validates that the body is an object and has an "owner" property that is a string, again so you won't have to.

It mounts the swagger-ui at `/` and loads up the generated swagger.json:

![screenshot](/screenshot.png?raw=true "Swagger")

The `PreHandlers` run before validation, and the `Handler` runs after validation is successful.
