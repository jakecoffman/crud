package crud

import (
	_ "embed"
	"fmt"
	"github.com/jakecoffman/crud/option"
	"regexp"
	"strings"
)

// Router is the main object that is used to generate swagger and holds the underlying router.
type Router struct {
	// Swagger is exposed so the user can edit additional optional fields.
	Swagger Swagger

	// The underlying router being used behind Adapter interface.
	adapter Adapter

	// used for automatically incrementing the model name, e.g. Model 1, Model 2.
	modelCounter int

	// options
	stripUnknown bool
	allowUnknown bool
}

type Adapter interface {
	Install(router *Router, spec *Spec) error
	Serve(swagger *Swagger, addr string) error
}

// NewRouter initializes a router.
func NewRouter(title, version string, adapter Adapter, options ...option.Option) *Router {
	r := &Router{
		Swagger: Swagger{
			OpenAPI: "3.0.3",
			Info:    Info{Title: title, Version: version},
			Paths:   map[string]*Path{},
			Components: Components{
				Schemas:         map[string]*JsonSchema{},
				Responses:       map[string]*JsonSchema{},
				Parameters:      map[string]*JsonSchema{},
				Examples:        map[string]*JsonSchema{},
				RequestBodies:   map[string]*JsonSchema{},
				Headers:         map[string]*JsonSchema{},
				SecuritySchemas: map[string]*JsonSchema{},
				Links:           map[string]*JsonSchema{},
				Callbacks:       map[string]*JsonSchema{},
			},
		},
		adapter:      adapter,
		modelCounter: 1,
		stripUnknown: true,
		allowUnknown: true,
	}
	for _, o := range options {
		if o.StripUnknown != nil {
			r.stripUnknown = *o.StripUnknown
		} else if o.AllowUnknown != nil {
			r.allowUnknown = *o.AllowUnknown
		}
	}
	return r
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
		operation.Responses = defaultResponse
		if spec.Responses != nil {
			operation.Responses = spec.Responses
		}
		operation.Tags = spec.Tags
		operation.Description = spec.Description
		operation.Summary = spec.Summary

		if spec.Validate.Path.Initialized() {
			params := spec.Validate.Path.ToSwaggerParameters("path")
			operation.Parameters = append(operation.Parameters, params...)
		}
		if spec.Validate.Query.Initialized() {
			params := spec.Validate.Query.ToSwaggerParameters("query")
			operation.Parameters = append(operation.Parameters, params...)
		}
		if spec.Validate.Header.Initialized() {
			params := spec.Validate.Header.ToSwaggerParameters("header")
			operation.Parameters = append(operation.Parameters, params...)
		}
		if spec.Validate.FormData.Initialized() {
			params := spec.Validate.FormData.ToSwaggerParameters("formData")
			operation.Parameters = append(operation.Parameters, params...)
		}
		if spec.Validate.Body.Initialized() {
			modelName := fmt.Sprintf("Model %v", r.modelCounter)
			r.Swagger.Components.Schemas[modelName] = spec.Validate.Body.ToJsonSchema()
			r.modelCounter++
			operation.RequestBody = &RequestBody{Content: map[string]MediaType{
				"*/*": {Schema: &Ref{fmt.Sprint("#/components/schemas/", modelName)}},
			}}
		}

		if err := r.adapter.Install(r, &spec); err != nil {
			return err
		}
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
	return r.adapter.Serve(&r.Swagger, addr)
}

// SwaggerPathPattern regex captures swagger path params.
var SwaggerPathPattern = regexp.MustCompile("\\{([^}]+)\\}")

func pathParms(swaggerUrl string) (params []string) {
	for _, p := range SwaggerPathPattern.FindAllString(swaggerUrl, -1) {
		params = append(params, p[1:len(p)-1])
	}
	return
}

// SwaggerUiTemplate contains the html for swagger UI.
//go:embed swaggerui.html
var SwaggerUiTemplate []byte
