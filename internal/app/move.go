package app

type Direction uint8

const (
	up Direction = iota
	down
	left
	right
)

type Move struct {
	PlayerID  int
	Direction Direction
}
