## crud

[![GoDoc](https://godoc.org/github.com/jakecoffman/crud?status.svg)](https://godoc.org/github.com/jakecoffman/crud)
[![Go](https://github.com/jakecoffman/crud/actions/workflows/go.yml/badge.svg)](https://github.com/jakecoffman/crud/actions/workflows/go.yml)

A Swagger/OpenAPI builder and validation library for servers.

Heavily inspired by [hapi](https://hapi.dev/) and the [hapi-swagger](https://github.com/glennjones/hapi-swagger) projects.

### Status

This project is not stable yet, API is still changing and there are missing features.

It's already pretty useful so if you are not risk averse then use it and pitch in!

### Why

Swagger is great, but up until now your options to use swagger are:

- Write it and then make your server match your spec.
- Write it and generate your server.
- Generate it from comments in your code.

None of these options seems like a great idea.

This project takes another approach: make a specification in Go code using nice builders where possible. The swagger is generated from this spec and validation is done before your handler gets called. 

This reduces boilerplate that you have to write and gives you nice documentation too!

### Getting started

Check the example directory for a simple example.

Start by getting the package `go get github.com/jakecoffman/crud`

Then in your `main.go`:

1. Create a router with `NewRouter`
2. Add routes with `Add`
3. Then call `Serve`

Routes are specifications that look like this:

```go
crud.Spec{
	Method:      "PATCH",
	Path:        "/widgets/{id}",
	PreHandlers: []gin.HandlerFunc{Auth},
	Handler:     CreateHandler,
	Description: "Adds a widget",
	Tags:        []string{"Widgets"},
	Validate: crud.Validate{
		Path: crud.Object(map[string]crud.Field{
			"id": crud.Number().Required().Description("ID of the widget")
        }),
		Body: crud.Object(map[string]crud.Field{
			"owner": crud.String().Required().Example("Bob").Description("Widget owner's name"),
		}),
	},
}
```

This will add a route `/widgets/:id` that responds to the PATCH method. It generates swagger and serves it at the root of the web application. It validates that the ID in the path is a number, so you don't have to. It also validates that the body is an object and has an "owner" property that is a string, again so you won't have to.

The `PreHandlers` run before validation, and the `Handler` runs after validation is successful.
