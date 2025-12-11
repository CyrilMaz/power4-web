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

type TemplateData struct {
	*game.Game
	Theme string
}

func Home(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	currentTheme := theme.GetTheme(r)
	data := TemplateData{
		Game:  Game,
		Theme: currentTheme,
	}
	mu.Unlock()
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Erreur template:", err)
	}
}

func Play(w http.ResponseWriter, r *http.Request) {
	colStr := r.URL.Query().Get("col")
	col, err := strconv.Atoi(colStr)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if col < 0 || col >= game.Columns {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mu.Lock()
	Game.Play(col)
	mu.Unlock()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UsePower(w http.ResponseWriter, r *http.Request) {
	power := r.URL.Query().Get("power")
	rowStr := r.URL.Query().Get("row")
	colStr := r.URL.Query().Get("col")

	row, errRow := strconv.Atoi(rowStr)
	col, errCol := strconv.Atoi(colStr)
	if errCol != nil || col < 0 || col >= game.Columns {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// row may be -1 for column-only powers
	if errRow != nil {
		row = -1
	}

	mu.Lock()
	// validation supplémentaire selon le pouvoir
	success := Game.UsePower(Game.Current, power, row, col)
	mu.Unlock()

	_ = success // pour l'instant on redirige toujours vers / ; on pourrait afficher un message
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ToggleTheme switches between light and dark
func ToggleTheme(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	if Game.Theme == "light" {
		Game.Theme = "dark"
	} else {
		Game.Theme = "light"
	}
	mu.Unlock()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Reset(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	Game = game.NewGame()
	mu.Unlock()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
