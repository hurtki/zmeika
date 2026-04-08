package app

type cord struct {
	x int
	y int
}

func (g *Game) applyMoves(moves []Move) {
	// seen is used to fastly find players on big field
	headsSeen := make(map[int]cord)

	for x, row := range g.plot {
		for y, cell := range row {
			if cell.Value != 0 && cell.IsHead {
				headsSeen[cell.PlayerID] = cord{x: x, y: y}
			}
			v := g.plot[x][y]
			// decrement all the values
			if v.Value > 0 {
				if v.Value == 1 {
					g.plot[x][y].IsHead = false
					g.plot[x][y].PlayerID = 0
				}
				g.plot[x][y].Value--
			}
		}
	}

	for _, m := range moves {
		if headCord, ok := headsSeen[m.PlayerID]; ok {
			ok := g.applyOneMoveWithHeadCord(m, headCord)
			if ok {
				delete(headsSeen, m.PlayerID)
				continue
			}
		}
	}

	// if there was no move, we also need to move it
	// going through every seen head, that wasn't already deleted
	for id, c := range headsSeen {
		m := Move{PlayerID: id}
		switch {
		// prevous was upper => go down
		case InGaps(cord{c.x - 1, c.y}, len(g.plot)) && g.plot[c.x-1][c.y].PlayerID == id:
			m.Direction = Down
		// previous was lower => go up
		case InGaps(cord{c.x + 1, c.y}, len(g.plot)) && g.plot[c.x+1][c.y].PlayerID == id:
			m.Direction = Up
		// previous was at left to head => go right
		case InGaps(cord{c.x, c.y - 1}, len(g.plot)) && g.plot[c.x][c.y-1].PlayerID == id:
			m.Direction = Right
		// previous was at right to head => go left
		case InGaps(cord{c.x, c.y + 1}, len(g.plot)) && g.plot[c.x][c.y+1].PlayerID == id:
			m.Direction = Left
		}
		g.applyOneMoveWithHeadCord(m, c)
	}
}

// aplies one move, when coordintares of player's head that
// did the move are known
func (g *Game) applyOneMoveWithHeadCord(move Move, c cord) (ok bool) {
	ok = true
	moveCord := cord{}
	switch move.Direction {
	case Up:
		moveCord = cord{x: c.x - 1, y: c.y}
	case Down:
		moveCord = cord{x: c.x + 1, y: c.y}
	case Left:
		moveCord = cord{x: c.x, y: c.y - 1}
	case Right:
		moveCord = cord{x: c.x, y: c.y + 1}
	}

	startCord := c // coordinates before tick

	if !InGaps(moveCord, len(g.plot)) {
		g.removePlayer(startCord, -1)
		return
	}

	startCell := g.plot[c.x][c.y]              // cell before tick
	moveCell := g.plot[moveCord.x][moveCord.y] // cell that move goes to

	if moveCell.PlayerID == startCell.PlayerID {
		if moveCell.Value+1 != startCell.Value {
			g.removePlayer(startCord, -1)
			return
		}
		// not "OK" case, when move goes backward
		return false
	}

	if moveCell.PlayerID > 0 && moveCell.Value > 0 {
		// if move is into someones head => kill both
		if moveCell.IsHead {
			g.removePlayer(moveCord, -1)
		}
		g.removePlayer(startCord, -1)
		return
	}

	g.plot[moveCord.x][moveCord.y].PlayerID = startCell.PlayerID
	g.plot[moveCord.x][moveCord.y].IsHead = true
	// should have value one bigger, then previous head cell
	// cause all the cells were decremented, before applying moves
	g.plot[moveCord.x][moveCord.y].Value = startCell.Value + 1

	// not forgetting to change previous head
	g.plot[startCord.x][startCord.y].IsHead = false
	return
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

	if v.Value == 0 || v.PlayerID != playerID {
		return
	}

	g.plot[c.x][c.y] = Cell{}
	g.removePlayer(cord{x: c.x - 1, y: c.y}, playerID)
	g.removePlayer(cord{x: c.x + 1, y: c.y}, playerID)
	g.removePlayer(cord{x: c.x, y: c.y + 1}, playerID)
	g.removePlayer(cord{x: c.x, y: c.y - 1}, playerID)
}
