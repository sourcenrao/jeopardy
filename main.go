package main

import (
	"database/sql"
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

	db, err := sql.Open("sqlite3", "./data/clues.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	board := NewGame(categoryCount, filepath, db)

	for _, category := range board.RoundOneColumns {
		for _, clue := range category.Clues {
			fmt.Println(clue.Value, clue.Question, clue.Answer)
		}
	}

	r := gin.Default()
	r.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to my Jeopardy game generator API, try /jeopardy to get unique data for a full game each refresh.")
	})
	r.GET("/jeopardy", func(c *gin.Context) {
		c.JSON(http.StatusOK, board)
		board = NewGame(categoryCount, filepath, db)
	})
	r.Run(address)

}

func NewGame(categoryCount int, filepath string, db *sql.DB) board.Board {
	board, err := board.NewBoard(categoryCount)
	if err != nil {
		log.Fatal(err)
	}

	err = board.LoadData(filepath, db)
	if err != nil {
		err = fmt.Errorf("data failed to load: %w", err)
		log.Fatal(err)
	}

	return board
}
