package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sourcenrao/jeopardy/board"
)

func main() {
	var categories int
	flag.IntVar(&categories, "c", 6, "number of categories per round (default 6)")
	flag.Parse()

	filepath := "./data/clues.db"

	var board board.Board

	err := board.LoadData(filepath, categories)
	if err != nil {
		err = fmt.Errorf("data failed to load: %w", err)
		log.Fatal(err)
	}

	// err = board.InitializeGame()
	// if err != nil {
	// 	err = fmt.Errorf("failed to initialize game: %w", err)
	// 	log.Fatal(err)
	// }

	fmt.Println(&board)
}
