package ws

import (
	"encoding/binary"

	"github.com/hurtki/zmeika/internal/app"
)

// Serializes whole plot into binary format
// One cell: [[2b player id][1 byte value][1 byte isHead]]
func SerializePlot(plot [][]app.Cell) []byte {
	res := make([]byte, 0, len(plot)*len(plot)*4)
	for _, row := range plot {
		for _, cell := range row {
			buf := make([]byte, 2)
			binary.LittleEndian.PutUint16(buf, uint16(cell.PlayerID))
			res = append(res, buf...)
			res = append(res, byte(uint8(cell.Value)))
			if cell.IsHead {
				res = append(res, 1)
			} else {
				res = append(res, 0)
			}
		}
	}
	return res
}
