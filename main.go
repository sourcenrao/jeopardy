package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Clue struct {
	Round    uint8
	Value    uint32
	Double   bool
	Category string
	Comments string
	Answer   string
	Question string
}

func main() {
	filepath := "./data/clues.db"

	clues, err := LoadData(filepath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(clues[0])
}

func LoadData(filename string) ([]Clue, error) {
	db, err := sql.Open("sqlite3", "./data/clues.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM clues")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var clues []Clue

	for rows.Next() {
		var clue Clue
		err = rows.Scan(&clue.Round, &clue.Value, &clue.Double, &clue.Category, &clue.Comments, &clue.Answer, &clue.Question)
		if err != nil {
			log.Fatal(err)
		}
		clues = append(clues, clue)
	}

	return clues, nil
}
