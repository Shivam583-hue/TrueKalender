package model

import (
	"fmt"
	"time"

	"github.com/Shivam583-hue/TrueKalender/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	minWidth  = 80
	minHeight = 24
)

var (
	dayHeaders  = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	currentTime = time.Now()
)

type model struct {
	quitting    bool
	help        help.Model
	lists       []list.Model
	focused     int
	loaded      bool
	width       int
	height      int
	cellWidth   int
	cellHeight  int
	panelWidth  int
	monthNumber int
	year        int
}

func New() *model {
	return &model{
		help:        help.New(),
		focused:     1,
		monthNumber: int(currentTime.Month()),
		year:        int(currentTime.Year()),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func getDays(monthNumber, year int) int {
	return time.Date(year, time.Month(monthNumber)+1, 0, 0, 0, 0, 0, time.Local).Day()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	numOfDays := getDays(m.monthNumber, m.year)
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
		case "L":
			m.monthNumber++
			if m.monthNumber > 12 {
				m.monthNumber = 1
				m.year++
			}
			m.focused = 1
		case "H":
			m.monthNumber--
			if m.monthNumber < 1 {
				m.monthNumber = 12
				m.year--
			}
			m.focused = 1
		}
	}
	return m, nil
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
		return styles.TooSmallStyle.
			Width(m.width).
			Height(m.height).
			Render(msg)
	}

	cw := m.cellWidth
	ch := m.cellHeight

	scaledCell := styles.CellStyle.Width(cw).Height(ch)
	scaledFocused := styles.FocusedCellStyle.Width(cw).Height(ch)
	scaledHeader := styles.HeaderStyle.Width(cw + 2)

	headers := make([]string, 7)
	for i, d := range dayHeaders {
		headers[i] = scaledHeader.Render(d)
	}

	numOfDays := getDays(m.monthNumber, m.year)
	rows := []string{lipgloss.JoinHorizontal(lipgloss.Top, headers...)}

	var cells []string
	firstWeekday := int(time.Date(m.year, time.Month(m.monthNumber), 1, 0, 0, 0, 0, time.Local).Weekday())
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
		styles.MonthTitleStyle.Render(time.Month(m.monthNumber).String()),
		styles.YearStyle.Render(fmt.Sprintf("%d", m.year)),
	)
	helpBar := styles.HelpBarStyle.Render(
		"h/← prev day  •  l/→ next day  •  j/↓ week down  •  k/↑ week up  •  H prev month  •  L next month  •  q quit",
	)

	gridWTitle := lipgloss.JoinVertical(lipgloss.Left, monthYear, grid)
	content := lipgloss.JoinHorizontal(lipgloss.Top, gridWTitle, "  ", m.renderPanel())

	return lipgloss.JoinVertical(lipgloss.Left, content, "", helpBar)
}
