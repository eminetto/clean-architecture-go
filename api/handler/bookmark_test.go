package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eminetto/clean-architecture-go/pkg/bookmark"
	"github.com/eminetto/clean-architecture-go/pkg/entity"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkIndex(t *testing.T) {
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
	assert.True(t, entity.IsValidID(b.ID.String()))
	assert.Equal(t, "http://github.com", b.Link)
	assert.False(t, b.CreatedAt.IsZero())

}
