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

	type column struct {
		status string
		tasks  []models.Task
		list   *widget.List
	}

	columns := []*column{}

	makeColumn := func(status string) *column {
		col := &column{status: status}

		col.list = widget.NewList(
			func() int { return len(col.tasks) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(i int, o fyne.CanvasObject) {
				t := col.tasks[i]
				o.(*widget.Label).SetText(
					fmt.Sprintf("%s – %s (%d%%)", t.ClientName, t.Title, t.Progress),
				)
			},
		)

		col.list.OnSelected = func(id int) {
			col.list.Selected = id
		}

		columns = append(columns, col)
		return col
	}

	refreshAll := func() {
		for _, c := range columns {
			c.tasks = loadKanbanTasks(c.status)
			c.list.Refresh()
		}
	}

	makeColumnUI := func(c *column) fyne.CanvasObject {
		moveBtn := widget.NewButton("⇒", func() {
			id := c.list.Selected
			if id < 0 || id >= len(c.tasks) {
				return
			}

			next := nextKanbanStatus(c.status)
			_, _ = db.DB.Exec(
				"UPDATE tasks SET status=? WHERE id=?",
				next,
				c.tasks[id].ID,
			)
			refreshAll()
		})

		return container.NewBorder(
			widget.NewLabel(c.status),
			moveBtn,
			nil,
			nil,
			c.list,
		)
	}

	colTodo := makeColumn("Por hacer")
	colDoing := makeColumn("En curso")
	colDone := makeColumn("Hecho")

	refreshAll()

	grid := container.NewGridWithColumns(
		3,
		makeColumnUI(colTodo),
		makeColumnUI(colDoing),
		makeColumnUI(colDone),
	)

	refreshBtn := widget.NewButton("Refrescar", refreshAll)

	return container.NewBorder(
		container.NewHBox(widget.NewLabel("Kanban"), refreshBtn),
		nil, nil, nil,
		grid,
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
