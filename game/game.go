package game

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
		{0, 1}, {1, 0}, {1, 1}, {1, -1},
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
