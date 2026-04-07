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
