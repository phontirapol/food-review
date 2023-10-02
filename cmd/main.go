package main

import (
	"food-review/pkg/db"
	"food-review/pkg/http"
)

func main() {
	http.StartServer()
	db.InitReviewDB()
}
