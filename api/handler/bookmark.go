package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/eminetto/clean-architecture-go/pkg/bookmark"
	"github.com/eminetto/clean-architecture-go/pkg/entity"
	"github.com/gorilla/mux"
)

func bookmarkIndex(service bookmark.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading bookmarks"
		var data []*entity.Bookmark
		var err error
		name := r.URL.Query().Get("name")
		switch {
		case name == "":
			data, err = service.FindAll()
		default:
			data, err = service.Search(name)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func bookmarkAdd(service bookmark.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding bookmark"
		var b *entity.Bookmark
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		b.ID, err = service.Store(b)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(b); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

//MakeBookmarkHandlers make url handlers
func MakeBookmarkHandlers(r *mux.Router, n negroni.Negroni, service bookmark.UseCase) {
	r.Handle("/v1/bookmark", n.With(
		negroni.Wrap(bookmarkIndex(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/bookmark", n.With(
		negroni.Wrap(bookmarkAdd(service)),
	)).Methods("POST", "OPTIONS")
}
