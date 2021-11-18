package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/sourcenrao/jeopardy/api"
	_ "modernc.org/sqlite"
)

var (
	categoryCount int
	address       string
)

func main() {

	flag.IntVar(&categoryCount, "c", 6, "number of categoryCount per round (default 6)")
	flag.StringVar(&address, "address", ":8080", "browser access address (default localhost:8080)")
	flag.Parse()

	db, err := sql.Open("sqlite", "./data/clues.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	b := api.NewGame(categoryCount, db)

	api.Server(categoryCount, b, address, db)

}
