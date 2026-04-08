package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/hurtki/ascii-snake/internal/client/ui"
)

func main() {
	p := tea.NewProgram(ui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("can't run", err)
		return
	}
}
