package crud

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

type Router struct {
	Swagger     string                `json:"swagger"`
	Title       string                `json:"title"`
	Version     string                `json:"version"`
	Paths       map[string]*Path      `json:"paths"`
	Definitions map[string]JsonSchema `json:"definitions"`

	Specs []Spec `json:"-"`
}

type JsonSchema struct {
	Type       string                `json:"type,omitempty"`
	Properties map[string]JsonSchema `json:"properties,omitempty"`
	Required   []string              `json:"required,omitempty"`
	Example    interface{}           `json:"example,omitempty"`
}

type Path struct {
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	Get         *Operation `json:"get,omitempty"`
	Post        *Operation `json:"post,omitempty"`
	Put         *Operation `json:"put,omitempty"`
	Delete      *Operation `json:"delete,omitempty"`
}

type Operation struct {
	Tags       []string            `json:"tags,omitempty"`
	Parameters []Parameter         `json:"parameters,omitempty"`
	Responses  map[string]Response `json:"responses"`
}

type Parameter struct {
	In   string `json:"in"`
	Name string `json:"name"`

	Type   string `json:"type,omitempty"`
	Schema Ref    `json:"schema,omitempty"`

	Required *bool `json:"required,omitempty"`
}

type Ref struct {
	Ref string `json:"$ref,omitempty"`
}

type Response struct {
	Schema      JsonSchema `json:"schema"`
	Description string     `json:"description"`
}

var DefaultResponse = map[string]Response{
	"default": {
		Schema:      JsonSchema{Type: "string"},
		Description: "Successful",
	},
}

func NewRouter(title, version string) *Router {
	return &Router{
		Swagger: "2.0",
		Title:   title,
		Version: version,
	}
}

func (r *Router) Add(specs ...Spec) {
	r.Specs = append(r.Specs, specs...)
}

type Spec struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Options     Options
	Summary     string
	Description string
}

type Options struct {
	Description string
	Tags        []string
	Validate    Validate
}

type Validate struct {
	Query    map[string]Field
	Body     map[string]Field
	Path     map[string]Field
	FormData map[string]Field
	Header   map[string]Field
}

func (r *Router) Serve(addr string, middlewares ...gin.HandlerFunc) error {
	// Build the routes and serve
	mux := gin.Default()
	for _, m := range middlewares {
		mux.Use(m)
	}

	modelCounter := 1
	r.Definitions = map[string]JsonSchema{}

	r.Paths = map[string]*Path{}
	for i := range r.Specs {
		spec := r.Specs[i]
		mux.Handle(spec.Method, spec.Path, func(c *gin.Context) {
			val := spec.Options.Validate
			if val.Query != nil {
				values := c.Request.URL.Query()
				for field, schema := range val.Query {
					if err := schema.Validate(values.Get(field)); err != nil {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, err.Error()))
						return
					}
				}
			}

			if val.Body != nil {
				var body map[string]interface{}
				if err := c.BindJSON(&body); err != nil {
					c.AbortWithStatusJSON(400, err.Error())
					return
				}
				for field, schema := range val.Body {
					if err := schema.Validate(body[field]); err != nil {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed for field %v: %v", field, err.Error()))
						return
					}
				}
				// TODO perhaps the user passes a struct to bind to instead?
				data, _ := json.Marshal(body)
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
			}

			if val.Path != nil {
				for field, schema := range val.Path {
					path := c.Param(field)
					if schema.Req != nil && *schema.Req && path == "" {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Missing path param"))
						return
					}
				}
			}

			spec.Handler(c)
		})

		if _, ok := r.Paths[spec.Path]; !ok {
			r.Paths[spec.Path] = &Path{
				Summary:     spec.Summary,
				Description: spec.Description,
			}
		}
		path := r.Paths[spec.Path]
		var operation *Operation
		switch strings.ToLower(spec.Method) {
		case "get":
			path.Get = &Operation{Responses: DefaultResponse}
			operation = path.Get
		case "post":
			path.Post = &Operation{Responses: DefaultResponse}
			operation = path.Post
		case "put":
			path.Put = &Operation{Responses: DefaultResponse}
			operation = path.Put
		case "delete":
			path.Delete = &Operation{Responses: DefaultResponse}
			operation = path.Delete
		default:
			panic("Unhandled method " + spec.Method)
		}
		operation.Tags = spec.Options.Tags

		if spec.Options.Validate.Path != nil {
			for name, field := range spec.Options.Validate.Path {
				operation.Parameters = append(operation.Parameters, Parameter{
					In:       "path",
					Name:     name,
					Type:     field.Type,
					Required: field.Req,
				})
			}
		}
		if spec.Options.Validate.Query != nil {
			for name, field := range spec.Options.Validate.Query {
				operation.Parameters = append(operation.Parameters, Parameter{
					In:       "query",
					Name:     name,
					Type:     field.Type,
					Required: field.Req,
				})
			}
		}
		if spec.Options.Validate.Body != nil {
			modelName := fmt.Sprintf("Model %v", modelCounter)
			parameter := Parameter{
				In:     "body",
				Name:   "body",
				Schema: Ref{fmt.Sprint("#/definitions/", modelName)},
			}
			r.Definitions[modelName] = ToJsonSchema(spec.Options.Validate.Body)
			modelCounter++
			operation.Parameters = append(operation.Parameters, parameter)
		}
	}

	mux.GET("/swagger.json", func(c *gin.Context) {
		c.JSON(200, r)
	})

	mux.GET("/", func(c *gin.Context) {
		c.Header("content-type", "text/html")
		_, err := c.Writer.Write(swaggerUiTemplate)
		if err != nil {
			panic(err)
		}
	})

	err := mux.Run(addr)
	return err
}

//go:embed swaggerui.html
var swaggerUiTemplate []byte
