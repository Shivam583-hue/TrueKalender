package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var models []tea.Model

func SetModels(m []tea.Model) {
	models = m
}

const (
	mainModel = 0
	formModel = 1
)

type Form struct {
	focused int
	title   textinput.Model
}

func NewForm(focused int) Form {
	ti := textinput.New()
	ti.Placeholder = "Enter title"
	ti.Focus()
	return Form{focused: focused, title: ti}
}

func (f Form) Init() tea.Cmd { return nil }

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return f, tea.Quit
		case "esc":
			return models[mainModel], nil
		case "enter":
			// return models[mainModel], f.createTask
			return models[mainModel], nil
		}
	}
	f.title, cmd = f.title.Update(msg)
	return f, cmd
}

// func (f Form) createTask() tea.Msg {
// 	// db.Insert(f.title.Value(), f.focused)
// 	// return types.Task{
// 	// 	TaskTitle: f.title.Value(),
// 	// 	Status:    f.focused,
// 	// }
// 	return
// }

func (f Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Create New Task",
		"",
		f.title.View(),
		"",
		"Enter to save • Esc to cancel",
	)
}
