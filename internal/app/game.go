package app

import (
	"fmt"
	"strings"
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
	t := time.NewTicker(time.Second * 2)

	for {
		<-t.C

		g.PrintPlot()

		// TICK TIME LOCK
		g.mu.Lock()

		moves := DeduplicateMoves(g.movesQueue)
		g.applyMoves(moves)

		// clear moves queue
		g.movesQueue = g.movesQueue[:0]

		for i, callback := range g.addQueue {
			id, ok := g.CreatePlayer()
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
