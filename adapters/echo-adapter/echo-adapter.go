package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jakecoffman/crud"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"net/url"
	"reflect"
)

type Adapter struct {
	Echo *echo.Echo
}

func New() *Adapter {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return &Adapter{
		Echo: e,
	}
}

func (a *Adapter) Install(r *crud.Router, spec *crud.Spec) error {
	middlewares := []echo.MiddlewareFunc{wrap(r, spec)}

	switch v := spec.PreHandlers.(type) {
	case nil:
		// ok
	case []echo.MiddlewareFunc:
		middlewares = append(middlewares, v...)
	case echo.MiddlewareFunc:
		middlewares = append(middlewares, v)
	case func(echo.HandlerFunc) echo.HandlerFunc:
		middlewares = append(middlewares, v)
	default:
		return fmt.Errorf("unexpected PreHandlers type: %v", reflect.TypeOf(spec.Handler))
	}

	var handler echo.HandlerFunc
	switch v := spec.Handler.(type) {
	case nil:
		return fmt.Errorf("handler must not be nil")
	case echo.HandlerFunc:
		handler = v
	case func(echo.Context) error:
		handler = v
	default:
		return fmt.Errorf("handler must be echo.HandlerFunc or func(echo.Context) error, got %v", reflect.TypeOf(spec.Handler))
	}

	a.Echo.Add(spec.Method, swaggerToEchoPattern(spec.Path), handler, middlewares...)
	return nil
}

func (a *Adapter) Serve(swagger *crud.Swagger, addr string) error {
	a.Echo.GET("/swagger.json", func(c echo.Context) error {
		return c.JSON(200, swagger)
	})

	a.Echo.GET("/", func(c echo.Context) error {
		return c.HTML(200, string(crud.SwaggerUiTemplate))
	})

	return a.Echo.Start(addr)
}

// converts swagger endpoints /widget/{id} to echo endpoints /widget/:id
func swaggerToEchoPattern(swaggerUrl string) string {
	return crud.SwaggerPathPattern.ReplaceAllString(swaggerUrl, ":$1")
}

func wrap(r *crud.Router, spec *crud.Spec) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			val := spec.Validate
			var query url.Values
			var body interface{}
			var path map[string]string

			// need this scope so the defers run before next is called
			err := func() error {
				if val.Path.Initialized() {
					path = map[string]string{}
					for _, key := range c.ParamNames() {
						path[key] = c.Param(key)
					}
				}

				if val.Body.Initialized() && val.Body.Kind() != crud.KindFile {
					if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
						_ = c.JSON(400, err.Error())
						return fmt.Errorf("failed to bind body %w", err)
					}

					defer func() {
						data, err := json.Marshal(body)
						if err != nil {
							panic(err)
						}
						c.Request().Body = ioutil.NopCloser(bytes.NewReader(data))
					}()
				}

				if val.Query.Initialized() {
					query = c.Request().URL.Query()
					defer func() {
						c.Request().URL.RawQuery = query.Encode()
					}()
				}

				if err := r.Validate(val, query, body, path); err != nil {
					_ = c.JSON(400, err.Error())
					return err
				}

				return nil
			}()

			if err != nil {
				return err
			}

			return next(c)
		}
	}
}
