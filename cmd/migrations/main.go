package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/eminetto/clean-architecture-go/config"
	_ "github.com/eminetto/clean-architecture-go/migrations"
	migrate "github.com/eminetto/mongo-migrate"
	"github.com/globalsign/mgo"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Missing options: up or down")
	}
	option := os.Args[1]

	session, err := mgo.Dial(config.MONGODB_HOST)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer session.Close()
	db := session.DB(config.MONGODB_DATABASE)
	migrate.SetDatabase(db)
	migrate.SetMigrationsCollection("migrations")
	migrate.SetLogger(log.New(os.Stdout, "INFO: ", 0))
	switch option {
	case "new":
		if len(os.Args) != 3 {
			log.Fatal("Should be: new description-of-migration")
		}
		fName := fmt.Sprintf("./migrations/%s_%s.go", time.Now().Format("20060102150405"), os.Args[2])
		from, err := os.Open("./migrations/template.go")
		if err != nil {
			log.Fatal("Should be: new description-of-migration")
		}
		defer from.Close()

		to, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("New migration created: %s\n", fName)
	case "up":
		err = migrate.Up(migrate.AllAvailable)
	case "down":
		err = migrate.Down(migrate.AllAvailable)
	}
	if err != nil {
		log.Fatal(err.Error())
	}
}
