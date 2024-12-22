package main

import (
	"context"
	"github.com/canonical/go-dqlite/v2/app"
	"log"
	"testing"
)

func Test(t *testing.T) {
	dir := "/tmp/dqlite-data"
	address := "127.0.0.1:9001" // Unique node address

	// Set up Dqlite application
	app, err := app.New(dir, app.WithAddress(address))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("App created")

	// Create a database 'my-database' or just open it if
	// it already exists.
	db, err := app.Open(context.Background(), "my-database")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database created")

	db.Close()
	app.Close()
}
