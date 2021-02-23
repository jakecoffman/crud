package widgets

import (
	"context"
	"github.com/jakecoffman/crud"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const port = "8888"

func TestRoutes(t *testing.T) {
	r := crud.NewRouter("test", "test")
	r.Add(Routes...)
	go func() {
		if err := r.Serve(":" + port); err != nil {
			log.Println(err)
		}
	}()

	res, err := get("/widgets?limit=a")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 400 {
		t.Fatal(res.StatusCode)
	}
	log.Println(string(data))
}

func get(url string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:"+port+url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
