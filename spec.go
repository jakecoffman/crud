package crud

import (
	"fmt"
	"strings"
)

// Spec is used to generate swagger paths and automatic handler validation
type Spec struct {
	// Method is the http method of the route, e.g. GET, POST
	Method string
	// Path is the URL path of the routes, w.g. /widgets/{id}
	Path string
	// PreHandlers run before validation. Good for authorization
	PreHandlers interface{}
	// Handler runs after validation. This is where you take over
	Handler interface{}
	// Description is the longer text that will appear in the Swagger under the endpoint
	Description string
	// Tags are how the Swagger groups paths together, e.g. []string{"Widgets"}
	Tags []string
	// Summary is a short description of what an endpoint does in the Swagger
	Summary string
	// Validate is used to automatically validate the various inputs to the endpoint
	Validate Validate
	// Responses specifies the responses in Swagger. If none provided a default is used.
	Responses map[string]Response
}

var methods = map[string]struct{}{
	"get":     {},
	"post":    {},
	"put":     {},
	"delete":  {},
	"options": {},
	"trace":   {},
	"patch":   {},
}

// Valid returns errors if the spec itself isn't valid. This helps finds bugs early.
func (s Spec) Valid() error {
	if _, ok := methods[strings.ToLower(s.Method)]; !ok {
		return fmt.Errorf("invalid method '%v'", s.Method)
	}

	params := pathParms(s.Path)
	if len(params) > 0 && !s.Validate.Path.Initialized() {
		return fmt.Errorf("path '%v' contains params but no path validation provided", s.Path)
	}

	if s.Validate.Path.Initialized() {
		if len(params) > 0 && s.Validate.Path.kind != KindObject {
			return fmt.Errorf("path must be an object")
		}
		// not ideal complexity but path params should be pretty small n
		for name := range s.Validate.Path.obj {
			var found bool
			for _, param := range params {
				if name == param {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("missing path param '%v' in url: '%v'", name, s.Path)
			}
		}

		for _, param := range params {
			if _, ok := s.Validate.Path.obj[param]; !ok {
				return fmt.Errorf("missing path validation '%v' in url: '%v'", param, s.Path)
			}
		}
	}

	return nil
}
