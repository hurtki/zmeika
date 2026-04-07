package app

import "errors"

var (
	ErrNoPlaceOnPlot = errors.New("there is no place for snake on game plot")
)

// AddPlayer blocks until new tick decides how to create player
// ErrNoPlaceOnPlot is returned if there is no place
// not for tick time!
func (g *Game) AddPlayer() (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	ch := make(chan struct {
		id  int
		err error
	}, 1)

	g.addQueueMu.Lock()
	g.addQueue = append(g.addQueue, func(playerID int, err error) {
		ch <- struct {
			id  int
			err error
		}{playerID, err}
	})
	g.addQueueMu.Unlock()

	res := <-ch
	return res.id, res.err
}
