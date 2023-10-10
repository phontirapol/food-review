package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

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

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	reviewIDstr := mux.Vars(r)["reviewID"]
	reviewIDu64, err := strconv.ParseUint(reviewIDstr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	db := h.ReviewDB.GetDB()
	targetReview, err := model.GetReview(db, uint(reviewIDu64))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.Template.ExecuteTemplate(w, "review.html", targetReview)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
