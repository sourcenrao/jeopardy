package board

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
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
	AllCategories      []string
	RoundOneCategories []string
	RoundTwoCategories []string
	FinalJeopardy      Clue
}

func (c *Board) LoadData(filename string) error {
	c.mu.Lock()
	db, err := sql.Open("sqlite3", "./data/clues.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load 12 categories of clues
	rows, err := db.Query("SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues WHERE CATEGORY IN	(SELECT CATEGORY FROM clues ORDER BY random() LIMIT 12)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	categories := make(map[string]bool)

	for rows.Next() {
		var clue Clue
		err = rows.Scan(&clue.Value, &clue.Category, &clue.Comments, &clue.Answer, &clue.Question)
		if err != nil {
			log.Fatal(err)
		}
		c.AllClues = append(c.AllClues, clue)
		categories[clue.Category] = true
	}

	c.AllCategories = make([]string, 0, len(categories))
	for k := range categories {
		c.AllCategories = append(c.AllCategories, k)
	}

	// Need to load final jeopardy question

	c.mu.Unlock()

	return nil
}

func (c *Board) InitializeGame() error {
	c.mu.Lock()

	fmt.Println(reflect.TypeOf(c.AllCategories))
	fmt.Println(c.AllCategories)
	c.RoundOneCategories = c.AllCategories[:7]
	c.RoundTwoCategories = c.AllCategories[6:]
	fmt.Println(c.RoundOneCategories)
	fmt.Println(c.RoundTwoCategories)

	c.mu.Unlock()
	return nil
}
