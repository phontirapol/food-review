package route

import (
	"fmt"
	"net/http"

	"food-review/pkg/model"
	"food-review/pkg/template"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func GetAllReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	allReviews, err := model.GetAllReviews()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = template.ExecuteTemplate(w, "reviews.html", allReviews)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
