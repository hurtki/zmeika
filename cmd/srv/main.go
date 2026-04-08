package main

import (
	"net/http"
	"time"

	"github.com/hurtki/zmeika/internal/app"
	"github.com/hurtki/zmeika/internal/domain"
	http_handlers "github.com/hurtki/zmeika/internal/handlers/http"
	"github.com/hurtki/zmeika/internal/ws"
	"golang.org/x/net/websocket"
)

func main() {
	game := app.InitGame(10)
	go game.Start(time.Second * 2)

	var wsHandler ws.Server
	usecase := domain.NewGameUsecase(game, &wsHandler)
	wsHandler = *ws.NewServer(usecase)

	joinHandler := http_handlers.NewJoinHandler(usecase)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		websocket.Handler(func(conn *websocket.Conn) {
			wsHandler.HandleWS(conn, token)
		}).ServeHTTP(w, r)
	})

	http.HandleFunc("/room", joinHandler.Join)

	http.ListenAndServe(":80", nil)
}
