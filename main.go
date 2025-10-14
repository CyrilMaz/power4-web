package main

import (
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

// --- Structure du jeu ---
type Game struct {
	Board   [Rows][Columns]int
	Current int
	Winner  int
	LastRow int
	LastCol int
}

func NewGame() *Game {
	return &Game{Current: 1, LastRow: -1, LastCol: -1}
}

func (g *Game) Play(col int) {
	if g.Winner != 0 || col < 0 || col >= Columns {
		return
	}
	for row := Rows - 1; row >= 0; row-- {
		if g.Board[row][col] == 0 {
			g.Board[row][col] = g.Current
			g.LastRow, g.LastCol = row, col
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

// --- Variables globales ---
var (
	game = NewGame()
	mu   sync.Mutex
	tmpl = template.Must(template.ParseFiles("templates/graphic.html"))
)

func main() {
	// Route principale : affiche le jeu
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		tmpl.Execute(w, game)
	})

	// Joue un coup dans une colonne
	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		col, _ := strconv.Atoi(r.URL.Query().Get("col"))
		mu.Lock()
		game.Play(col)
		mu.Unlock()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Réinitialise la partie
	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		game = NewGame()
		mu.Unlock()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Sert les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
