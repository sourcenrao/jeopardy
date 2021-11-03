package board

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
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

// Create function for column generation

func (c *Board) LoadData(filename string) error {

	db, err := sql.Open("sqlite3", "./data/clues.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load numCategories*2 categories of clues for both rounds
	totalCategoriesStr := strconv.Itoa(c.NumCategories * 2)
	rows, err := db.Query(fmt.Sprintf(`SELECT CATEGORY FROM clues ORDER BY random() LIMIT %s`, totalCategoriesStr))
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

	// Find 6 questions of 200 value with distinct categories
	q := `SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues WHERE CATEGORY = ? AND VALUE = ? ORDER BY random() LIMIT 1`

	for _, category := range RoundOneCategories {
		var column BoardColumn
		column.Category = category
		column.Clues = make([]Clue, 0)
		for _, value := range roundOneValues {
			val := int(value)
			row, err := db.Query(q, category, val)
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
		c.RoundOneColumns = append(c.RoundOneColumns, column)
	}

	for _, category := range RoundTwoCategories {
		var column BoardColumn
		column.Category = category
		column.Clues = make([]Clue, 0)
		for _, value := range roundTwoValues {
			val := int(value)
			row, err := db.Query(q, category, val)
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
		c.RoundTwoColumns = append(c.RoundTwoColumns, column)
	}

	// Final Jeopardy clue
	finalq := `SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues WHERE VALUE > 2000 ORDER BY random() LIMIT 1`

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
