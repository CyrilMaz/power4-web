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

type Game struct {
	Board        [Rows][Columns]int
	Current      int
	Winner       int
	LastRow      int
	LastCol      int
	WinningCells [][2]int
}

func NewGame() *Game {
	return &Game{
		Current: 1,
		LastRow: -1,
		LastCol: -1,
	}
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
		{0, 1},  // horizontal
		{1, 0},  // vertical
		{1, 1},  // diagonale ↘
		{1, -1}, // diagonale ↙
	}

	for _, d := range directions {
		cells := [][2]int{{r, c}}
		cells = append(cells, g.collect(r, c, d[0], d[1], player)...)
		cells = append(cells, g.collect(r, c, -d[0], -d[1], player)...)

		if len(cells) >= 4 {
			g.WinningCells = cells
			return true
		}
	}
	return false
}

// Helper pour collecter les positions alignées
func (g *Game) collect(r, c, dr, dc, player int) [][2]int {
	var cells [][2]int
	for {
		r += dr
		c += dc
		if r < 0 || r >= Rows || c < 0 || c >= Columns {
			break
		}
		if g.Board[r][c] != player {
			break
		}
		cells = append(cells, [2]int{r, c})
	}
	return cells
}

var (
	game = NewGame()
	mu   sync.Mutex
)

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"isWinning": func(cells [][2]int, r, c int) bool {
		for _, cell := range cells {
			if cell[0] == r && cell[1] == c {
				return true
			}
		}
		return false
	},
}).ParseFiles("templates/graphic.html"))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if err := tmpl.Execute(w, game); err != nil {
			log.Println("Erreur template:", err)
		}
	})

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		col, _ := strconv.Atoi(r.URL.Query().Get("col"))
		mu.Lock()
		game.Play(col)
		mu.Unlock()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		game = NewGame()
		mu.Unlock()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Serveur sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
