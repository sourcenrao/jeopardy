package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sourcenrao/jeopardy/board"
)

var (
	filepath string = "./data/clues.db"
	numClues int    = 10
	welcome  string = `Welcome to my Jeopardy game generator API, try /jeopardy to get unique data for a full game each refresh. 
You can use /jeopardy/(category name) to view 10 random clues from that category or view more than 10 using /jeopardy/(category name)/(quantity)
Flags:	-c <int>		| number of categories per round (default 6)
	-address <address:port>	| browser access address (default localhost:8080)`
)

func Server(categoryCount int, b board.Board, address string, db *sql.DB) {
	r := gin.Default()
	r.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, welcome)
	})
	r.GET("/jeopardy/:category", func(c *gin.Context) {
		httpStatus, cluesForCat := board.GetCluesForCategory(strings.ToUpper(c.Param("category")), numClues, db)
		if cluesForCat == nil {
			c.IndentedJSON(httpStatus, "You may have a typo, no resutls were found! Casing doesn't matter.")
		} else {
			c.IndentedJSON(httpStatus, cluesForCat)
		}
	})
	r.GET("/jeopardy/:category/:numClues", func(c *gin.Context) {
		tip := ""
		numClues, err := strconv.Atoi(c.Param("numClues"))
		if err != nil {
			numClues = 10
			tip = "WARNING: You may have entered an invalid number of clues, using default of 10."
		}
		httpStatus, cluesForCat := board.GetCluesForCategory(strings.ToUpper(c.Param("category")), numClues, db)
		if cluesForCat == nil {
			c.IndentedJSON(httpStatus, "You may have a typo, no resutls were found! Casing doesn't matter.")
		} else if tip != "" {
			cluesForCat[0].Comments = tip
			c.IndentedJSON(httpStatus, cluesForCat)
		} else {
			c.IndentedJSON(httpStatus, cluesForCat)
		}
	})
	r.GET("/jeopardy", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, b)
		b = NewGame(categoryCount, db)
	})
	r.Run(address)
}

func NewGame(categoryCount int, db *sql.DB) board.Board {
	b, err := board.NewBoard(categoryCount)
	if err != nil {
		log.Fatal(err)
	}

	err = b.LoadData(filepath, db)
	if err != nil {
		err = fmt.Errorf("data failed to load: %w", err)
		log.Fatal(err)
	}

	return b
}
