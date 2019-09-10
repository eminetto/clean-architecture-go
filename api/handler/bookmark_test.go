package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/eminetto/clean-architecture-go/pkg/bookmark"
	"github.com/eminetto/clean-architecture-go/pkg/entity"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkIndex(t *testing.T) {
	repo := bookmark.NewInmemRepository()
	service := bookmark.NewService(repo)
	r := mux.NewRouter()
	n := negroni.New()
	MakeBookmarkHandlers(r, *n, service)
	path, err := r.GetRoute("bookmarkIndex").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/bookmark", path)
	b := &entity.Bookmark{
		Name:        "Elton Minetto",
		Description: "Minetto's page",
		Link:        "http://www.eltonminetto.net",
		Tags:        []string{"golang", "php", "linux", "mac"},
		Favorite:    true,
	}
	_, _ = service.Store(b)
	ts := httptest.NewServer(bookmarkIndex(service))
	defer ts.Close()
	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestBookmarkIndexNotFound(t *testing.T) {
	repo := bookmark.NewInmemRepository()
	service := bookmark.NewService(repo)
	ts := httptest.NewServer(bookmarkIndex(service))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?name=github")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestBookmarkSearch(t *testing.T) {
	repo := bookmark.NewInmemRepository()
	service := bookmark.NewService(repo)
	b := &entity.Bookmark{
		Name:        "Elton Minetto",
		Description: "Minetto's page",
		Link:        "http://www.eltonminetto.net",
		Tags:        []string{"golang", "php", "linux", "mac"},
		Favorite:    true,
	}
	_, _ = service.Store(b)
	ts := httptest.NewServer(bookmarkIndex(service))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?name=minetto")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestBookmarkAdd(t *testing.T) {
	repo := bookmark.NewInmemRepository()
	service := bookmark.NewService(repo)
	r := mux.NewRouter()
	n := negroni.New()
	MakeBookmarkHandlers(r, *n, service)
	path, err := r.GetRoute("bookmarkAdd").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/bookmark", path)

	h := bookmarkAdd(service)

	ts := httptest.NewServer(h)
	defer ts.Close()
	payload := fmt.Sprintf(`{
  "name": "Github",
  "description": "Github site",
  "link": "http://github.com",
  "tags": [
    "git",
    "social"
  ]
}`)
	resp, _ := http.Post(ts.URL+"/v1/bookmark", "application/json", strings.NewReader(payload))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var b *entity.Bookmark
	json.NewDecoder(resp.Body).Decode(&b)
	assert.True(t, entity.IsValidID(b.ID.Hex()))
	assert.Equal(t, "http://github.com", b.Link)
	assert.False(t, b.CreatedAt.IsZero())
}

func TestBookmarkFind(t *testing.T) {
	repo := bookmark.NewInmemRepository()
	service := bookmark.NewService(repo)
	r := mux.NewRouter()
	n := negroni.New()
	MakeBookmarkHandlers(r, *n, service)
	path, err := r.GetRoute("bookmarkFind").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/bookmark/{id}", path)
	b := &entity.Bookmark{
		ID:          entity.NewID(),
		Name:        "Elton Minetto",
		Description: "Minetto's page",
		Link:        "http://www.eltonminetto.net",
		Tags:        []string{"golang", "php", "linux", "mac"},
		Favorite:    true,
	}
	bID, _ := service.Store(b)
	handler := bookmarkFind(service)
	r.Handle("/v1/bookmark/{id}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/bookmark/" + bID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var d *entity.Bookmark
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)
	assert.Equal(t, bID, d.ID)
}

func TestBookmarkRemove(t *testing.T) {
	repo := bookmark.NewInmemRepository()
	service := bookmark.NewService(repo)
	r := mux.NewRouter()
	n := negroni.New()
	MakeBookmarkHandlers(r, *n, service)
	path, err := r.GetRoute("bookmarkDelete").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/bookmark/{id}", path)
	b := &entity.Bookmark{
		Name:        "Elton Minetto",
		Description: "Minetto's page",
		Link:        "http://www.eltonminetto.net",
		Tags:        []string{"golang", "php", "linux", "mac"},
		Favorite:    false,
	}
	bID, _ := service.Store(b)
	handler := bookmarkDelete(service)
	req, _ := http.NewRequest("DELETE", "/v1/bookmark/"+bID.Hex(), nil)
	r.Handle("/v1/bookmark/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
