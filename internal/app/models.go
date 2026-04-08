package app

type Cell struct {
	// if there is not player here => 0
	PlayerID int
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
	up Direction = iota
	down
	left
	right
)

// Abstract move for game input
type Move struct {
	PlayerID  int
	Direction Direction
}
