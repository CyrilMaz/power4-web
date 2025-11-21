package game

const (
	Rows    = 6
	Columns = 7
)

type Power struct {
	Name     string
	Uses     int
	MaxUses  int
	Cooldown int
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
	}

	if success {
		g.Powers[player][powerIndex].Uses--
		g.Current = 3 - g.Current
	}

	return success
}

func (g *Game) destroyPiece(row, col int) bool {
	if row < 0 || row >= Rows || col < 0 || col >= Columns {
		return false
	}
	if g.Board[row][col] == 0 || g.Board[row][col] == g.Current {
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
	if g.Board[row][col] == 0 || g.Board[row+1][col] == 0 {
		return false
	}
	g.Board[row][col], g.Board[row+1][col] = g.Board[row+1][col], g.Board[row][col]
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
	}
	return false
}

func (g *Game) applyGravity() {
	for col := 0; col < Columns; col++ {
		writePos := Rows - 1
		for row := Rows - 1; row >= 0; row-- {
			if g.Board[row][col] != 0 {
				if row != writePos {
					g.Board[writePos][col] = g.Board[row][col]
					g.Board[row][col] = 0
				}
				writePos--
			}
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
