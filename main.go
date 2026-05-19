package main

import (
	"fmt"
	"os"

	"github.com/Shivam583-hue/TrueKalender/model"
	tea "github.com/charmbracelet/bubbletea"
)

var models []tea.Model

func SetModels(m []tea.Model) {
	models = m
}

func main() {
	mainModel := model.New()
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
