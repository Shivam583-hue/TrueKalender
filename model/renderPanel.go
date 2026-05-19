package model

import (
	"fmt"

	"github.com/Shivam583-hue/TrueKalender/styles"
	"github.com/charmbracelet/lipgloss"
)

// mock data
type taskKey struct {
	month int
	year  int
	day   int
}

var tasks = map[taskKey][]string{
	{month: int(currentTime.Month()), year: int(currentTime.Year()), day: 1}:     {"read bittorrent code", "complete math"},
	{month: int(currentTime.Month()), year: int(currentTime.Year()), day: 3}:     {"revise graphs"},
	{month: int(currentTime.Month()), year: int(currentTime.Year()), day: 7}:     {"start truekalender"},
	{month: int(currentTime.Month()), year: int(currentTime.Year()), day: 15}:    {"dentist appt", "grocery run", "call mom"},
	{month: int(currentTime.Month()) + 1, year: int(currentTime.Year()), day: 5}: {"next month task"},
}

func (m model) renderPanel() string {
	title := styles.PanelTitleStyle.Render(fmt.Sprintf("Day %d", m.focused))

	key := taskKey{month: m.monthNumber, year: m.year, day: m.focused}
	dayTasks, ok := tasks[key]

	var body string
	if !ok || len(dayTasks) == 0 {
		body = styles.EmptyStyle.Render("no tasks")
	} else {
		for _, t := range dayTasks {
			body += styles.TaskStyle.Render("• "+t) + "\n"
			body += styles.TaskDescStyle.Render("Todo") + "\n\n"
		}
	}

	return styles.PanelStyle.Width(m.panelWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left, title, "", body),
	)
}
