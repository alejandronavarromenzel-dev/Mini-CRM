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

	makeColumn := func(status string) fyne.CanvasObject {
		var tasks []models.Task
		selected := -1

		list := widget.NewList(
			func() int { return len(tasks) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(i int, o fyne.CanvasObject) {
				t := tasks[i]
				o.(*widget.Label).SetText(
					fmt.Sprintf("%s – %s (%d%%)", t.ClientName, t.Title, t.Progress),
				)
			},
		)

		list.OnSelected = func(id int) {
			selected = id
		}

		refresh := func() {
			tasks = loadKanbanTasks(status)
			selected = -1
			list.Refresh()
		}

		moveBtn := widget.NewButton("⇒", func() {
			if selected < 0 || selected >= len(tasks) {
				return
			}

			next := nextKanbanStatus(status)
			_, _ = db.DB.Exec(
				"UPDATE tasks SET status=? WHERE id=?",
				next,
				tasks[selected].ID,
			)
			refresh()
		})

		refresh()

		return container.NewBorder(
			widget.NewLabel(status),
			moveBtn,
			nil,
			nil,
			list,
		)
	}

	return container.NewGridWithColumns(
		3,
		makeColumn("Por hacer"),
		makeColumn("En curso"),
		makeColumn("Hecho"),
	)
}

func loadKanbanTasks(status string) []models.Task {
	rows, err := db.DB.Query(`
		SELECT t.id, c.name, t.title, t.progress
		FROM tasks t
		JOIN clients c ON c.id = t.client_id
		WHERE t.status=?
		ORDER BY c.name, t.title
	`, status)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []models.Task
	for rows.Next() {
		var t models.Task
		rows.Scan(&t.ID, &t.ClientName, &t.Title, &t.Progress)
		out = append(out, t)
	}
	return out
}

func nextKanbanStatus(current string) string {
	switch current {
	case "Por hacer":
		return "En curso"
	case "En curso":
		return "Hecho"
	default:
		return "Hecho"
	}
}
