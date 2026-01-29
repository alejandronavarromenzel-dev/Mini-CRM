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
   TIPO
   ========================= */

type taskBlock struct {
	status        string
	list          *widget.List
	tasks         []models.Task
	selectedIndex int
}

/* =========================
   VISTA PRINCIPAL
   ========================= */

func tareasView() fyne.CanvasObject {

	porHacer := newTaskBlock("Por hacer")
	enCurso := newTaskBlock("En curso")
	hecho := newTaskBlock("Hecho")

	refreshAll := func() {
		porHacer.refresh()
		enCurso.refresh()
		hecho.refresh()
	}

	addBtn := widget.NewButton("Nueva tarea", func() {
		showNewTaskForm(refreshAll)
	})

	refreshAll()

	return container.NewVBox(
		container.NewHBox(addBtn),
		porHacer.view("→ En curso", "En curso", refreshAll),
		enCurso.view("→ Hecho", "Hecho", refreshAll),
		hecho.view("Archivar", "ARCHIVE", refreshAll),
	)
}

/* =========================
   BLOQUE
   ========================= */

func newTaskBlock(status string) *taskBlock {
	b := &taskBlock{
		status:        status,
		selectedIndex: -1,
	}

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

	b.list.OnSelected = func(id int) {
		b.selectedIndex = id
	}

	return b
}

func (b *taskBlock) refresh() {
	b.tasks = loadTasksByStatus(b.status)
	b.selectedIndex = -1
	b.list.Refresh()
}

func (b *taskBlock) view(
	buttonLabel string,
	nextStatus string,
	refreshAll func(),
) fyne.CanvasObject {

	btn := widget.NewButton(buttonLabel, func() {

		id := b.selectedIndex
		if id < 0 || id >= len(b.tasks) {
			return
		}

		if nextStatus == "ARCHIVE" {
			dialog.ShowConfirm(
				"Archivar tarea",
				"¿Deseas archivar esta tarea?",
				func(ok bool) {
					if !ok {
						return
					}
					_, _ = db.DB.Exec(
						"UPDATE tasks SET archived=1 WHERE id=?",
						b.tasks[id].ID,
					)
					refreshAll()
				},
				w,
			)
			return
		}

		_, _ = db.DB.Exec(
			"UPDATE tasks SET status=? WHERE id=?",
			nextStatus,
			b.tasks[id].ID,
		)
		refreshAll()
	})

	return container.NewBorder(
		widget.NewLabelWithStyle(
			b.status,
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
		btn,
		nil,
		nil,
		b.list,
	)
}

/* =========================
   ALTA
   ========================= */

func showNewTaskForm(onSave func()) {
	clients := loadClients()
	clientNames := []string{}
	clientMap := map[string]int64{}

	for _, c := range clients {
		clientNames = append(clientNames, c.Name)
		clientMap[c.Name] = c.ID
	}

	clientSel := widget.NewSelect(clientNames, nil)
	title := widget.NewEntry()
	tag := widget.NewEntry()
	owner := widget.NewEntry()

	items := []*widget.FormItem{
		{Text: "Cliente", Widget: clientSel},
		{Text: "Tarea", Widget: title},
		{Text: "Etiqueta", Widget: tag},
		{Text: "Responsable", Widget: owner},
	}

	dialog.ShowForm(
		"Nueva tarea",
		"Guardar",
		"Cancelar",
		items,
		func(ok bool) {
			if !ok {
				return
			}

			_, _ = db.DB.Exec(
				`INSERT INTO tasks (client_id, title, tag, owner, status, archived)
				 VALUES (?,?,?,?, 'Por hacer', 0)`,
				clientMap[clientSel.Selected],
				title.Text,
				tag.Text,
				owner.Text,
			)
			onSave()
		},
		w,
	)
}

/* =========================
   QUERY
   ========================= */

func loadTasksByStatus(status string) []models.Task {
	rows, err := db.DB.Query(`
		SELECT t.id, c.name, t.title, t.tag, t.owner
		FROM tasks t
		JOIN clients c ON c.id = t.client_id
		WHERE t.status=? AND t.archived=0
		ORDER BY c.name, t.title
	`, status)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []models.Task
	for rows.Next() {
		var t models.Task
		rows.Scan(&t.ID, &t.ClientName, &t.Title, &t.Tag, &t.Owner)
		out = append(out, t)
	}
	return out
}
