package board

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"
)

type Clue struct {
	Value    uint32
	Category string
	Comments string
	Answer   string
	Question string
}

type Board struct {
	mu              sync.Mutex
	RoundOneColumns []BoardColumn
	RoundTwoColumns []BoardColumn
	FinalJeopardy   Clue
}

type BoardColumn struct {
	Category string
	Clues    map[int]Clue
}

func (c *Board) LoadData(filename string, numCategories int) error {
	c.mu.Lock()

	db, err := sql.Open("sqlite3", "./data/clues.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load numCategories*2 categories of clues for both rounds
	numCategoriesStr := strconv.Itoa(numCategories * 2)
	rows, err := db.Query(fmt.Sprintf(`SELECT CATEGORY FROM clues ORDER BY random() LIMIT %s`, numCategoriesStr))
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

	RoundOneCategories := allCategories[:numCategories+1]
	RoundTwoCategories := allCategories[numCategories:]
	RoundOneValues := []int{200, 400, 600, 800, 1000}
	RoundTwoValues := []int{400, 800, 1200, 1600, 2000}

	// Find 6 questions of 200 value with distinct categories
	q := `SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues WHERE CATEGORY = ? AND VALUE = ? ORDER BY random() LIMIT 1`

	for _, category := range RoundOneCategories {
		var column BoardColumn
		column.Category = category
		column.Clues = make(map[int]Clue, 5)
		for _, value := range RoundOneValues {
			fmt.Print(category)
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
			column.Clues[int(tempClue.Value)] = tempClue
		}
		c.RoundOneColumns = append(c.RoundOneColumns, column)
		fmt.Print(c.RoundOneColumns)
	}

	for _, category := range RoundTwoCategories {
		var column BoardColumn
		column.Category = category
		column.Clues = make(map[int]Clue, 5)
		for _, value := range RoundTwoValues {
			fmt.Print(category)
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
			column.Clues[int(tempClue.Value)] = tempClue
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

	c.mu.Unlock()

	return nil
}

// func (c *Board) InitializeGame() error {
// 	c.mu.Lock()

// 	c.mu.Unlock()
// 	return nil
// }
