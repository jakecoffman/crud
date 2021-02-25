package main

import (
	"github.com/jakecoffman/crud"
	"github.com/jakecoffman/crud/example/widgets"
	"log"
)

func main() {
	r := crud.NewRouter("Widget API", "1.0.0")

	if err := r.Add(widgets.Routes...); err != nil {
		log.Fatal(err)
	}

	log.Println("Serving http://127.0.0.1:8080")
	err := r.Serve(":8080")
	if err != nil {
		log.Println(err)
	}
}
