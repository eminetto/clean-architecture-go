package handler

import (
	"encoding/json"
	"fmt"
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

func bookmarkFind(service bookmark.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading bookmark"
		vars := mux.Vars(r)
		fmt.Println(vars)
		id := vars["id"]
		data, err := service.Find(entity.StringToID(id))
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

func bookmarkDelete(service bookmark.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing bookmark"
		vars := mux.Vars(r)
		id := vars["id"]
		err := service.Delete(entity.StringToID(id))
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
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
	)).Methods("GET", "OPTIONS").Name("bookmarkIndex")

	r.Handle("/v1/bookmark", n.With(
		negroni.Wrap(bookmarkAdd(service)),
	)).Methods("POST", "OPTIONS").Name("bookmarkAdd")

	r.Handle("/v1/bookmark/{id}", n.With(
		negroni.Wrap(bookmarkFind(service)),
	)).Methods("GET", "OPTIONS").Name("bookmarkFind")

	r.Handle("/v1/bookmark/{id}", n.With(
		negroni.Wrap(bookmarkDelete(service)),
	)).Methods("DELETE", "OPTIONS").Name("bookmarkDelete")
}
