package app

func (g *Game) CreatePlayer() (int, bool) {
	rows := len(g.plot)
	if rows == 0 {
		return 0, false
	}
	cols := len(g.plot[0])

	isFree := func(x, y int) bool {
		if g.plot[x][y].Value != 0 {
			return false
		}

		// zone of 1 cell around
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				nx, ny := x+dx, y+dy
				if nx < 0 || ny < 0 || nx >= rows || ny >= cols {
					continue
				}
				if g.plot[nx][ny].Value != 0 {
					return false
				}
			}
		}

		return true
	}

	checkHorizontal := func(x, y int) bool {
		if y+snakeLength > cols {
			return false
		}

		for i := range snakeLength {
			if !isFree(x, y+i) {
				return false
			}
		}
		return true
	}

	checkVertical := func(x, y int) bool {
		if x+snakeLength > rows {
			return false
		}

		for i := range snakeLength {
			if !isFree(x+i, y) {
				return false
			}
		}
		return true
	}

	place := func(id int, x, y int, dx, dy int) {
		for i := range snakeLength {
			cx := x + dx*i
			cy := y + dy*i

			g.plot[cx][cy].PlayerID = id
			g.plot[cx][cy].Value = i + 1
			g.plot[cx][cy].IsHead = false
		}

		hx := x + dx*(snakeLength-1)
		hy := y + dy*(snakeLength-1)
		g.plot[hx][hy].IsHead = true
		g.plot[hx][hy].Value = snakeLength
	}

	id := g.cntr + 1
	g.cntr = id

	for i := range rows {
		for j := range cols {

			// horizontal
			if checkHorizontal(i, j) {
				place(id, i, j, 0, 1)
				return id, true
			}

			// vertical
			if checkVertical(i, j) {
				place(id, i, j, 1, 0)
				return id, true
			}
		}
	}

	return 0, false
}
