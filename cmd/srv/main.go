package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hurtki/ascii-snake/internal/app"
	"github.com/hurtki/ascii-snake/internal/domain"
	http_handlers "github.com/hurtki/ascii-snake/internal/handlers/http"
	"github.com/hurtki/ascii-snake/internal/ws"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelDebug}))
	game := app.InitGame(20)
	go game.Start(time.Second / 3)

	var wsHandler ws.Server
	usecase := domain.NewGameUsecase(game, &wsHandler)
	wsHandler = *ws.NewServer(usecase, logger)

	joinHandler := http_handlers.NewJoinHandler(usecase)

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("Upgrade error: %v\n", err)
			return // Ответ со статусом ошибки отправится автоматически
		}
		fmt.Println("token got:", token)

		wsHandler.HandleWS(conn, token)
	})

	http.HandleFunc("GET /room", joinHandler.Join)

	http.ListenAndServe(":80", nil)
}
