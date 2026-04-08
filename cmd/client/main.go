package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/term"
)

type RoomDTO struct {
	Token    string `json:"token"`
	PlayerID int    `json:"player_id"`
	MapSize  int    `json:"map_size"`
}

func main() {
	var addr string
	fmt.Print("Enter remote addr (e.g. localhost:8080) -> ")
	fmt.Scanln(&addr)

	dto, err := fetchRoomData(addr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	wsURL := fmt.Sprintf("ws://%s/ws?token=%s", addr, dto.Token)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected! Use WASD to move.")

	// raw mode: read char by char without Enter
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Failed to set raw mode: %v\n", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	frames := make(chan []byte)
	errors := make(chan error, 2)

	// goroutine: read WebSocket frames
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				errors <- err
				return
			}
			frames <- message
		}
	}()

	// goroutine: read keyboard input
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				continue
			}
			key := buf[0]
			if key == 3 || key == 'q' {
				errors <- fmt.Errorf("user quit")
				return
			}
			move := mapMove(key)
			if move == 0xFF {
				continue
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, []byte{move}); err != nil {
				errors <- err
				return
			}
		}
	}()

	// save cursor position before first render
	os.Stdout.WriteString("\033[s")

	for {
		select {
		case frame := <-frames:
			// restore cursor position and clear below
			os.Stdout.WriteString("\033[u\033[J")
			PrintPlot(frame, dto.PlayerID, dto.MapSize)
		case err := <-errors:
			term.Restore(int(os.Stdin.Fd()), oldState)
			os.Stdout.WriteString("\r\nDisconnected: " + err.Error() + "\r\n")
			return
		}
	}
}

func fetchRoomData(addr string) (*RoomDTO, error) {
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
	return &dto, nil
}

func mapMove(key byte) byte {
	switch key {
	case 'w':
		return 0x0
	case 's':
		return 0x1
	case 'a':
		return 0x2
	case 'd':
		return 0x3
	default:
		return 0xFF
	}
}

func PrintPlot(plot []byte, myID int, size int) {
	var b strings.Builder
	const cellSize = 4

	for i := 0; i < len(plot); i += cellSize {
		if i+cellSize > len(plot) {
			break
		}
		playerID := binary.LittleEndian.Uint16(plot[i : i+2])
		value := plot[i+2]
		isHead := plot[i+3] == 0x1

		symbol := ". "
		if value > 0 {
			color := ""
			reset := ""
			if int(playerID) == myID {
				color = "\033[32m"
				reset = "\033[0m"
			}
			char := fmt.Sprintf("%d ", value%10)
			if isHead {
				char = "H "
			}
			symbol = color + char + reset
		}

		b.WriteString(symbol)

		cellIndex := i / cellSize
		if (cellIndex+1)%size == 0 {
			b.WriteString("\r\n")
		}
	}

	os.Stdout.WriteString(b.String())
}
