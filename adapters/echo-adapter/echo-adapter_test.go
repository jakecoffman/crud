package adapter

import (
	"fmt"
	"github.com/jakecoffman/crud"
	"github.com/jakecoffman/crud/adapters/echo-adapter/example/widgets"
	"net/http"
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

const address = "127.0.0.1:8080"

func TestExampleServer(t *testing.T) {
	r := crud.NewRouter("Widget API", "1.0.0", New())

	if err := r.Add(widgets.Routes...); err != nil {
		t.Fatal(err)
	}

	go func() {
		err := r.Serve(address)
		if err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("GET /widgets", func(t *testing.T) {
		// enforces the limit
		res := get("/widgets?limit=100")
		if res.StatusCode != 400 {
			t.Error(res.StatusCode)
		}

		res = get("/widgets?limit=25")
		if res.StatusCode != 200 {
			t.Error(res.StatusCode)
		}
	})
}

func get(url string) *http.Response {
	url = fmt.Sprintf("http://%s%s", address, url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return res
}
