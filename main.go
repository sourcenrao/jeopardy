package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sourcenrao/jeopardy/board"
	"github.com/sourcenrao/jeopardy/web"
)

var (
	categories int
	address    string
	filepath   string = "./data/clues.db"
)

func main() {

	flag.IntVar(&categories, "c", 6, "number of categories per round (default 6)")
	flag.StringVar(&address, "address", ":8000", "browser access address")
	flag.Parse()

	board, err := board.NewBoard(categories)
	if err != nil {
		log.Fatal(err)
	}

	err = board.LoadData(filepath)
	if err != nil {
		err = fmt.Errorf("data failed to load: %w", err)
		log.Fatal(err)
	}

	for _, category := range board.RoundOneColumns {
		fmt.Println(category.Category)
		for _, clue := range category.Clues {
			fmt.Println(clue.Value, clue.Question, clue.Answer)
		}
	}

	game, err := web.NewGame(board)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(game)
}
