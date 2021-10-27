package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
}

func main() {
	filepath := "./data/clues.db"

	var clues Board

	err := clues.LoadData(filepath)
	if err != nil {
		fmt.Errorf("Data failed to load: %w", err)
	}

	err2 := clues.InitializeGame()
	if err2 != nil {
		fmt.Errorf("Failed to initialize game: %w", err)
	}
}

func (c *Board) InitializeGame() error {
	c.mu.Lock()
	numCat := len(c.AllCategories)
	fmt.Println(numCat)
	sixCategories := make([]string, 6) // Six categories per round
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 6; i += 1 {
		num := r.Intn(numCat - i)
		fmt.Println(num, c.AllCategories[num])
		sixCategories = append(sixCategories, c.AllCategories[num])
		c.AllCategories = append(c.AllCategories[:num], c.AllCategories[num+1:]...)
	}
	fmt.Println(sixCategories)
	c.mu.Unlock()
	return nil
}

func (c *Board) LoadData(filename string) error {
	c.mu.Lock()
	db, err := sql.Open("sqlite3", "./data/clues.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues")
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

	// file, err := os.Create("log.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.SetOutput(file)

	c.AllCategories = make([]string, len(categories))
	fmt.Println("length of allcategories: ", len(categories))
	for k := range categories {
		c.AllCategories = append(c.AllCategories, k)
		// log.Println(k)
	}

	c.mu.Unlock()

	return nil
}
