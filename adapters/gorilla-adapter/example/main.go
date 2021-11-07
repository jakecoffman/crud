package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/crud"
	"github.com/jakecoffman/crud/adapters/gorilla-adapter"
	"log"
	"math/rand"
	"net/http"
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
	Handler: func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(r.URL.Query())
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
	PreHandlers: func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rand.Intn(2) == 0 {
				w.WriteHeader(http.StatusTeapot)
				_, _ = w.Write([]byte("Random rejection from PreHandler"))
				return
			}
			next.ServeHTTP(w, r)
		})
	},
	Handler: func(w http.ResponseWriter, r *http.Request) {
		var widget interface{}
		if err := json.NewDecoder(r.Body).Decode(&widget); err != nil {
			w.WriteHeader(400)
			_ = json.NewEncoder(w).Encode("Failed reading body: " + err.Error())
			return
		}
		_ = json.NewEncoder(w).Encode(widget)
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
	Handler: func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(mux.Vars(r))
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
	Handler: func(w http.ResponseWriter, r *http.Request) {
		var widget interface{}
		if err := json.NewDecoder(r.Body).Decode(&widget); err != nil {
			w.WriteHeader(400)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
		_ = json.NewEncoder(w).Encode(widget)
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
	Handler: func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(mux.Vars(r))
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
