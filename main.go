package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sourcenrao/jeopardy/board"
)

var (
	categoryCount int
	address       string
	filepath      string = "./data/clues.db"
)

func main() {

	flag.IntVar(&categoryCount, "c", 6, "number of categoryCount per round (default 6)")
	flag.StringVar(&address, "address", ":8080", "browser access address (default localhost:8080)")
	flag.Parse()

	board := NewGame(categoryCount, filepath)

	for _, category := range board.RoundOneColumns {
		for _, clue := range category.Clues {
			fmt.Println(clue.Value, clue.Question, clue.Answer)
		}
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/jeopardy", func(c *gin.Context) {
		c.JSON(http.StatusOK, board)
		board = NewGame(categoryCount, filepath)
	})
	r.Run(address)

}

func NewGame(categoryCount int, filepath string) board.Board {
	board, err := board.NewBoard(categoryCount)
	if err != nil {
		log.Fatal(err)
	}

	err = board.LoadData(filepath)
	if err != nil {
		err = fmt.Errorf("data failed to load: %w", err)
		log.Fatal(err)
	}

	return board
}
