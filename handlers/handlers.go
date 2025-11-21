package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/CyrilMaz/power4-web/game"
	"github.com/CyrilMaz/power4-web/theme"
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

type PageData struct {
	*game.Game
	Theme string
}

func Home(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	data := PageData{Game: Game, Theme: theme.GetTheme(r)}
	if err := tmpl.Execute(w, data); err != nil {
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

func UsePower(w http.ResponseWriter, r *http.Request) {
	power := r.URL.Query().Get("power")
	row, _ := strconv.Atoi(r.URL.Query().Get("row"))
	col, _ := strconv.Atoi(r.URL.Query().Get("col"))

	mu.Lock()
	Game.UsePower(Game.Current, power, row, col)
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Reset(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	Game = game.NewGame()
	mu.Unlock()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
