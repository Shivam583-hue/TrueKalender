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
	purple = lipgloss.Color("135")

	focusedCellStyle = lipgloss.NewStyle().
				Width(6).
				Height(3).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("255")).
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("55")).
				Bold(true).
				Align(lipgloss.Center, lipgloss.Center)

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

	panelStyle = lipgloss.NewStyle().
			Width(28).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple)

	panelTitleStyle = lipgloss.NewStyle().
			Background(purple).
			Foreground(lipgloss.Color("255")).
			Bold(true).
			Padding(0, 1)

	taskStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			PaddingLeft(1)

	taskDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(1)

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			PaddingLeft(1)
)

var dayHeaders = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// mock data
var tasks = map[int][]string{
	1:  {"read bittorrent code", "complete math"},
	3:  {"revise graphs"},
	7:  {"start truekalender"},
	15: {"dentist appt", "grocery run", "call mom"},
}

type model struct {
	quitting bool
	help     help.Model
	lists    []list.Model
	focused  int
	loaded   bool
}

func New() *model {
	return &model{
		help:    help.New(),
		focused: 1,
	}
}

func SetModels(m []tea.Model) {
	models = m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "l", "right":
			m.focused = min(m.focused+1, 31)
		case "h", "left":
			m.focused = max(m.focused-1, 1)
		case "j", "down":
			m.focused = min(m.focused+7, 31)
		case "k", "up":
			m.focused = max(m.focused-7, 1)
		}
	}
	return m, nil
}

func (m model) renderPanel() string {
	title := panelTitleStyle.Render(fmt.Sprintf("Day %d", m.focused))

	dayTasks, ok := tasks[m.focused]

	var body string
	if !ok || len(dayTasks) == 0 {
		body = emptyStyle.Render("no tasks")
	} else {
		for _, t := range dayTasks {
			body += taskStyle.Render("• "+t) + "\n"
			body += taskDescStyle.Render("Todo") + "\n\n"
		}
	}

	return panelStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, title, "", body),
	)
}

func (m model) View() string {
	if !m.loaded {
		return "Loading..."
	}
	if m.quitting {
		return ""
	}

	headers := make([]string, 7)
	for i, d := range dayHeaders {
		headers[i] = headerStyle.Render(d)
	}
	rows := []string{lipgloss.JoinHorizontal(lipgloss.Top, headers...)}

	var cells []string
	for day := 1; day <= 31; day++ {
		style := cellStyle
		if day == m.focused {
			style = focusedCellStyle
		}
		cells = append(cells, style.Render(fmt.Sprintf("%d", day)))
		if day%7 == 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
			cells = nil
		}
	}
	if len(cells) > 0 {
		for len(cells) < 7 {
			cells = append(cells, cellStyle.Render(""))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return lipgloss.JoinHorizontal(lipgloss.Top, grid, "  ", m.renderPanel())
}

func main() {
	mainModel := New()
	p := tea.NewProgram(mainModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
