package board

import (
	"database/sql"
	"log"
	"sync"
)

type Clue struct {
	Value    uint32 `json:"value"`
	Category string `json:"category"`
	Comments string `json:"comment"`
	Answer   string `json:"answer"`
	Question string `json:"question"`
}

type Board struct {
	RoundOneColumns []BoardColumn `json:"round1"`
	RoundTwoColumns []BoardColumn `json:"round2"`
	FinalJeopardy   Clue          `json:"final"`
	NumCategories   int           `json:"numcategories"`
}

type BoardColumn struct {
	Clues []Clue `json:"clues"`
}

/* Example:
{
	"round1" : [
		"clues" : [
			{
				Clue object
			},
			{
				Clue object
			},
			{etc}
			],
			[
				{},{},{}
			],
			[etc]
	],
	"round2" : [columns of clues],
	"final" : {final jeopardy clue}
}
*/

var (
	roundOneValues [5]int = [5]int{200, 400, 600, 800, 1000}
	roundTwoValues [5]int = [5]int{400, 800, 1200, 1600, 2000}
	categoryq      string = `SELECT CATEGORY FROM clues ORDER BY random() LIMIT ?`
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

func (c *Board) LoadData(filename string, db *sql.DB) error {

	// Load numCategories*2 categories of clues for both rounds
	totalCategories := c.NumCategories * 2
	rows, err := db.Query(categoryq, totalCategories)
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

	RoundOneCategories := allCategories[:c.NumCategories]
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

	return nil
}

func GetRoundColumns(categories []string, values []int, db *sql.DB) ([]BoardColumn, error) {
	var roundColumns []BoardColumn
	// Turn this loop into a goroutine
	cats := make(chan BoardColumn, len(categories))
	var wg sync.WaitGroup
	for _, category := range categories {
		wg.Add(1)
		go func(category string) {
			var column BoardColumn
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
			// move append outside of loop and get columns from channel
			cats <- column
		}(category)

		go func() {
			cat := <-cats
			roundColumns = append(roundColumns, cat)
			wg.Done()
		}()
	}
	wg.Wait()
	return roundColumns, nil
}
