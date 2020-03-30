package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/eminetto/clean-architecture-go/pkg/bookmark/mock"
	"github.com/eminetto/clean-architecture-go/pkg/entity"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkIndex(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
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
	service.EXPECT().
		FindAll().
		Return([]*entity.Bookmark{b}, nil)
	apitest.New().
		Handler(bookmarkIndex(service)).
		Get("/v1/bookmark").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestBookmarkIndexNotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	service.EXPECT().
		Search("github").
		Return(nil, entity.ErrNotFound)

	apitest.New().
		Handler(bookmarkIndex(service)).
		Get("/v1/bookmark").
		Query("name", "github").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestBookmarkSearch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	b := &entity.Bookmark{
		Name:        "Elton Minetto",
		Description: "Minetto's page",
		Link:        "http://www.eltonminetto.net",
		Tags:        []string{"golang", "php", "linux", "mac"},
		Favorite:    true,
	}
	service.EXPECT().
		Search("minetto").
		Return([]*entity.Bookmark{b}, nil)
	apitest.New().
		Handler(bookmarkIndex(service)).
		Get("/v1/bookmark").
		Query("name", "minetto").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestBookmarkAdd(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeBookmarkHandlers(r, *n, service)
	path, err := r.GetRoute("bookmarkAdd").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/bookmark", path)
	service.EXPECT().
		Store(gomock.Any()).
		Return(entity.NewID(), nil)
	payload := fmt.Sprintf(`{
			"name": "Github",
			"description": "Github site",
			"link": "http://github.com",
			"tags": [
			  "git",
			  "social"
			]
		  }`)

	apitest.New().
		Handler(bookmarkAdd(service)).
		Post("/v1/bookmark").
		JSON(payload).
		Expect(t).
		Assert(jsonpath.Equal(`$.link`, "http://github.com")).
		Status(http.StatusCreated).
		End()
}

func TestBookmarkFind(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
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
	service.EXPECT().
		Find(b.ID).
		Return(b, nil)

	ts := httptest.NewServer(r)
	defer ts.Close()
	apitest.New().
		Handler(r).
		Get("/v1/bookmark/" + b.ID.String()).
		Expect(t).
		Assert(jsonpath.Equal(`$.id`, b.ID.String())).
		Status(http.StatusOK).
		End()
}

func TestBookmarkRemove(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeBookmarkHandlers(r, *n, service)
	path, err := r.GetRoute("bookmarkDelete").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/bookmark/{id}", path)
	b := &entity.Bookmark{
		ID:          entity.NewID(),
		Name:        "Elton Minetto",
		Description: "Minetto's page",
		Link:        "http://www.eltonminetto.net",
		Tags:        []string{"golang", "php", "linux", "mac"},
		Favorite:    false,
	}
	service.EXPECT().Delete(b.ID).Return(nil)
	ts := httptest.NewServer(r)
	defer ts.Close()

	apitest.New().
		Handler(r).
		Delete(ts.URL + "/v1/bookmark/" + b.ID.String()).
		Expect(t).
		Status(http.StatusOK).
		End()
}
