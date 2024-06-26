package route

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"food-review/pkg/db"
	"food-review/pkg/model"
	"food-review/pkg/template"
)

type Handler struct {
	Template     template.Templater
	ReviewDB     db.ReviewDBOpener
	DictionaryDB db.DictionaryDBOpener
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	db := h.ReviewDB.GetDB()
	targetReview, err := model.GetReview(db, uint(reviewIDu64))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("No Review with this ID"))
		return
	}

	err = h.Template.ExecuteTemplate(w, "review.html", targetReview)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) GetReviewsByKeyword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	reviewKeyword := r.URL.Query().Get("query")
	dict := h.DictionaryDB.GetDB()
	exist, err := model.KeywordExists(dict, reviewKeyword)
	if !exist {
		if err == nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Keyword not in dictionary"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		return
	}

	db := h.ReviewDB.GetDB()
	targetReviews, err := model.GetReviewsByKeyword(db, reviewKeyword)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("No review you are looking for"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.Template.ExecuteTemplate(w, "reviews_keyword.html", targetReviews)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) AccessReviewEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	reviewIDstr := mux.Vars(r)["reviewID"]
	reviewIDu64, err := strconv.ParseUint(reviewIDstr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	db := h.ReviewDB.GetDB()
	targetReview, err := model.GetReview(db, uint(reviewIDu64))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("No Review with this ID"))
		return
	}

	err = h.Template.ExecuteTemplate(w, "edit.html", targetReview)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) EditReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	reviewIDstr := mux.Vars(r)["reviewID"]
	reviewIDu64, err := strconv.ParseUint(reviewIDstr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	reviewBody, _ := ioutil.ReadAll(r.Body)
	db := h.ReviewDB.GetDB()
	err = model.UpdateReview(db, uint(reviewIDu64), reviewBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
