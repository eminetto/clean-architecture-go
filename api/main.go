package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/eminetto/clean-architecture-go/api/handler"
	"github.com/eminetto/clean-architecture-go/pkg/bookmark"
	"github.com/eminetto/clean-architecture-go/pkg/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/juju/mgosession"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	env := os.Getenv("BOOKMARK_ENV")
	if env == "" {
		env = "dev"
	}
	err := godotenv.Load("config/" + env + ".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer session.Close()

	r := mux.NewRouter()

	cPool, err := strconv.Atoi(os.Getenv("MONGODB_CONNECTION_POOL"))
	if err != nil {
		log.Println(err.Error())
		cPool = 10
	}
	mPool := mgosession.NewPool(nil, session, cPool)
	defer mPool.Close()

	bookmarkRepo := bookmark.NewMongoRepository(mPool)
	bookmarkService := bookmark.NewService(bookmarkRepo)

	//handlers
	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),
		negroni.NewLogger(),
	)
	//bookmark
	handler.MakeBookmarkHandlers(r, *n, bookmarkService)

	http.Handle("/", r)
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + os.Getenv("API_PORT"),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
