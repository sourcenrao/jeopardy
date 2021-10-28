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
	mu                 sync.Mutex
	AllClues           []Clue
	RoundOneCategories []string
	RoundTwoCategories []string
	FinalJeopardy      Clue
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

	c.RoundOneCategories = allCategories[:numCategories+1]
	c.RoundTwoCategories = allCategories[numCategories:]

	fmt.Println(c.RoundOneCategories)
	fmt.Println(c.RoundTwoCategories)

	// Find 6 questions of 200 value with distinct categories
	// `SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues WHERE CATEGORY IN (SELECT CATEGORY FROM clues ORDER BY random() LIMIT %s)`

	// finalRow, err := db.Query("SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues WHERE VALUE = 200 LIMIT 1)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer finalRow.Close()
	// fmt.Println(finalRow)

	// for finalRow.Next() {
	// 	var clue Clue
	// 	err = finalRow.Scan(&clue.Value, &clue.Category, &clue.Comments, &clue.Answer, &clue.Question)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	c.FinalJeopardy = clue
	// }

	c.mu.Unlock()

	return nil
}

// func (c *Board) InitializeGame() error {
// 	c.mu.Lock()

// 	c.mu.Unlock()
// 	return nil
// }
