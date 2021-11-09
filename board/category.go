package board

import (
	"database/sql"
	"net/http"
)

var (
	catq string = `SELECT VALUE, CATEGORY, COMMENTS, ANSWER, QUESTION FROM clues WHERE CATEGORY = ? ORDER BY random() LIMIT ?`
)

func GetCluesForCategory(category string, numClues int, db *sql.DB) (int, []Clue) {
	var cluesForCat []Clue

	rows, err := db.Query(catq, category, numClues)
	if err != nil {
		return http.StatusBadRequest, cluesForCat
	}
	defer rows.Close()

	for rows.Next() {
		var tempClue Clue
		err = rows.Scan(&tempClue.Value, &tempClue.Category, &tempClue.Comments, &tempClue.Answer, &tempClue.Question)
		if err != nil {
			return http.StatusBadRequest, cluesForCat
		}
		cluesForCat = append(cluesForCat, tempClue)
	}

	return http.StatusOK, cluesForCat
}
