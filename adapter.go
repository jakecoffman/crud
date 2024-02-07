package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

type ServeMuxAdapter struct {
	Engine *http.ServeMux
}

func NewServeMuxAdapter() *ServeMuxAdapter {
	return &ServeMuxAdapter{
		Engine: http.NewServeMux(),
	}
}

type MiddlewareFunc func(http.Handler) http.Handler

func (a *ServeMuxAdapter) Install(r *Router, spec *Spec) error {
	middlewares := []MiddlewareFunc{
		validateHandlerMiddleware(r, spec),
	}

	switch v := spec.PreHandlers.(type) {
	case nil:
	case []MiddlewareFunc:
		middlewares = append(middlewares, v...)
	case MiddlewareFunc:
		middlewares = append(middlewares, v)
	case func(http.Handler) http.Handler:
		middlewares = append(middlewares, v)
	default:
		return fmt.Errorf("PreHandlers must be MiddlewareFunc, got: %v", reflect.TypeOf(spec.Handler))
	}

	var finalHandler http.Handler
	switch v := spec.Handler.(type) {
	case nil:
		return fmt.Errorf("handler must not be nil")
	case http.HandlerFunc:
		finalHandler = v
	case func(http.ResponseWriter, *http.Request):
		finalHandler = http.HandlerFunc(v)
	case http.Handler:
		finalHandler = v
	default:
		return fmt.Errorf("handler must be http.HandlerFunc, got %v", reflect.TypeOf(spec.Handler))
	}

	// install the route, use a subrouter so the "use" is scoped
	path := fmt.Sprintf("%s %s", spec.Method, spec.Path)
	subrouter := http.NewServeMux()
	subrouter.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		handler := finalHandler
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		handler.ServeHTTP(w, r)
	})
	a.Engine.Handle(path, subrouter)

	return nil
}

func (a *ServeMuxAdapter) Serve(swagger *Swagger, addr string) error {
	a.Engine.HandleFunc("GET /swagger.json", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(swagger)
	})

	a.Engine.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		_, err := w.Write(SwaggerUiTemplate)
		if err != nil {
			panic(err)
		}
	})

	return http.ListenAndServe(addr, a.Engine)
}

func validateHandlerMiddleware(router *Router, spec *Spec) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := spec.Validate
			var query url.Values
			var body interface{}
			var path map[string]string

			if val.Path.Initialized() {
				path = map[string]string{}
				for name := range val.Path.obj {
					path[name] = r.PathValue(name)
				}
			}

			var rewriteBody bool
			if val.Body.Initialized() && val.Body.Kind() != KindFile {
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					w.WriteHeader(400)
					_ = json.NewEncoder(w).Encode("failure decoding body: " + err.Error())
					return
				}
				rewriteBody = true
			}

			var rewriteQuery bool
			if val.Query.Initialized() {
				query = r.URL.Query()
				rewriteQuery = true
			}

			if err := router.Validate(val, query, body, path); err != nil {
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(err.Error())
				return
			}

			// Validate can strip values that are not valid, so we rewrite them
			// after validation is complete. Can't use defer as in other adapters
			// because next.ServeHTTP calls the next handler and defer hasn't
			// run yet.
			if rewriteBody {
				data, _ := json.Marshal(body)
				_ = r.Body.Close()
				r.Body = io.NopCloser(bytes.NewReader(data))
			}
			if rewriteQuery {
				r.URL.RawQuery = query.Encode()
			}

			next.ServeHTTP(w, r)
		})
	}
}
