package crud

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Router struct {
	Swagger     string                `json:"swagger"`
	Info        Info                  `json:"info"`
	Paths       map[string]*Path      `json:"paths"`
	Definitions map[string]JsonSchema `json:"definitions"`

	Specs []Spec      `json:"-"`
	Mux   *gin.Engine `json:"-"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type JsonSchema struct {
	Type        string                `json:"type,omitempty"`
	Properties  map[string]JsonSchema `json:"properties,omitempty"`
	Required    []string              `json:"required,omitempty"`
	Example     interface{}           `json:"example,omitempty"`
	Description string                `json:"description,omitempty"`
	Minimum     float64               `json:"minimum,omitempty"`
	Maximum     float64               `json:"maximum,omitempty"`
	Enum        []interface{}         `json:"enum,omitempty"`
}

type Path struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
	Patch  *Operation `json:"patch,omitempty"`
}

type Operation struct {
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses"`
	Description string              `json:"description"`
	Summary     string              `json:"summary"`
}

type Parameter struct {
	In   string `json:"in"`
	Name string `json:"name"`

	Type   string `json:"type,omitempty"`
	Schema *Ref   `json:"schema,omitempty"`

	Required    *bool         `json:"required,omitempty"`
	Description string        `json:"description,omitempty"`
	Minimum     *float64      `json:"minimum,omitempty"`
	Maximum     *float64      `json:"maximum,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
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
		Info: Info{
			Title:   title,
			Version: version,
		},
		Mux: gin.Default(),
	}
}

func (r *Router) Add(specs ...Spec) {
	r.Specs = append(r.Specs, specs...)
}

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

type Validate struct {
	Query    map[string]Field
	Body     map[string]Field
	Path     map[string]Field
	FormData map[string]Field
	Header   map[string]Field
}

func (r *Router) Use(middlewares ...gin.HandlerFunc) {
	r.Mux.Use(middlewares...)
}

func (r *Router) Serve(addr string) error {
	modelCounter := 1
	r.Definitions = map[string]JsonSchema{}

	r.Paths = map[string]*Path{}
	for i := range r.Specs {
		spec := r.Specs[i]

		handlers := []gin.HandlerFunc{preHandler(spec)}
		handlers = append(handlers, spec.PreHandlers...)
		handlers = append(handlers, spec.Handler)

		r.Mux.Handle(spec.Method, swaggerToGinPattern(spec.Path), handlers...)

		if _, ok := r.Paths[spec.Path]; !ok {
			r.Paths[spec.Path] = &Path{}
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
		case "patch":
			path.Patch = &Operation{Responses: DefaultResponse}
			operation = path.Patch
		case "delete":
			path.Delete = &Operation{Responses: DefaultResponse}
			operation = path.Delete
		default:
			panic("Unhandled method " + spec.Method)
		}
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
			modelName := fmt.Sprintf("Model %v", modelCounter)
			parameter := Parameter{
				In:     "body",
				Name:   "body",
				Schema: &Ref{fmt.Sprint("#/definitions/", modelName)},
			}
			r.Definitions[modelName] = ToJsonSchema(spec.Validate.Body)
			modelCounter++
			operation.Parameters = append(operation.Parameters, parameter)
		}
	}

	r.Mux.GET("/swagger.json", func(c *gin.Context) {
		c.JSON(200, r)
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

func preHandler(spec Spec) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := spec.Validate
		if val.Query != nil {
			values := c.Request.URL.Query()
			for field, schema := range val.Query {
				// query values are always strings, so we must try to convert
				queryValue := values.Get(field)

				// don't try to convert if the field is empty
				if queryValue == "" {
					if schema.required != nil && *schema.required {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, ErrRequired))
					}
					return
				}
				var convertedValue interface{}
				switch schema.kind {
				case KindBoolean:
					if queryValue == "true" {
						convertedValue = true
					} else if queryValue == "false" {
						convertedValue = false
					} else {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, ErrWrongType))
						return
					}
				case KindString:
					convertedValue = queryValue
				case KindNumber:
					var err error
					convertedValue, err = strconv.ParseFloat(queryValue, 64)
					if err != nil {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, ErrWrongType))
						return
					}
				case KindInteger:
					var err error
					convertedValue, err = strconv.Atoi(queryValue)
					if err != nil {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, ErrWrongType))
						return
					}
				case KindArray:
					// TODO I'm not sure how this works yet
					c.AbortWithStatusJSON(http.StatusNotImplemented, "TODO")
					return
				default:
					c.AbortWithStatusJSON(400, fmt.Sprintf("Validation not possible due to kind: %v", schema.kind))
				}
				if err := schema.Validate(convertedValue); err != nil {
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
				if schema.required != nil && *schema.required && path == "" {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Missing path param"))
					return
				}
			}
		}
	}
}

var swaggerPathPattern = regexp.MustCompile("\\{([^}]+)\\}")

func swaggerToGinPattern(ginUrl string) string {
	return swaggerPathPattern.ReplaceAllString(ginUrl, ":$1")
}

//go:embed swaggerui.html
var swaggerUiTemplate []byte
