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
