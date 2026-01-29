package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"minicrm/internal/db"
	"minicrm/internal/models"
)

func tareasView() fyne.CanvasObject {

	type block struct {
		status string
		tasks  []models.Task
		list   *widget.List
	}

	blocks := []*block{}

	makeBlock := func(status string) *block {
		b := &block{status: status}

		b.list = widget.NewList(
			func() int { return len(b.tasks) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(i int, o fyne.CanvasObject) {
				t := b.tasks[i]
				o.(*widget.Label).SetText(
					fmt.Sprintf(
						"%s – %s\nEtiqueta: %s | Resp: %s",
						t.ClientName,
						t.Title,
						t.Tag,
						t.Owner,
					),
				)
			},
		)

		blocks = append(blocks, b)
		return b
	}

	refreshAll := func() {
		for _, b := range blocks {
			b.tasks = loadTasksByStatus(b.status)
			b.list.Refresh()
		}
	}

	makeBlockUI := func(b *block) fyne.CanvasObject {

		var actionBtn *widget.Button

		switch b.status {
		case "Por hacer":
			actionBtn = widget.NewButton("→ En curso", func() {
				moveSelectedTask(b, "En curso", refreshAll)
			})
		case "En curso":
			actionBtn = widget.NewButton("→ Hecho", func() {
				moveSelectedTask(b, "Hecho", refreshAll)
			})
		case "Hecho":
			actionBtn = widget.NewButton("Archivar", func() {
				archiveSelectedTask(b, refreshAll)
			})
		}

		return container.NewBorder(
			widget.NewLabelWithStyle(b.status, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			actionBtn,
			nil,
			nil,
			b.list,
		)
	}

	porHacer := makeBlock("Por hacer")
	enCurso := makeBlock("En curso")
	hecho := makeBlock("Hecho")

	refreshAll()

	addBtn := widget.NewButton("Nueva tarea", func() {
		showNewTaskForm(refreshAll)
	})

	return container.NewVBox(
		container.NewHBox(addBtn),
		makeBlockUI(porHacer),
		makeBlockUI(enCurso),
		makeBlockUI(hecho),
	)
}

/* ---------- acciones ---------- */

func moveSelectedTask(b *struct {
	status string
	tasks  []models.Task
	list   *widget.List
}, newStatus string, refresh func()) {

	id := b.list.
