package game

const (
	Rows    = 6
	Columns = 7
)

type Power struct {
	Name    string
	Uses    int
	MaxUses int
}

type Game struct {
	Board        [Rows][Columns]int
	Current      int
	Winner       int
	LastRow      int
	LastCol      int
	WinningCells [][2]int
	Powers       map[int][]Power
}

func NewGame() *Game {
	return &Game{
		Current: 1,
		LastRow: -1,
		LastCol: -1,
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
				g.WinningCells = nil
				g.Current = 3 - g.Current
			}
			return
		}
	}
}

func (g *Game) UsePower(player int, powerName string, row, col int) bool {
	if g.Winner != 0 || player != g.Current {
		return false
	}
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
		g.Powers[player][powerIndex].Uses--
		g.WinningCells = nil
		g.Current = 3 - g.Current
	}
	return success
}

func (g *Game) destroyPiece(row, col int) bool {
	if row < 0 || row >= Rows || col < 0 || col >= Columns {
		return false
	}
	val := g.Board[row][col]
	if val == 0 || val == 3 || val == g.Current {
		return false
	}
	g.Board[row][col] = 0
	g.applyGravity()
	return true
}

func (g *Game) swapPieces(row, col int) bool {
	if row < 0 || row >= Rows-1 || col < 0 || col >= Columns {
		return false
	}
	v1 := g.Board[row][col]
	v2 := g.Board[row+1][col]
	if v1 == 0 || v2 == 0 || v1 == 3 || v2 == 3 {
		return false
	}
	g.Board[row][col], g.Board[row+1][col] = v2, v1
	return true
}

func (g *Game) blockColumn(col int) bool {
	if col < 0 || col >= Columns {
		return false
	}
	for row := Rows - 1; row >= 0; row-- {
		if g.Board[row][col] == 0 {
			g.Board[row][col] = 3
			return true
		}
		if g.Board[row][col] == 3 {
			return false
		}
	}
	return false
}

func (g *Game) applyGravity() {
	for col := 0; col < Columns; col++ {
		writePos := Rows - 1
		for row := Rows - 1; row >= 0; row-- {
			val := g.Board[row][col]
			if val == 3 {
				writePos = row - 1
				continue
			}
			if val != 0 {
				if row != writePos {
					g.Board[writePos][col] = val
					g.Board[row][col] = 0
				}
				writePos--
			}
		}
		for r := writePos; r >= 0; r-- {
			if g.Board[r][col] != 3 {
				g.Board[r][col] = 0
			}
		}
	}
}

func (g *Game) checkWin(r, c int) bool {
	player := g.Board[r][c]
	if player == 0 || player == 3 {
		return false
	}
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}

	for _, d := range directions {
		line := make([][2]int, 0, 8)
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
		for i := 0; i < len(line)/2; i++ {
			line[i], line[len(line)-1-i] = line[len(line)-1-i], line[i]
		}
		line = append(line, [2]int{r, c})
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
