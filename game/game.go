package game

const (
	Rows    = 6
	Columns = 7
)

type Power struct {
	Name    string
	Uses    int
	MaxUses int
	// Cooldown could be added later
}

type Game struct {
	Board        [Rows][Columns]int
	Current      int
	Winner       int
	LastRow      int
	LastCol      int
	WinningCells [][2]int
	Powers       map[int][]Power
	Theme        string
}

func NewGame() *Game {
	return &Game{
		Current: 1,
		LastRow: -1,
		LastCol: -1,
		Theme:   "light", // default theme
		Powers: map[int][]Power{
			1: {
				{Name: "Détruire", Uses: 2, MaxUses: 2},
				{Name: "Échanger", Uses: 1, MaxUses: 1},
				{Name: "Bloquer", Uses: 1, MaxUses: 1},
			},
			2: {
				{Name: "Détruire", Uses: 2, MaxUses: 2},
				{Name: "Échanger", Uses: 1, MaxUses: 1},
				{Name: "Bloquer", Uses: 1, MaxUses: 1},
			},
		},
	}
}

// Play pose un jeton dans la colonne si possible. Si victoire -> Winner, sinon change de joueur.
func (g *Game) Play(col int) {
	if g.Winner != 0 || col < 0 || col >= Columns {
		return
	}
	for row := Rows - 1; row >= 0; row-- {
		// une case bloquée (3) empêche la pose au-dessus d'elle
		if g.Board[row][col] == 0 {
			// poser sur cette case uniquement si ce n'est pas un bloc fixe
			g.Board[row][col] = g.Current
			g.LastRow, g.LastCol = row, col
			if g.checkWin(row, col) {
				g.Winner = g.Current
			} else {
				// nettoie les anciennes cellules gagnantes si le coup n'est pas gagnant
				g.WinningCells = nil
				g.Current = 3 - g.Current
			}
			return
		}
		// si case == 3 (blocage) -> on ne peut pas passer au-dessus, donc continuer vers le haut
	}
}

// UsePower vérifie et exécute le pouvoir demandé par `player`.
// row/col sont les coordonnées ciblées selon le pouvoir.
func (g *Game) UsePower(player int, powerName string, row, col int) bool {
	if g.Winner != 0 || player != g.Current {
		return false
	}
	// recherche du pouvoir disponible
	powerIndex := -1
	for i, p := range g.Powers[player] {
		if p.Name == powerName && p.Uses > 0 {
			powerIndex = i
			break
		}
	}
	if powerIndex == -1 {
		return false
	}

	success := false
	switch powerName {
	case "Détruire":
		success = g.destroyPiece(row, col)
	case "Échanger":
		success = g.swapPieces(row, col)
	case "Bloquer":
		success = g.blockColumn(col)
	default:
		return false
	}

	if success {
		// décrémenter l'usage
		g.Powers[player][powerIndex].Uses--

		// 🔥 Vérifier si un alignement apparaît APRÈS le pouvoir
		for r := 0; r < Rows; r++ {
			for c := 0; c < Columns; c++ {
				if g.Board[r][c] == player && g.checkWin(r, c) {
					g.Winner = player
					return true
				}
			}
		}

		// Sinon : nettoyer et passer au joueur suivant
		g.WinningCells = nil
		g.Current = 3 - g.Current
	}
	return success
}

// destroyPiece : détruit une pièce adverse (ne peut pas détruire une case = 3 ni une case vide ni la sienne)
func (g *Game) destroyPiece(row, col int) bool {
	if row < 0 || row >= Rows || col < 0 || col >= Columns {
		return false
	}
	val := g.Board[row][col]
	if val == 0 || val == 3 || val == g.Current {
		return false
	}
	// détruire
	g.Board[row][col] = 0
	g.applyGravity()
	return true
}

// swapPieces : échange la pièce à (row,col) avec celle en dessous (row+1,col)
// n'autorise pas d'échange impliquant une case vide ou un bloc (3)
func (g *Game) swapPieces(row, col int) bool {
	if row < 0 || row >= Rows-1 || col < 0 || col >= Columns {
		return false
	}
	v1 := g.Board[row][col]
	v2 := g.Board[row+1][col]
	if v1 == 0 || v2 == 0 || v1 == 3 || v2 == 3 {
		return false
	}
	// échange
	g.Board[row][col], g.Board[row+1][col] = v2, v1
	return true
}

// blockColumn : place un bloc fixe (valeur 3) sur la première case vide la plus basse
// Le bloc ne doit pas tomber (applyGravity le respectera)
func (g *Game) blockColumn(col int) bool {
	if col < 0 || col >= Columns {
		return false
	}
	for row := Rows - 1; row >= 0; row-- {
		if g.Board[row][col] == 0 {
			g.Board[row][col] = 3
			return true
		}
		// si rencontre un bloc 3 déjà présent, on considère la colonne pleine
		if g.Board[row][col] == 3 {
			return false
		}
	}
	return false
}

// applyGravity fait tomber les jetons (1 et 2) mais laisse les blocs (3) en place
func (g *Game) applyGravity() {
	for col := 0; col < Columns; col++ {
		writePos := Rows - 1
		// parcourir de bas en haut
		for row := Rows - 1; row >= 0; row-- {
			val := g.Board[row][col]
			if val == 3 {
				// bloc fixe : le poser là et avancer writePos au dessus du bloc
				writePos = row - 1
				continue
			}
			if val != 0 {
				// déplacer la pièce vers writePos si nécessaire
				if row != writePos {
					g.Board[writePos][col] = val
					g.Board[row][col] = 0
				}
				writePos--
			}
		}
		// nettoyer au-dessus
		for r := writePos; r >= 0; r-- {
			if g.Board[r][col] != 3 { // ne pas écraser un bloc déjà présent (mais il n'y en aura pas)
				g.Board[r][col] = 0
			}
		}
	}
}

// checkWin : similaire à la version robuste — ignore les blocs (3) et sélectionne exactement 4 contiguës contenant (r,c)
func (g *Game) checkWin(r, c int) bool {
	player := g.Board[r][c]
	if player == 0 || player == 3 {
		return false
	}
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}

	for _, d := range directions {
		line := make([][2]int, 0, 8)
		// aller vers la direction négative
		r0, c0 := r, c
		for {
			r0 -= d[0]
			c0 -= d[1]
			if r0 < 0 || r0 >= Rows || c0 < 0 || c0 >= Columns {
				break
			}
			if g.Board[r0][c0] != player {
				break
			}
			line = append(line, [2]int{r0, c0})
		}
		// inverser le préfixe pour garder l'ordre
		for i := 0; i < len(line)/2; i++ {
			line[i], line[len(line)-1-i] = line[len(line)-1-i], line[i]
		}
		// ajouter la cellule courante
		line = append(line, [2]int{r, c})
		// aller vers la direction positive
		r1, c1 := r, c
		for {
			r1 += d[0]
			c1 += d[1]
			if r1 < 0 || r1 >= Rows || c1 < 0 || c1 >= Columns {
				break
			}
			if g.Board[r1][c1] != player {
				break
			}
			line = append(line, [2]int{r1, c1})
		}

		if len(line) >= 4 {
			// trouver position du jeton joué
			idx := -1
			for i, cell := range line {
				if cell[0] == r && cell[1] == c {
					idx = i
					break
				}
			}
			if idx == -1 {
				continue
			}
			start := idx - 3
			if start < 0 {
				start = 0
			}
			if start+4 > len(line) {
				start = len(line) - 4
			}
			w := make([][2]int, 0, 4)
			for i := start; i < start+4; i++ {
				w = append(w, line[i])
			}
			g.WinningCells = w
			return true
		}
	}
	return false
}
