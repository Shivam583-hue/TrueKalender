package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var models []tea.Model

var (
	purple       = lipgloss.Color("135")
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	cellStyle = lipgloss.NewStyle().
			Width(6).
			Height(3).
			Border(lipgloss.NormalBorder()).
			BorderForeground(purple).
			Foreground(purple).
			Align(lipgloss.Center, lipgloss.Center)

	headerStyle = lipgloss.NewStyle().
			Width(8).
			Height(1).
			Bold(true).
			Foreground(purple).
			Align(lipgloss.Center)
)

var dayHeaders = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

type model struct {
	quitting bool
	help     help.Model
	lists    []list.Model
	loaded   bool
	calendar []string
}

func New() *model {
	return &model{help: help.New()}
}

func SetModels(m []tea.Model) {
	models = m
}

func (m *model) initCalendar(width, height int) {
	headers := make([]string, 7)
	for i, d := range dayHeaders {
		headers[i] = headerStyle.Render(d)
	}
	m.calendar = append(m.calendar, lipgloss.JoinHorizontal(lipgloss.Top, headers...))

	var cells []string
	for day := 1; day <= 31; day++ {
		cell := cellStyle.Render(fmt.Sprintf("%d", day))
		cells = append(cells, cell)

		if day%7 == 0 {
			m.calendar = append(m.calendar, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
			cells = nil
		}
	}
	if len(cells) > 0 {
		for len(cells) < 7 {
			cells = append(cells, cellStyle.Render(""))
		}
		m.calendar = append(m.calendar, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initCalendar(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if !m.loaded {
		return "Loading..."
	}
	return lipgloss.JoinVertical(lipgloss.Left, m.calendar...)
}

func main() {
	mainModel := New()
	p := tea.NewProgram(mainModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
