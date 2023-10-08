package main

import (
	"log"

	"food-review/pkg/http"
	"food-review/pkg/template"
)

func main() {
	tmpl := &template.Template{}
	err := tmpl.LoadTemplates("template/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.StartServer(tmpl)
}
