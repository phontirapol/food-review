package http

import (
	"fmt"
	"log"
	"net/http"

	"food-review/pkg/route"
	"food-review/pkg/template"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func initNewRouter(tmpl *template.Template) *mux.Router {
	newRouter := mux.NewRouter()
	var templater template.Templater = tmpl

	handler := &route.Handler{
		Template: templater,
	}

	newRouter.HandleFunc("/", handler.Index).Methods("GET")
	newRouter.HandleFunc("/reviews", handler.GetAllReviews).Methods("GET")

	return newRouter
}

func StartServer(tmpl *template.Template) {
	newRouter := initNewRouter(tmpl)
	http.Handle("/", newRouter)

	fmt.Println("Someone has entered your website")

	err := http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
