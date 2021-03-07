package main

import (
	"encoding/json"
	"github.com/jakecoffman/crud"
	"github.com/jakecoffman/crud/adapters/echo-adapter"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	r := crud.NewRouter("Widget API", "1.0.0", adapter.New())

	if err := r.Add(Routes...); err != nil {
		log.Fatal(err)
	}

	log.Println("Serving http://127.0.0.1:8080")
	err := r.Serve("127.0.0.1:8080")
	if err != nil {
		log.Println(err)
	}
}

var tags = []string{"Widgets"}

var Routes = []crud.Spec{{
	Method: "GET",
	Path:   "/widgets",
	Handler: func(c echo.Context) error {
		return c.JSON(200, c.Request().URL.Query())
	},
	Description: "Lists widgets",
	Tags:        tags,
	Validate: crud.Validate{
		Query: crud.Object(map[string]crud.Field{
			"limit": crud.Number().Required().Min(0).Max(25).Description("Records to return"),
			"ids":   crud.Array().Items(crud.Number()),
		}),
	},
}, {
	Method: "POST",
	Path:   "/widgets",
	Handler: func(c echo.Context) error {
		var widget interface{}
		if err := c.Bind(&widget); err != nil {
			return err
		}
		return c.JSON(200, widget)
	},
	Description: "Adds a widget",
	Tags:        tags,
	Validate: crud.Validate{
		Body: crud.Object(map[string]crud.Field{
			"name":       crud.String().Required().Example("Bob"),
			"arrayMatey": crud.Array().Items(crud.Number()),
		}),
	},
	Responses: map[string]crud.Response{
		"200": {
			Schema: crud.JsonSchema{
				Type: crud.KindObject,
				Properties: map[string]crud.JsonSchema{
					"hello": {Type: crud.KindString},
				},
			},
			Description: "OK",
		},
	},
}, {
	Method: "GET",
	Path:   "/widgets/{id}",
	Handler: func(c echo.Context) error {
		return c.JSON(200, c.Param("id"))
	},
	Description: "Updates a widget",
	Tags:        tags,
	Validate: crud.Validate{
		Path: crud.Object(map[string]crud.Field{
			"id": crud.Number().Required(),
		}),
	},
}, {
	Method: "PUT",
	Path:   "/widgets/{id}",
	Handler: func(c echo.Context) error {
		// for some reason I can't bind to interface{} here but I can above
		var widget struct {
			json.RawMessage
		}
		if err := c.Bind(&widget); err != nil {
			return err
		}
		return c.JSON(200, widget)
	},
	Description: "Updates a widget",
	Tags:        tags,
	Validate: crud.Validate{
		Path: crud.Object(map[string]crud.Field{
			"id": crud.Number().Required(),
		}),
		Body: crud.Object(map[string]crud.Field{
			"name": crud.String().Required(),
		}),
	},
}, {
	Method: "DELETE",
	Path:   "/widgets/{id}",
	Handler: func(c echo.Context) error {
		return c.JSON(200, c.Param("id"))
	},
	Description: "Deletes a widget",
	Tags:        tags,
	Validate: crud.Validate{
		Path: crud.Object(map[string]crud.Field{
			"id": crud.Number().Required(),
		}),
	},
},
}
