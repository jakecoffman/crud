package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/crud"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type Adapter struct {
	Engine *mux.Router
}

func New() *Adapter {
	return &Adapter{
		Engine: mux.NewRouter(),
	}
}

func (a *Adapter) Install(r *crud.Router, spec *crud.Spec) error {
	handlers := []mux.MiddlewareFunc{
		validateHandlerMiddleware(r, spec),
	}

	switch spec.PreHandlers.(type) {
	case nil:
	case []mux.MiddlewareFunc:
	case mux.MiddlewareFunc:
	case func(http.HandlerFunc) http.HandlerFunc:
	default:
		return fmt.Errorf("PreHandlers must be mux.MiddlewareFunc, got: %v", reflect.TypeOf(spec.Handler))
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

	// without Subrouter, Use would affect all other routes too
	subRouter := a.Engine.Path(spec.Path).Methods(spec.Method).Subrouter()
	subRouter.Use(handlers[:len(handlers)-1]...)
	subRouter.Handle("", finalHandler)

	return nil
}

func (a *Adapter) Serve(swagger *crud.Swagger, addr string) error {
	a.Engine.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(swagger)
	})

	a.Engine.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		_, err := w.Write(crud.SwaggerUiTemplate)
		if err != nil {
			panic(err)
		}
	})

	return http.ListenAndServe(addr, a.Engine)
}

func validateHandlerMiddleware(router *crud.Router, spec *crud.Spec) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := spec.Validate
			var query url.Values
			var body interface{}
			var path map[string]string

			if val.Path.Initialized() {
				path = map[string]string{}
				for key, value := range mux.Vars(r) {
					path[key] = value
				}
			}

			if val.Body.Initialized() && val.Body.Kind() != crud.KindFile {
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					w.WriteHeader(400)
					_ = json.NewEncoder(w).Encode(err.Error())
					return
				}
				defer func() {
					data, _ := json.Marshal(body)
					_ = r.Body.Close()
					r.Body = ioutil.NopCloser(bytes.NewReader(data))
				}()
			}

			if val.Query.Initialized() {
				query = r.URL.Query()
				defer func() {
					r.URL.RawQuery = query.Encode()
				}()
			}

			if err := router.Validate(val, query, body, path); err != nil {
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(err.Error())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
