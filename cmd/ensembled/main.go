package main

import (
	"ensemble/internal/tui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.CreateTui(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("an error!: %v", err)
		os.Exit(1)
	}
}
