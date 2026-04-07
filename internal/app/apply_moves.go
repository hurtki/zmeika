package app

type cord struct {
	x int
	y int
}

func (g *Game) applyMoves(moves []Move) {

	// 0:0 is left top pixel

	// seen is used to fastly find players on big field
	headsSeen := make(map[int]cord)

	for x, row := range g.plot {
		for y, cell := range row {
			if cell.Value != 0 && cell.IsHead {
				headsSeen[cell.PlayerID] = cord{x: x, y: y}
			}
			// DECREMENT ALL THE VALUES
			g.plot[x][y].Value -= 1
		}
	}

	for _, m := range moves {
		if cord, ok := headsSeen[m.PlayerID]; ok {
			g.applyOneMoveWithHeadCord(m, cord)
			continue
		}
	}
}

// aplies one move, when coordintares of player's head that
// did the move are known
func (g *Game) applyOneMoveWithHeadCord(move Move, c cord) {
	moveCord := cord{}
	switch move.Direction {
	case up:
		moveCord = cord{x: c.x - 1, y: c.y}
	case down:
		moveCord = cord{x: c.x + 1, y: c.y}
	case left:
		moveCord = cord{x: c.x, y: c.y - 1}
	case right:
		moveCord = cord{x: c.x - 1, y: c.y + 1}
	}

	if !InGaps(moveCord, len(g.plot)) {
		g.removePlayer(c, -1)
		return
	}

	startCell := g.plot[c.x][c.y]
	moveCell := g.plot[moveCord.x][moveCord.y]

	if moveCell.PlayerID == startCell.PlayerID {
		return
	}

	if moveCell.PlayerID > 0 && moveCell.Value > 0 {
		// if move is into someones head => kill both
		if moveCell.IsHead {
			g.removePlayer(moveCord, -1)
		}
		g.removePlayer(c, -1)
		return
	}

	g.plot[moveCord.x][moveCord.y].PlayerID = startCell.PlayerID
	g.plot[moveCord.x][moveCord.y].IsHead = true
	// should have value one bigger, then previous head cell
	// cause all the cells were decremented, before applying moves
	g.plot[moveCord.x][moveCord.y].Value = startCell.Value + 1
}

// Removes all the nearby cells with playerID
// if playerID -1, init playerID will be ID that stays on cord
func (g *Game) removePlayer(c cord, playerID int) {
	if !InGaps(c, len(g.plot)) {
		return
	}

	v := g.plot[c.x][c.y]

	if playerID == -1 {
		playerID = v.PlayerID
	}

	if v.Value == 0 {
		return
	}

	if v.PlayerID == playerID {
		g.plot[c.x][c.y] = Cell{}
	}
	g.removePlayer(cord{x: c.x - 1, y: c.y}, playerID)
	g.removePlayer(cord{x: c.x - 1, y: c.y + 1}, playerID)
	g.removePlayer(cord{x: c.x, y: c.y + 1}, playerID)
	g.removePlayer(cord{x: c.x, y: c.y}, playerID)
}
