package model

import (
	"fmt"

	"github.com/Shivam583-hue/TrueKalender/db"
	"github.com/Shivam583-hue/TrueKalender/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m model) renderPanel() string {
	title := styles.PanelTitleStyle.Render(fmt.Sprintf("Day %d", m.focused))

	dayTasks, err := db.FetchByDate(m.year, m.monthNumber, m.focused)

	var body string
	if err != nil || len(dayTasks) == 0 {
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
