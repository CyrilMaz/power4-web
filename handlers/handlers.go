package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"power-web/game"
)

var (
	Game = game.NewGame()
	mu   sync.Mutex
)

var tmpl = template.Must(template.New("graphic.html").Funcs(template.FuncMap{
	"isWinning": func(cells [][2]int, r, c int) bool {
		for _, cell := range cells {
			if cell[0] == r && cell[1] == c {
				return true
			}
		}
		return false
	},
}).ParseFiles("templates/graphic.html"))

func Home(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if err := tmpl.Execute(w, Game); err != nil {
		log.Println("Erreur template:", err)
	}
}

func Play(w http.ResponseWriter, r *http.Request) {
	col, _ := strconv.Atoi(r.URL.Query().Get("col"))
	mu.Lock()
	Game.Play(col)
	mu.Unlock()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Reset(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	Game = game.NewGame()
	mu.Unlock()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
