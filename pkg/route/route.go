package route

import (
	"fmt"
	"net/http"

	"food-review/pkg/db"
	"food-review/pkg/model"
	"food-review/pkg/template"
)

type Handler struct {
	Template template.Templater
	ReviewDB db.ReviewDBOpener
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func (h *Handler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	db := h.ReviewDB.GetDB()
	allReviews, err := model.GetAllReviews(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.Template.ExecuteTemplate(w, "reviews.html", allReviews)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
