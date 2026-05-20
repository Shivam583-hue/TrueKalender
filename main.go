package main

import (
	"fmt"
	"os"

	"github.com/Shivam583-hue/TrueKalender/model"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	mainModel := model.New()
	models := []tea.Model{mainModel, model.NewForm(0)}
	model.SetModels(models)
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
