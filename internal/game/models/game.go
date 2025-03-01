package models

type Game struct {
	ID     string     `json:"id"`
	Board  [][]string `json:"board"`
	Turn   string     `json:"turn"`
	Status string     `json:"status"`
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
