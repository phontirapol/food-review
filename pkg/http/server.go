package http

import (
	"fmt"
	"log"
	"net/http"

	"food-review/pkg/db"
	"food-review/pkg/route"
	"food-review/pkg/template"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func initNewRouter() *mux.Router {
	newRouter := mux.NewRouter()

	tmpl := &template.Template{}
	if err := tmpl.LoadTemplates("template/*.html"); err != nil {
		log.Fatal(err)
	}
	var templater template.Templater = tmpl

	reviewDB := db.InitReviewDB()
	var reviewDBOpener db.ReviewDBOpener = reviewDB

	handler := &route.Handler{
		Template: templater,
		ReviewDB: reviewDBOpener,
	}

	newRouter.HandleFunc("/", handler.Index).Methods("GET")
	newRouter.HandleFunc("/reviews", handler.GetAllReviews).Methods("GET")

	return newRouter
}

func StartServer() {
	newRouter := initNewRouter()
	http.Handle("/", newRouter)

	fmt.Println("Someone has entered your website")

	err := http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
