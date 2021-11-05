package web

import (
	"net/http"

	"github.com/sourcenrao/jeopardy/board"
)

type JHandler struct {
	Message string
	Board   board.Board
}

func Server(j *JHandler) {

}

func NewGame(board board.Board) (JHandler, error) {
	var game JHandler
	game.Board = board
	game.Message = "new"
	return game, nil
}

func JServe(w http.ResponseWriter, r *http.Request, j *JHandler) {
	w.Write([]byte(j.Message))
}
