package crud

import "testing"

func TestDuplicateRouteError(t *testing.T) {
	r := NewRouter("", "", &TestAdapter{})

	if err := r.Add(Spec{
		Method: "GET",
		Path:   "/widgets/{id}",
		Validate: Validate{Path: Object(map[string]Field{
			"id": Number(),
		})},
	}); err != nil {
		t.Errorf("expected no error, got %v", err.Error())
	}

	// duplicate of the one above
	if err := r.Add(Spec{
		Method: "GET",
		Path:   "/widgets/{id}",
		Validate: Validate{Path: Object(map[string]Field{
			"id": Number(),
		})},
	}); err == nil {
		t.Errorf("expected error")
	}
}
