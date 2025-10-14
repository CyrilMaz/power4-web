package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

const (
	Rows    = 6
	Columns = 7
)

type Game struct {
	Board   [Rows][Columns]int `json:"board"`
	Current int                `json:"current"`
	Winner  int                `json:"winner"`
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
	game = NewGame()
	mu   sync.Mutex
)

func main() {
	// Serve the static HTML page at /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		// passer éventuellement l'état du jeu au template
		tmpl.Execute(w, game)
	})

	// Serve static files (script.js, style.css) directly
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))

	// Return current game state as JSON
	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	})

	// Play a column. Accepts GET ?column= or POST JSON {"column":n}
	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		var col int
		if r.Method == http.MethodPost {
			var body struct {
				Column int `json:"column"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
				col = body.Column
			}
		} else {
			col, _ = strconv.Atoi(r.URL.Query().Get("column"))
		}

		mu.Lock()
		game.Play(col)
		mu.Unlock()

		// return updated state
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	})

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var tmpl = template.Must(template.ParseFiles("templates/graphic.html"))
