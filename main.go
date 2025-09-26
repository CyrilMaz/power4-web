package main

import (
	"html/template"
	"log"
	"net/http"
)

const (
	Rows    = 6
	Columns = 7
)

type Game struct {
	Board   [Rows][Columns]int
	Current int
	Winner  int
}

func NewGame() *Game {
	return &Game{Current: 1}
}

func (g *Game) Play(col int) {
	if g.Winner != 0 {
		return
	}
	if col < 0 || col >= Columns {
		return
	}
	for row := Rows - 1; row >= 0; row-- {
		if g.Board[row][col] == 0 {
			g.Board[row][col] = g.Current
			if g.checkWin(row, col) {
				g.Winner = g.Current
			} else {
				g.Current = 3 - g.Current
			}
			return
		}
	}
}

func (g *Game) checkWin(r, c int) bool {
	player := g.Board[r][c]
	directions := [][2]int{
		{0, 1},
		{1, 0},
		{1, 1},
		{1, -1},
	}
	for _, d := range directions {
		count := 1
		count += g.countDirection(r, c, d[0], d[1], player)
		count += g.countDirection(r, c, -d[0], -d[1], player)
		if count >= 4 {
			return true
		}
	}
	return false
}

func (g *Game) countDirection(r, c, dr, dc, player int) int {
	count := 0
	for {
		r += dr
		c += dc
		if r < 0 || r >= Rows || c < 0 || c >= Columns {
			break
		}
		if g.Board[r][c] != player {
			break
		}
		count++
	}
	return count
}

var (
	game *Game
	tmpl *template.Template
)

func main() {
	game = NewGame()
	tmpl = template.Must(template.New("index").Parse(htmlTemplate))

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
