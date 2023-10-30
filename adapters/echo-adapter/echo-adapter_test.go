package adapter

import (
	"github.com/jakecoffman/crud"
	"github.com/jakecoffman/crud/adapters/echo-adapter/example/widgets"
	"net/http/httptest"
	"testing"
)

func TestSwaggerToEcho(t *testing.T) {
	if "/widgets/:id" != swaggerToEchoPattern("/widgets/{id}") {
		t.Error(swaggerToEchoPattern("/widgets/{id}"))
	}
	if "/widgets/:id/sub/:subId" != swaggerToEchoPattern("/widgets/{id}/sub/{subId}") {
		t.Error(swaggerToEchoPattern("/widgets/{id}/sub/{subId}"))
	}
}

func TestExampleServer(t *testing.T) {
	adapter := New()
	router := crud.NewRouter("Widget API", "1.0.0", adapter)

	if err := router.Add(widgets.Routes...); err != nil {
		t.Fatal(err)
	}

	t.Run("GET /widgets", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/widgets?limit=100", nil)
		w := httptest.NewRecorder()
		adapter.Echo.ServeHTTP(w, r)

		if w.Result().StatusCode != 400 {
			t.Error(w.Result().StatusCode)
		}

		r = httptest.NewRequest("GET", "/widgets?limit=25", nil)
		w = httptest.NewRecorder()
		adapter.Echo.ServeHTTP(w, r)

		if w.Result().StatusCode != 200 {
			t.Error(w.Result().StatusCode)
		}
	})
}
