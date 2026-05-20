package model

import (
	"database/sql"
	"fmt"

	"github.com/Shivam583-hue/TrueKalender/db"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

type syncResultMsg string

func syncToKanban(year, month, day int) tea.Cmd {
	return func() tea.Msg {
		tasks, err := db.FetchByDate(year, month, day)
		if err != nil {
			return syncResultMsg("sync failed: could not fetch tasks")
		}
		if len(tasks) == 0 {
			return syncResultMsg("nothing to sync")
		}

		kanban, err := sql.Open("sqlite3", "./task.db")
		if err != nil {
			return syncResultMsg("sync failed: could not open kanban db")
		}
		defer kanban.Close()

		for _, title := range tasks {
			_, err := kanban.Exec(
				"INSERT INTO tasks(title, status) VALUES(?, ?)", title, 0,
			)
			if err != nil {
				return syncResultMsg(fmt.Sprintf("sync failed: %v", err))
			}
		}

		return syncResultMsg(fmt.Sprintf("synced %d tasks to kanban", len(tasks)))
	}
}
