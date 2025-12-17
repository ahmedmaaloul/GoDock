package main

import (
	"fmt"
	"os"

	dockerwrapper "github.com/ahmedmaaloul/godock-tui-manager/docker"
	"github.com/ahmedmaaloul/godock-tui-manager/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize Docker Client
	dc, err := dockerwrapper.NewClient()
	if err != nil {
		fmt.Printf("Error initializing Docker client: %v\n", err)
		os.Exit(1)
	}

	// Initialize UI Model
	m := ui.NewModel(dc)

	// Start Bubble Tea Program
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
