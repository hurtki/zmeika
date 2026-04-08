package app

import "errors"

type Cell struct {
	// if there is not player here => 0
	PlayerID int // uint16
	// if this cell is head of the snake
	IsHead bool

	// if 0 then there is nothing there
	// if cell is head => same as snakeLength
	// "depth" of snake [head[3] <- 2 <- 1 ]
	Value int
}

// Plot is a field of cells
// Coordinates here are indecies p[x][y] is (x, y) point
type Plot [][]Cell

// Direction is used to determine a move on the Plot
type Direction uint8

const (
	Up Direction = iota
	Down
	Left
	Right
)

func NewDirection(d uint8) (Direction, error) {
	if d > 3 {
		return Direction(0), errors.New("not existing direction")
	}
	return Direction(d), nil
}

// Abstract move for game input
type Move struct {
	PlayerID  int
	Direction Direction
}
