package crud

import "testing"

func TestSwaggerToGin(t *testing.T) {
	if "/widgets/:id" != swaggerToGinPattern("/widgets/{id}") {
		t.Error(swaggerToGinPattern("/widgets/{id}"))
	}
	if "/widgets/:id/sub/:subId" != swaggerToGinPattern("/widgets/{id}/sub/{subId}") {
		t.Error(swaggerToGinPattern("/widgets/{id}/sub/{subId}"))
	}
}

func TestValidationOnAdd(t *testing.T) {
	r := NewRouter("", "")
	if err := r.Add(Spec{}); err == nil {
		t.Errorf("expected error")
	}

	if err := r.Add(Spec{
		Method: "GET",
		Path:   "/widgets",
		Validate: Validate{Path: Object(map[string]Field{
			"id": Number(),
		})},
	}); err == nil {
		t.Errorf("expected error")
	}

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
