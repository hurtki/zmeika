package app

import (
	"fmt"
	"strings"
)

func DeduplicateMoves(moves []Move) []Move {
	seen := make(map[int]struct{}, len(moves))
	res := make([]Move, 0, len(moves))

	for _, m := range moves {
		_, ok := seen[m.PlayerID]
		if !ok {
			res = append(res, m)
			seen[m.PlayerID] = struct{}{}
		}
	}
	return res
}

func InGaps(c cord, size int) bool {
	if c.x < 0 || c.x >= size {
		return false
	}
	if c.y < 0 || c.y >= size {
		return false
	}
	return true
}

func (g *Game) PrintPlot() {
	var b strings.Builder

	for i := 0; i < len(g.plot); i++ {
		for j := 0; j < len(g.plot[i]); j++ {
			c := g.plot[i][j]

			switch {
			case c.Value == 0:
				b.WriteString(". ")
			case c.IsHead:
				b.WriteString("H ")
			default:
				// тело змейки
				fmt.Fprintf(&b, "%d ", c.Value)
			}
		}
		b.WriteString("\n")
	}

	fmt.Print(b.String())
}
