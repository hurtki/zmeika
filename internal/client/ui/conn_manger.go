package ui

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/gorilla/websocket"
)

type RoomDTO struct {
	Token    string `json:"token"`
	PlayerID int    `json:"player_id"`
	MapSize  int    `json:"map_size"`
}

type ConnectionManager struct {
	wsConn   *websocket.Conn
	playerID int
	mapSize  int
}

// tea messages

type FrameMsg []byte
type WsErrMsg struct{ Err error }

func NewConnectionManager(addr string) (*ConnectionManager, error) {
	res, err := http.Get("http://" + addr + "/room")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var dto RoomDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		return nil, fmt.Errorf("invalid server response: %s", string(body))
	}

	wsURL := fmt.Sprintf("ws://%s/ws?token=%s", addr, dto.Token)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("can't establish websocket connection: %w", err)
	}

	return &ConnectionManager{
		wsConn:   conn,
		playerID: dto.PlayerID,
		mapSize:  dto.MapSize,
	}, nil
}

// Listen starts reading WebSocket frames and returns a bubbletea Cmd.
// Each frame is dispatched as FrameMsg; errors as WsErrMsg.
func (cm *ConnectionManager) Listen() tea.Cmd {
	return func() tea.Msg {
		_, msg, err := cm.wsConn.ReadMessage()
		if err != nil {
			return WsErrMsg{Err: err}
		}
		return FrameMsg(msg)
	}
}

// SendMove sends a direction byte over the WebSocket.
func (cm *ConnectionManager) SendMove(dir byte) tea.Cmd {
	return func() tea.Msg {
		if err := cm.wsConn.WriteMessage(websocket.BinaryMessage, []byte{dir}); err != nil {
			return WsErrMsg{Err: err}
		}
		return nil
	}
}

// Close closes the underlying WebSocket connection.
func (cm *ConnectionManager) Close() {
	cm.wsConn.Close()
}

// MapMove converts a key string to a protocol direction byte.
// Returns 0xFF for unknown keys.
func MapMove(key string) byte {
	switch key {
	case "w", "k":
		return 0x0
	case "s", "j":
		return 0x1
	case "a", "h":
		return 0x2
	case "d", "l":
		return 0x3
	default:
		return 0xFF
	}
}

// RenderPlot converts the raw binary plot frame into a string
// ready to be displayed by lipgloss / bubbletea.
func RenderPlot(plot []byte, myID int, size int) string {
	const cellSize = 4
	var b strings.Builder

	for i := 0; i+cellSize <= len(plot); i += cellSize {
		playerID := binary.LittleEndian.Uint16(plot[i : i+2])
		value := plot[i+2]
		isHead := plot[i+3] == 0x1

		symbol := ". "
		if value > 0 {
			var char string
			if isHead {
				char = "H "
			} else {
				char = fmt.Sprintf("%d ", value%10)
			}

			if int(playerID) == myID {
				symbol = "\033[32m" + char + "\033[0m"
			} else {
				symbol = char
			}
		}

		b.WriteString(symbol)

		if cellIndex := i / cellSize; (cellIndex+1)%size == 0 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
