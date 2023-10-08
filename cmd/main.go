package main

import (
	"log"

	"food-review/pkg/http"
	"food-review/pkg/template"
)

func main() {
	err := template.LoadTemplates("template/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.StartServer()
}
