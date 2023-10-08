package http

import (
	"fmt"
	"log"
	"net/http"

	"food-review/pkg/route"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func initNewRouter() *mux.Router {
	newRouter := mux.NewRouter()
	newRouter.HandleFunc("/", route.Index).Methods("GET")
	newRouter.HandleFunc("/reviews", route.GetAllReviews).Methods("GET")

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
