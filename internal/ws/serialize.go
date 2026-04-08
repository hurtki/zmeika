package ws

import (
	"encoding/binary"

	"github.com/hurtki/zmeika/internal/app"
)

func SerializePlot(plot [][]app.Cell) []byte {
	res := make([]byte, len(plot)*len(plot)*4)
	for _, row := range plot {
		for _, cell := range row {
			binary.LittleEndian.PutUint16(res, uint16(cell.PlayerID))
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
