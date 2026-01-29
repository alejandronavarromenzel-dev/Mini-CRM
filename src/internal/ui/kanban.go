package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"minicrm/internal/db"
	"minicrm/internal/models"
)

func kanbanView() fyne.CanvasObject {
	statuses := []string{"Por hacer", "En curso", "Hecho"}

	makeColumn := func(status string) fyne.CanvasObject {
		var selectedIndex = -1
		tasks := loadTasksByStatus(status)

		list := widget.NewList(
			func() int { return len(tasks) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(i int, o fyne.CanvasObject) {
				t := tasks[i]
				o.(*widget.Label).SetText(fmt.Sprintf("%s (%d%%)", t.Title, t.Progress))
			},
		)

		list.OnSelected = func(id int) {
			selectedIndex = id
		}

		moveBtn := widget.NewButton("Mover →", func() {
			if selectedIndex < 0 || selectedIndex >= len(tasks) {
				return
			}

			next := nextStatus(status)
			_, _ = db.DB.Exec(
				"UPDATE tasks SET status=? WHERE id=?",
				next,
				tasks[selectedIndex].ID,
			)

			// recargar columna
			tasks = loadTasksByStatus(status)
			selectedIndex = -1
			list.Refresh()
		})

		return container.NewBorder(
			widget.NewLabel(status),
			moveBtn,
			nil,
			nil,
			list,
		)
	}

	cols := []fyne.CanvasObject{}
	for _, s := range statuses {
		cols = append(cols, makeColumn(s))
	}

	return container.NewGridWithColumns(3, cols...)
}

func loadTasksByStatus(status string) []models.Task {
	rows, err := db.DB.Query(
		`SELECT id, title, progress FROM tasks WHERE status=?`,
		status,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []models.Task
	for rows.Next() {
		var t models.Task
		rows.Scan(&t.ID, &t.Title, &t.Progress)
		out = append(out, t)
	}
	return out
}

func nextStatus(current string) string {
	switch current {
	case "Por hacer":
		return "En curso"
	case "En curso":
		return "Hecho"
	default:
		return "Hecho"
	}
}
