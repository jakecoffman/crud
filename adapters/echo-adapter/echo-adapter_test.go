package adapter

import "testing"

func TestSwaggerToEcho(t *testing.T) {
	if "/widgets/:id" != swaggerToEchoPattern("/widgets/{id}") {
		t.Error(swaggerToEchoPattern("/widgets/{id}"))
	}
	if "/widgets/:id/sub/:subId" != swaggerToEchoPattern("/widgets/{id}/sub/{subId}") {
		t.Error(swaggerToEchoPattern("/widgets/{id}/sub/{subId}"))
	}
}
