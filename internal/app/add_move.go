package app

// adds move to queue that will be applied in tick time
// can't return error, but there is not guarantee, that move will be applied
// *for example if there were two, only first one will be applied
// not for tick time!
func (g *Game) AddMove(move Move) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.addQueueMu.Lock()
	g.movesQueue = append(g.movesQueue, move)
	g.addQueueMu.Unlock()
}
