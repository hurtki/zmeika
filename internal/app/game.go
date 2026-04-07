package app

import (
	"sync"
	"time"
)

var snakeLength = 3

type Cell struct {
	Value int // if cell is head => same as snakeLength
	// then it goes down
	// if 0 then there is nothing there

	PlayerID int
	IsHead   bool
}

type Game struct {
	plot [][]Cell

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
}

func InitGame(size int) *Game {
	plot := make([][]Cell, size)
	for i := range size {
		plot[i] = make([]Cell, size)
	}
	return &Game{
		plot: plot,
	}
}

func (g *Game) Start() {
	t := time.NewTicker(time.Second)

	for {
		<-t.C
		// TICK TIME LOCK
		g.mu.Lock()

		moves := DeduplicateMoves(g.movesQueue)
		g.applyMoves(moves)

		g.mu.Unlock()

	}
}
