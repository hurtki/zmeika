package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/hurtki/zmeika/internal/app"
)

func main() {
	game := app.InitGame(25)
	go game.Start()

	for range 10 {
		go game.AddPlayer()
		go func() {
			dirs := []app.Direction{app.Down, app.Left, app.Up, app.Right}
			for _, d := range dirs {
				time.Sleep(time.Second * 4)
				game.AddMove(app.Move{PlayerID: 10, Direction: d})
			}
		}()
	}

	time.Sleep(time.Second / 2)
	go func() {
		for {
			game.PrintPlot()
			time.Sleep(time.Second)
			deleteCoupleLines(25)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	// var input string
	// fmt.Printf("->")
	// fmt.Scanln(&input)
}

func deleteCoupleLines(n int) {
	for range n {
		fmt.Print("\033[F\033[K")
	}
}
