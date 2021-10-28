package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sourcenrao/jeopardy/board"
)

func main() {
	filepath := "./data/clues.db"

	var board board.Board

	err := board.LoadData(filepath)
	if err != nil {
		err = fmt.Errorf("data failed to load: %w", err)
		log.Fatal(err)
	}

	err = board.InitializeGame()
	if err != nil {
		err = fmt.Errorf("failed to initialize game: %w", err)
		log.Fatal(err)
	}
}
