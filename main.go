package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var models []tea.Model

var (
	currentTime = time.Now()
	monthNumber = int(currentTime.Month())
	year        = int(currentTime.Year())
)

func getDays() int {
	return time.Date(year, time.Month(monthNumber)+1, 0, 0, 0, 0, 0, time.Local).Day()
}

const (
	minWidth  = 80
	minHeight = 24
)

var (
	purple = lipgloss.Color("135")

	focusedCellStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("255")).
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("55")).
				Bold(true).
				Align(lipgloss.Center, lipgloss.Center)

	cellStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(purple).
			Foreground(purple).
			Align(lipgloss.Center, lipgloss.Center)
	monthTitleStyle = lipgloss.NewStyle().
			Background(purple).
			Foreground(lipgloss.Color("255")).
			Bold(true).
			Padding(0, 2)

	yearStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(purple)

	helpBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(purple).
			Align(lipgloss.Center)

	panelStyle = lipgloss.NewStyle().
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

	tooSmallStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Align(lipgloss.Center, lipgloss.Center)
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
	quitting   bool
	help       help.Model
	lists      []list.Model
	focused    int
	loaded     bool
	width      int
	height     int
	cellWidth  int
	cellHeight int
	panelWidth int
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
	numOfDays := getDays()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		panelWidth := m.width * 28 / 100
		gridWidth := m.width - panelWidth - 4
		cellW := gridWidth/7 - 2
		if cellW < 4 {
			cellW = 4
		}
		cellH := (m.height-6)/6 - 2
		if cellH < 1 {
			cellH = 1
		}
		m.cellWidth = cellW
		m.cellHeight = cellH
		m.panelWidth = panelWidth

		if !m.loaded {
			m.loaded = true
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "l", "right":
			m.focused = min(m.focused+1, numOfDays)
		case "h", "left":
			m.focused = max(m.focused-1, 1)
		case "j", "down":
			m.focused = min(m.focused+7, numOfDays)
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

	return panelStyle.Width(m.panelWidth).Render(
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

	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf(
			"Terminal too small: %dx%d\nMinimum required: %dx%d\nPlease resize your terminal.",
			m.width, m.height, minWidth, minHeight,
		)
		return tooSmallStyle.
			Width(m.width).
			Height(m.height).
			Render(msg)
	}

	cw := m.cellWidth
	ch := m.cellHeight

	scaledCell := cellStyle.Width(cw).Height(ch)
	scaledFocused := focusedCellStyle.Width(cw).Height(ch)
	scaledHeader := headerStyle.Width(cw + 2)

	headers := make([]string, 7)
	for i, d := range dayHeaders {
		headers[i] = scaledHeader.Render(d)
	}

	numOfDays := getDays()
	rows := []string{lipgloss.JoinHorizontal(lipgloss.Top, headers...)}

	var cells []string
	firstWeekday := int(time.Date(year, time.Month(monthNumber), 1, 0, 0, 0, 0, time.Local).Weekday())
	for i := 0; i < firstWeekday; i++ {
		cells = append(cells, scaledCell.Render(""))
	}
	for day := 1; day <= numOfDays; day++ {
		style := scaledCell
		if day == m.focused {
			style = scaledFocused
		}
		cells = append(cells, style.Render(fmt.Sprintf("%d", day)))
		if (firstWeekday+day)%7 == 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
			cells = nil
		}
	}
	if len(cells) > 0 {
		for len(cells) < 7 {
			cells = append(cells, scaledCell.Render(""))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	monthYear := lipgloss.JoinHorizontal(lipgloss.Left,
		monthTitleStyle.Render(currentTime.Month().String()),
		yearStyle.Render(fmt.Sprintf("%d", year)),
	)

	helpBar := helpBarStyle.Render(
		"h/← left  •  l/→ right  •  j/↓ week down  •  k/↑ week up  •  q quit",
	)

	gridWTitle := lipgloss.JoinVertical(lipgloss.Left, monthYear, grid)
	content := lipgloss.JoinHorizontal(lipgloss.Top, gridWTitle, "  ", m.renderPanel())

	return lipgloss.JoinVertical(lipgloss.Left, content, "", helpBar)
}

func main() {
	mainModel := New()
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
