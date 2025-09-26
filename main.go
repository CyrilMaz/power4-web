package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Constantes du jeu
const (
	Rows    = 6
	Columns = 7
)

// Structure du jeu
type Game struct {
	Board   [Rows][Columns]int // 0 = vide, 1 = joueur 1, 2 = joueur 2
	Current int                // joueur actuel
	Winner  int                // 0 = pas encore, 1 ou 2 = gagnant, -1 = égalité
}

// Crée une nouvelle partie
func NewGame() *Game {
	return &Game{Current: 1}
}

// Jouer un coup dans une colonne
func (g *Game) Play(col int) {
	if g.Winner != 0 {
		return // partie terminée
	}
	if col < 0 || col >= Columns {
		return
	}
	for row := Rows - 1; row >= 0; row-- {
		if g.Board[row][col] == 0 {
			g.Board[row][col] = g.Current
			if g.checkWin(row, col) {
				g.Winner = g.Current
			} else if g.isDraw() {
				g.Winner = -1
			} else {
				g.Current = 3 - g.Current // alterne joueur 1 <-> 2
			}
			return
		}
	}
}

// Vérifie une victoire après un coup
func (g *Game) checkWin(r, c int) bool {
	player := g.Board[r][c]
	directions := [][2]int{
		{0, 1},  // horizontal
		{1, 0},  // vertical
		{1, 1},  // diagonale ↘
		{1, -1}, // diagonale ↙
	}
	for _, d := range directions {
		count := 1
		// vers l’avant
		count += g.countDirection(r, c, d[0], d[1], player)
		// vers l’arrière
		count += g.countDirection(r, c, -d[0], -d[1], player)
		if count >= 4 {
			return true
		}
	}
	return false
}

// Compte les jetons alignés dans une direction
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

// Vérifie si la grille est pleine
func (g *Game) isDraw() bool {
	for c := 0; c < Columns; c++ {
		if g.Board[0][c] == 0 {
			return false
		}
	}
	return true
}

// Variables globales
var (
	game *Game
	tmpl *template.Template
)

// HTML embarqué
const htmlTemplate = `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Puissance 4</title>
    <style>
        table { border-collapse: collapse; }
        td { width: 50px; height: 50px; text-align: center; border: 1px solid black; }
        .p1 { background: red; }
        .p2 { background: yellow; }
    </style>
</head>
<body>
    <h1>Puissance 4</h1>
    {{if eq .Winner 0}}
        <p>Joueur actuel : {{.Current}}</p>
    {{else if eq .Winner -1}}
        <p>Égalité !</p>
    {{else}}
        <p>🎉 Joueur {{.Winner}} a gagné !</p>
    {{end}}
    <form method="POST" action="/play">
        <table>
            {{range $i, $row := .Board}}
            <tr>
                {{range $j, $cell := $row}}
                <td class="{{if eq $cell 1}}p1{{else if eq $cell 2}}p2{{end}}"></td>
                {{end}}
            </tr>
            {{end}}
        </table>
        {{if eq .Winner 0}}
        <p>
            Choisir une colonne :
            <select name="column">
                <option value="0">0</option>
                <option value="1">1</option>
                <option value="2">2</option>
                <option value="3">3</option>
                <option value="4">4</option>
                <option value="5">5</option>
                <option value="6">6</option>
            </select>
            <button type="submit">Jouer</button>
        </p>
        {{end}}
    </form>
    <form method="POST" action="/reset">
        <button type="submit">🔄 Nouvelle partie</button>
    </form>
</body>
</html>
`

// Handlers HTTP
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, game)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		colStr := r.FormValue("column")
		col, err := strconv.Atoi(colStr)
		if err == nil {
			game.Play(col)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		game = NewGame()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	game = NewGame()
	tmpl = template.Must(template.New("index").Parse(htmlTemplate))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/reset", resetHandler)

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
