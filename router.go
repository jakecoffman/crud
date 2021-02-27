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

// Add routes to the swagger spec and installs a handler with built-in validation. Some validation of the
// route itself occurs on Add so this is the kind of error that can be returned.
func (r *Router) Add(specs ...Spec) error {
	for i := range specs {
		spec := specs[i]

		if err := spec.Valid(); err != nil {
			return err
		}

		if _, ok := r.Swagger.Paths[spec.Path]; !ok {
			r.Swagger.Paths[spec.Path] = &Path{}
		}
		path := r.Swagger.Paths[spec.Path]
		var operation *Operation
		switch strings.ToLower(spec.Method) {
		case "get":
			if path.Get != nil {
				return fmt.Errorf("duplicate GET on route %v", spec.Path)
			}
			path.Get = &Operation{}
			operation = path.Get
		case "post":
			if path.Post != nil {
				return fmt.Errorf("duplicate POST on route %v", spec.Path)
			}
			path.Post = &Operation{}
			operation = path.Post
		case "put":
			if path.Put != nil {
				return fmt.Errorf("duplicate PUT on route %v", spec.Path)
			}
			path.Put = &Operation{}
			operation = path.Put
		case "patch":
			if path.Patch != nil {
				return fmt.Errorf("duplicate PATCH on route %v", spec.Path)
			}
			path.Patch = &Operation{}
			operation = path.Patch
		case "options":
			if path.Options != nil {
				return fmt.Errorf("duplicate PATCH on route %v", spec.Path)
			}
			path.Options = &Operation{}
			operation = path.Options
		case "delete":
			if path.Delete != nil {
				return fmt.Errorf("duplicate DELETE on route %v", spec.Path)
			}
			path.Delete = &Operation{}
			operation = path.Delete
		default:
			panic("Unhandled method " + spec.Method)
		}
		operation.Responses = DefaultResponse
		if spec.Responses != nil {
			operation.Responses = spec.Responses
		}
		operation.Tags = spec.Tags
		operation.Description = spec.Description
		operation.Summary = spec.Summary

		if spec.Validate.Path.Initialized() {
			params := ToSwaggerParameters(spec.Validate.Path, "path")
			operation.Parameters = append(operation.Parameters, params...)
		}
		if spec.Validate.Query.Initialized() {
			params := ToSwaggerParameters(spec.Validate.Query, "query")
			operation.Parameters = append(operation.Parameters, params...)
		}
		if spec.Validate.Header.Initialized() {
			params := ToSwaggerParameters(spec.Validate.Header, "header")
			operation.Parameters = append(operation.Parameters, params...)
		}
		if spec.Validate.FormData.Initialized() {
			params := ToSwaggerParameters(spec.Validate.Header, "formData")
			operation.Parameters = append(operation.Parameters, params...)
		}
		if spec.Validate.Body.Initialized() {
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

		// Finally add the route to gin. This is done last because gin will panic on duplicate entries.
		handlers := []gin.HandlerFunc{validationMiddleware(spec)}
		handlers = append(handlers, spec.PreHandlers...)
		handlers = append(handlers, spec.Handler)

		r.Mux.Handle(spec.Method, swaggerToGinPattern(spec.Path), handlers...)
	}
	return nil
}

// Validate are optional fields that will be used during validation. Leave unneeded
// properties nil and they will be ignored.
type Validate struct {
	Query    Field
	Body     Field
	Path     Field
	FormData Field
	Header   Field
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

func swaggerToGinPattern(swaggerUrl string) string {
	return swaggerPathPattern.ReplaceAllString(swaggerUrl, ":$1")
}

func pathParms(swaggerUrl string) (params []string) {
	for _, p := range swaggerPathPattern.FindAllString(swaggerUrl, -1) {
		params = append(params, p[1:len(p)-1])
	}
	return
}

//go:embed swaggerui.html
var swaggerUiTemplate []byte
