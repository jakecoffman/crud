package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jakecoffman/crud"
	"io/ioutil"
	"net/url"
	"reflect"
)

type Adapter struct {
	Engine *gin.Engine
}

func New() *Adapter {
	return &Adapter{
		Engine: gin.Default(),
	}
}

func (a *Adapter) Install(r *crud.Router, spec *crud.Spec) error {
	handlers := []gin.HandlerFunc{wrap(r, spec)}

	switch v := spec.PreHandlers.(type) {
	case nil:
		// ok
	case []gin.HandlerFunc:
		handlers = append(handlers, v...)
	case gin.HandlerFunc:
		handlers = append(handlers, v)
	case func(*gin.Context):
		handlers = append(handlers, v)
	default:
		return fmt.Errorf("unexpected PreHandlers type: %v", reflect.TypeOf(spec.Handler))
	}

	switch v := spec.Handler.(type) {
	case nil:
		return fmt.Errorf("handler must not be nil")
	case gin.HandlerFunc:
		handlers = append(handlers, v)
	case func(*gin.Context):
		handlers = append(handlers, v)
	default:
		return fmt.Errorf("handler must be gin.HandlerFunc or func(*gin.Context), got %v", reflect.TypeOf(spec.Handler))
	}

	a.Engine.Handle(spec.Method, swaggerToGinPattern(spec.Path), handlers...)
	return nil
}

func (a *Adapter) Serve(swagger *crud.Swagger, addr string) error {
	a.Engine.GET("/swagger.json", func(c *gin.Context) {
		c.JSON(200, swagger)
	})

	a.Engine.GET("/", func(c *gin.Context) {
		c.Header("content-type", "text/html")
		_, err := c.Writer.Write(crud.SwaggerUiTemplate)
		if err != nil {
			panic(err)
		}
	})

	return a.Engine.Run(addr)
}

// converts swagger endpoints /widget/{id} to gin endpoints /widget/:id
func swaggerToGinPattern(swaggerUrl string) string {
	return crud.SwaggerPathPattern.ReplaceAllString(swaggerUrl, ":$1")
}

func wrap(r *crud.Router, spec *crud.Spec) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := spec.Validate
		var query url.Values
		var body interface{}
		var path map[string]string

		if val.Path.Initialized() {
			path = map[string]string{}
			for _, param := range c.Params {
				path[param.Key] = param.Value
			}
		}

		if val.Body.Initialized() && val.Body.Kind() != crud.KindFile {
			if err := c.Bind(&body); err != nil {
				c.AbortWithStatusJSON(400, err.Error())
				return
			}
			defer func() {
				data, _ := json.Marshal(body)
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
			}()
		}

		if val.Query.Initialized() {
			query = c.Request.URL.Query()
			defer func() {
				c.Request.URL.RawQuery = query.Encode()
			}()
		}

		if err := r.Validate(val, query, body, path); err != nil {
			c.AbortWithStatusJSON(400, err.Error())
		}
	}
}
