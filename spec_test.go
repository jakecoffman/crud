package crud

import "testing"

func TestSpec_Valid(t *testing.T) {
	specs := []Spec{
		{},
		{
			Method: "GET",
			Path:   "/1",
			Validate: Validate{Path: Object(map[string]Field{
				"id": Number(),
			})},
		},
		{
			Method: "GET",
			Path:   "/{id}",
		},
		{
			Method: "GET",
			Path:   "/3/{id}/4/{ok}",
			Validate: Validate{Path: Object(map[string]Field{
				"id": Number(),
			})},
		},
	}

	for _, spec := range specs {
		if err := spec.Valid(); err == nil {
			t.Errorf("expected error for path %v", spec.Path)
		}
	}
}
