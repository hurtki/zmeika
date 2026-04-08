package main

import (
	"fmt"
	"time"

	"github.com/hurtki/zmeika/internal/app"
)

func main() {
	game := app.InitGame(10)
	go game.Start()
	id, err := game.AddPlayer()
	fmt.Println("id:", id, "err", err)

	time.Sleep(time.Second * 100)
}
