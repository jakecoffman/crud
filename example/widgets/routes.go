package widgets

import (
	"github.com/jakecoffman/crud"
)

var tags = []string{"Widgets"}

var Routes = []crud.Spec{{
	Method:      "GET",
	Path:        "/widgets",
	Handler:     ListHandler,
	Description: "Lists widgets",
	Tags:        tags,
	Validate: crud.Validate{
		Query: map[string]crud.Field{
			"limit": crud.Number().Required().Min(0).Max(25).Description("Records to return"),
		},
	},
}, {
	Method:      "POST",
	Path:        "/widgets",
	Handler:     CreateHandler,
	Description: "Adds a widget",
	Tags:        tags,
	Validate: crud.Validate{
		Body: map[string]crud.Field{
			"name": crud.String().Required().Example("Bob"),
		},
	},
}, {
	Method:      "GET",
	Path:        "/widgets/{id}",
	Handler:     GetHandler,
	Description: "Updates a widget",
	Tags:        tags,
	Validate: crud.Validate{
		Path: map[string]crud.Field{
			"id": crud.Number().Required(),
		},
	},
}, {
	Method:      "PUT",
	Path:        "/widgets/{id}",
	Handler:     UpdateHandler,
	Description: "Updates a widget",
	Tags:        tags,
	Validate: crud.Validate{
		Path: map[string]crud.Field{
			"id": crud.Number().Required(),
		},
		Body: map[string]crud.Field{
			"name": crud.String().Required(),
		},
	},
}, {
	Method:      "DELETE",
	Path:        "/widgets/{id}",
	Handler:     DeleteHandler,
	Description: "Deletes a widget",
	Tags:        tags,
	Validate: crud.Validate{
		Path: map[string]crud.Field{
			"id": crud.Number().Required(),
		},
	},
},
}
