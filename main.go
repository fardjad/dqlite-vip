package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/canonical/go-dqlite/v2/app"
	"golang.org/x/sys/unix"
)

func main() {
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

	// Execute a SQL command on the database.
	// Creates a table 'my_table'
	if _, err := db.Exec("CREATE TABLE my_table (n INT)"); err != nil {
		log.Fatal(err)
	}
	log.Println("Table created")

	// wait until we received a termination signal
	ch := make(chan os.Signal, 32)
	signal.Notify(ch, unix.SIGPWR)
	signal.Notify(ch, unix.SIGINT)
	signal.Notify(ch, unix.SIGQUIT)
	signal.Notify(ch, unix.SIGTERM)

	<-ch

	db.Close()
	app.Close()
}
