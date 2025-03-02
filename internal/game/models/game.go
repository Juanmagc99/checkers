package models

import "errors"

type GameStatus string

const (
	StatusWaiting    GameStatus = "waiting"
	StatusInProgress GameStatus = "in_progress"
	StatusFinished   GameStatus = "finished"
)

type Game struct {
	ID           string     `json:"id"`
	Board        [][]string `json:"board"`
	Turn         string     `json:"turn"`
	Status       GameStatus `json:"status"`
	Player1Token string     `json:"player1_token,omitempty"`
	Player2Token string     `json:"player2_token,omitempty"`
}

/*
Create initial board with black and whites pieces
*/
func InitBoard() [][]string {
	board := make([][]string, 8)

	for i := 0; i < 8; i++ {
		board[i] = make([]string, 8)
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 8; j++ {
			if (i+j)%2 == 1 {
				board[i][j] = "B"
			}
		}
	}

	for i := 5; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if (i+j)%2 == 1 {
				board[i][j] = "W"
			}
		}
	}

	return board
}

func (s GameStatus) IsValid() bool {
	switch s {
	case StatusWaiting, StatusInProgress, StatusFinished:
		return true
	default:
		return false
	}
}

func (g *Game) ToSafeGame(requestToken string) (string, error) {

	if requestToken == g.Player1Token {
		g.Player2Token = ""
		return "player1", nil
	} else if requestToken == g.Player2Token {
		g.Player1Token = ""
		return "player2", nil
	} else {
		return "", errors.New("unauthorized access")
	}

}
