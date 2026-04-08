package app

import (
	"sync"
	"time"
)

var snakeLength = 6

type Game struct {
	plot Plot

	// ID for new plater
	// increment after adding a new one
	// starts with 1, cause 0 is zero value
	cntr int

	// add a new player requests
	// slice of callback functions
	addQueue   []func(int, error)
	addQueueMu sync.Mutex

	// apply move queue
	// slice of moves that came to server
	movesQueue   []Move
	movesQueueMu sync.Mutex

	// used on tick "stop the world"
	mu sync.RWMutex

	// After every tick one struct is sent to notify readers to check plot
	AfterTickCh chan struct{}
}

func InitGame(size int) *Game {
	plot := make([][]Cell, size)
	for i := range size {
		plot[i] = make([]Cell, size)
	}
	return &Game{
		plot:        plot,
		cntr:        1,
		AfterTickCh: make(chan struct{}, 10000),
	}
}

func (g *Game) Start(tickTime time.Duration) {
	t := time.NewTicker(tickTime)

	for {
		<-t.C

		// TICK TIME LOCK
		g.mu.Lock()

		moves := DeduplicateMoves(g.movesQueue)
		g.applyMoves(moves)

		// clear moves queue
		g.movesQueue = g.movesQueue[:0]

		for i, callback := range g.addQueue {
			id, ok := g.createPlayer()
			if !ok {
				for j := i; i < len(g.addQueue); i++ {
					g.addQueue[j](0, ErrNoPlaceOnPlot)
				}
				break
			}
			callback(id, nil)
		}
		g.addQueue = g.addQueue[:0]

		g.mu.Unlock()
	}
}

func (g *Game) GetTickMap() [][]Cell {
	<-g.AfterTickCh
	g.mu.RLock()
	plot := g.plot
	g.mu.RUnlock()
	return plot
}
