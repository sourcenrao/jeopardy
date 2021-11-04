package board

import (
	"database/sql"
	"fmt"
	"log"
)

type Clue struct {
	Value    uint32
	Category string
	Comments string
	Answer   string
	Question string
}

type Board struct {
	RoundOneColumns []BoardColumn
	RoundTwoColumns []BoardColumn
	FinalJeopardy   Clue
	NumCategories   int
}

type BoardColumn struct {
	Category string
	Clues    []Clue
}

var (
	roundOneValues [5]int = [5]int{200, 400, 600, 800, 1000}
	roundTwoValues [5]int = [5]int{400, 800, 1200, 1600, 2000}
	q              string = `SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues WHERE CATEGORY = ? AND VALUE = ? ORDER BY random() LIMIT 1`
	finalq         string = `SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues WHERE VALUE > 2000 ORDER BY random() LIMIT 1`
)

// NewBoard that checks category count and constructs empty board
func NewBoard(numCategories int) (Board, error) {
	if numCategories < 3 {
		numCategories = 3
	} else if numCategories > 8 {
		numCategories = 8
	}
	var c Board
	c.NumCategories = int(numCategories)
	return c, nil
}

func (c *Board) LoadData(filename string) error {
	db, err := sql.Open("sqlite3", "./data/clues.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load numCategories*2 categories of clues for both rounds
	totalCategories := c.NumCategories * 2
	rows, err := db.Query(fmt.Sprintf(`SELECT CATEGORY FROM clues ORDER BY random() LIMIT %d`, totalCategories))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var allCategories []string
	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			log.Fatal(err)
		}
		allCategories = append(allCategories, category)
	}

	RoundOneCategories := allCategories[:c.NumCategories+1]
	RoundTwoCategories := allCategories[c.NumCategories:]

	c.RoundOneColumns, err = GetRoundColumns(RoundOneCategories, roundOneValues[:], db)
	if err != nil {
		log.Fatal(err)
	}
	c.RoundTwoColumns, err = GetRoundColumns(RoundTwoCategories, roundTwoValues[:], db)
	if err != nil {
		log.Fatal(err)
	}

	// Final Jeopardy clue

	row, err := db.Query(finalq)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	row.Next()
	err = row.Scan(&c.FinalJeopardy.Value, &c.FinalJeopardy.Category, &c.FinalJeopardy.Comments, &c.FinalJeopardy.Answer, &c.FinalJeopardy.Question)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(c)

	return nil
}

func GetRoundColumns(categories []string, values []int, db *sql.DB) ([]BoardColumn, error) {
	var roundColumns []BoardColumn
	for _, category := range categories {
		var column BoardColumn
		column.Category = category
		column.Clues = make([]Clue, 0)
		for _, value := range values {
			row, err := db.Query(q, category, value)
			if err != nil {
				log.Fatal(err)
			}
			defer row.Close()
			row.Next()
			var tempClue Clue
			err = row.Scan(&tempClue.Value, &tempClue.Category, &tempClue.Comments, &tempClue.Answer, &tempClue.Question)
			if err != nil {
				log.Fatal(err)
			}
			column.Clues = append(column.Clues, tempClue)
		}
		roundColumns = append(roundColumns, column)
	}
	return roundColumns, nil
}
