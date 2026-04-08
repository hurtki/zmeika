package app

import "fmt"

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
			v := g.plot[x][y]
			// DECREMENT ALL THE VALUES
			if v.Value > 0 {
				if v.Value == 1 {
					g.plot[x][y].IsHead = false
					g.plot[x][y].PlayerID = 0
				}
				g.plot[x][y].Value--
			}
		}
	}
	fmt.Println("print after decrementing all")
	g.PrintPlot()

	for _, m := range moves {
		if headCord, ok := headsSeen[m.PlayerID]; ok {
			g.applyOneMoveWithHeadCord(m, headCord)
			delete(headsSeen, m.PlayerID)
			continue
		} else {
			continue
		}
	}

	// if there was no move, we also need to move it
	// going through every seen head, that wasn't already deleted
	for id, c := range headsSeen {
		m := Move{PlayerID: id}
		switch {
		// prevous was upper => go down
		case InGaps(cord{c.x - 1, c.y}, len(g.plot)) && g.plot[c.x-1][c.y].PlayerID == id:
			m.Direction = down
		// previous was lower => go up
		case InGaps(cord{c.x + 1, c.y}, len(g.plot)) && g.plot[c.x+1][c.y].PlayerID == id:
			m.Direction = up
		// previous was at left to head => go right
		case InGaps(cord{c.x, c.y - 1}, len(g.plot)) && g.plot[c.x][c.y-1].PlayerID == id:
			m.Direction = right
		// previous was at right to head => go left
		case InGaps(cord{c.x, c.y + 1}, len(g.plot)) && g.plot[c.x][c.y+1].PlayerID == id:
			m.Direction = left
		}
		fmt.Println("after checking direction decided: direction: ", m.Direction, "applying one move from", c)
		g.applyOneMoveWithHeadCord(m, c)
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
		moveCord = cord{x: c.x, y: c.y + 1}
	}

	if !InGaps(moveCord, len(g.plot)) {
		fmt.Println("not in gaps, removing cause rached border:", c)
		g.removePlayer(c, -1)
		return
	}

	startCell := g.plot[c.x][c.y]
	moveCell := g.plot[moveCord.x][moveCord.y]

	if moveCell.PlayerID == startCell.PlayerID {
		return
	}

	if moveCell.PlayerID > 0 && moveCell.Value > 0 {
		fmt.Println("collision")
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
	fmt.Println("updated after move to free cell", g.plot[moveCord.x][moveCord.y])

	fmt.Println("removed head from prevous on", c.x, c.y)
	// not forgetting to change previous head
	g.plot[c.x][c.y].IsHead = false
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
