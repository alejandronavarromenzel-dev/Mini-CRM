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

/* =========================
   TIPOS
   ========================= */

type taskBlock struct {
	status string
	tasks  []models.Task
	list   *widget.List
}

/* =========================
   VISTA PRINCIPAL
   ========================= */

func tareasView() fyne.CanvasObject {

	blocks := []*taskBlock{}

	makeBlock := func(status string) *taskBlock {
		b := &taskBlock{status: status}

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

	makeBlockUI := func(b *taskBlock) fyne.CanvasObject {

		var actionBtn *widget.Button

		switch b.status {
		case "Por hacer":
			actionBtn = w
