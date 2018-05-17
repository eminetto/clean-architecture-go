package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/eminetto/clean-architecture-go/pkg/bookmark"
	"github.com/eminetto/clean-architecture-go/pkg/entity"
	"github.com/joho/godotenv"
	"github.com/juju/mgosession"
	mgo "gopkg.in/mgo.v2"
)

func handleParams() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("Invalid query")
	}
	return os.Args[1], nil
}

func main() {
	query, err := handleParams()
	if err != nil {
		log.Fatal(err.Error())
	}
	env := os.Getenv("BOOKMARK_ENV")
	if env == "" {
		env = "dev"
	}
	err = godotenv.Load("config/" + env + ".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer session.Close()

	cPool, err := strconv.Atoi(os.Getenv("MONGODB_CONNECTION_POOL"))
	if err != nil {
		log.Println(err.Error())
		cPool = 10
	}
	mPool := mgosession.NewPool(nil, session, cPool)
	defer mPool.Close()

	bookmarkRepo := bookmark.NewMongoRepository(mPool)
	bookmarkService := bookmark.NewService(bookmarkRepo)
	all, err := bookmarkService.Search(query)
	if err != nil {
		log.Fatal(err)
	}
	if len(all) == 0 {
		log.Fatal(entity.ErrNotFound.Error())
	}
	for _, j := range all {
		fmt.Printf("%s %s %v \n", j.Name, j.Link, j.Tags)
	}
}
