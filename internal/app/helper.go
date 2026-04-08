package app

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
func (g *Game) GetMapSize() int {
	return len(g.plot)
}
