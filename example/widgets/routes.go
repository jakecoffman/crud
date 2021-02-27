package widgets

import (
	"github.com/gin-gonic/gin"
	"github.com/jakecoffman/crud"
)

var tags = []string{"Widgets"}

var Routes = []crud.Spec{{
	Method: "GET",
	Path:   "/widgets",
	Handler: func(c *gin.Context) {
		c.JSON(200, c.Request.URL.Query())
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
	Handler: func(c *gin.Context) {
		var widget interface{}
		if err := c.BindJSON(&widget); err != nil {
			return
		}
		c.JSON(200, widget)
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
	Handler: func(c *gin.Context) {
		c.JSON(200, c.Params)
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
	Handler: func(c *gin.Context) {
		var widget interface{}
		if err := c.BindJSON(&widget); err != nil {
			return
		}
		c.JSON(200, widget)
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
	Handler: func(c *gin.Context) {
		c.JSON(200, c.Params)
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
