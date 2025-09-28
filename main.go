package main

import (
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
			//afficheage de la grille
			http.HundleFunc("/grid", func(w ResponseWriter, r *Request) {
				fmt.Fprint(<w, "<table border='1'>")
				for _, r := range Rows {
					fmt.Fprint(w, "<tr>")
					for _, c := range r {
						color := "white"
						if g.Board[r][c] == 1 {
							color = "red"
						} else if g.Board[r][c] == 2 {
							color = "yellow"
						}
						fmt.Fprintf(w, "<td style='width:50px;height:50px;background-color:%s'></td>", color)
					}
					fmt.Fprint(w, "</tr>")
				}
				fmt.Fprint(w, "</table>")
			})
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
)

func main() {
	game = NewGame()

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
