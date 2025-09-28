package main

import (
	"fmt"
	"strconv"
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
	input int
)

func main() {
	game = NewGame()
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<table border="1">`)
		for r := 0; r < Rows; r++ {
			fmt.Fprint(w, "<tr>")
			for c := 0; c < Columns; c++ {
				color := "white"
				if game.Board[r][c] == 1 {
					color = "red"
				} else if game.Board[r][c] == 2 {
					color = "yellow"
				}
				if game.Winner != 0 {
					fmt.Fprintf(w, `<td style='width:50px;height:50px;background-color:%v'></td>`, color)
				} else {
					fmt.Fprintf(w, 
						`<td style='width:50px;height:50px;background-color:%v'>
						<a href="/play?column=%v">%v<a>
						</td>`, color, c, c)
				}
			}
			fmt.Fprint(w, "</tr>")
		}
		fmt.Fprint(w, "</table>")
		if game.Winner != 0 {
			fmt.Fprintf(w, `<p style="font-size:50px;font-width:bold">Player %v win</p>`, game.Winner)
		}
	})

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		col, _ := strconv.Atoi(r.URL.Query().Get("column"))
		game.Play(col)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
