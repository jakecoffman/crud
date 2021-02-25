package crud

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

// Router is the main object that is used to generate swagger and holds the underlying router.
type Router struct {
	// Swagger is exposed so the user can edit additional optional fields.
	Swagger Swagger

	// Mux is the underlying router being used. The user can add middlewares and use other features.
	Mux *gin.Engine

	// used for automatically incrementing the model name, e.g. Model 1, Model 2.
	modelCounter int
}

// NewRouter initializes a router.
func NewRouter(title, version string) *Router {
	return &Router{
		Swagger: Swagger{
			Swagger:     "2.0",
			Info:        Info{Title: title, Version: version},
			Paths:       map[string]*Path{},
			Definitions: map[string]JsonSchema{},
		},
		Mux:          gin.Default(),
		modelCounter: 1,
	}
}

// Add routes to the swagger spec and installs a handler with built-in validation.
func (r *Router) Add(specs ...Spec) {
	for i := range specs {
		spec := specs[i]

		handlers := []gin.HandlerFunc{validationMiddleware(spec)}
		handlers = append(handlers, spec.PreHandlers...)
		handlers = append(handlers, spec.Handler)

		r.Mux.Handle(spec.Method, swaggerToGinPattern(spec.Path), handlers...)

		if _, ok := r.Swagger.Paths[spec.Path]; !ok {
			r.Swagger.Paths[spec.Path] = &Path{}
		}
		path := r.Swagger.Paths[spec.Path]
		var operation *Operation
		switch strings.ToLower(spec.Method) {
		case "get":
			path.Get = &Operation{}
			operation = path.Get
		case "post":
			path.Post = &Operation{}
			operation = path.Post
		case "put":
			path.Put = &Operation{}
			operation = path.Put
		case "patch":
			path.Patch = &Operation{}
			operation = path.Patch
		case "delete":
			path.Delete = &Operation{}
			operation = path.Delete
		default:
			panic("Unhandled method " + spec.Method)
		}
		operation.Responses = DefaultResponse
		operation.Tags = spec.Tags
		operation.Description = spec.Description
		operation.Summary = spec.Summary

		if spec.Validate.Path != nil {
			for name, field := range spec.Validate.Path {
				param := Parameter{
					In:          "path",
					Name:        name,
					Type:        field.kind,
					Required:    field.required,
					Description: field.description,
					Enum:        field.enum,
					Minimum:     field.min,
					Maximum:     field.max,
				}
				operation.Parameters = append(operation.Parameters, param)
			}
		}
		if spec.Validate.Query != nil {
			for name, field := range spec.Validate.Query {
				param := Parameter{
					In:          "query",
					Name:        name,
					Type:        field.kind,
					Required:    field.required,
					Description: field.description,
					Enum:        field.enum,
					Minimum:     field.min,
					Maximum:     field.max,
				}
				operation.Parameters = append(operation.Parameters, param)
			}
		}
		if spec.Validate.Body != nil {
			modelName := fmt.Sprintf("Model %v", r.modelCounter)
			parameter := Parameter{
				In:     "body",
				Name:   "body",
				Schema: &Ref{fmt.Sprint("#/definitions/", modelName)},
			}
			r.Swagger.Definitions[modelName] = ToJsonSchema(spec.Validate.Body)
			r.modelCounter++
			operation.Parameters = append(operation.Parameters, parameter)
		}
	}
}

// Spec is used to generate swagger paths and automatic handler validation.
type Spec struct {
	Method      string
	Path        string
	PreHandlers []gin.HandlerFunc
	Handler     gin.HandlerFunc
	Description string
	Tags        []string
	Summary     string

	Validate Validate
}

// Validate are optional fields that will be used during validation. Leave unneeded
// properties nil and they will be ignored.
type Validate struct {
	Query    map[string]Field
	Body     map[string]Field
	Path     map[string]Field
	FormData map[string]Field
	Header   map[string]Field
}

// Serve installs the swagger and the swagger-ui and runs the server.
func (r *Router) Serve(addr string) error {
	r.Mux.GET("/swagger.json", func(c *gin.Context) {
		c.JSON(200, r.Swagger)
	})

	r.Mux.GET("/", func(c *gin.Context) {
		c.Header("content-type", "text/html")
		_, err := c.Writer.Write(swaggerUiTemplate)
		if err != nil {
			panic(err)
		}
	})

	err := r.Mux.Run(addr)
	return err
}

// we need to convert swagger endpoints /widget/{id} to gin endpoints /widget/:id
var swaggerPathPattern = regexp.MustCompile("\\{([^}]+)\\}")

func swaggerToGinPattern(ginUrl string) string {
	return swaggerPathPattern.ReplaceAllString(ginUrl, ":$1")
}

//go:embed swaggerui.html
var swaggerUiTemplate []byte
