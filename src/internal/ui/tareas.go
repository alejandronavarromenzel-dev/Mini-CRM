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
   VISTA PRINCIPAL TAREAS
   ========================= */

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
			widget.NewLabelWithStyle(
				b.status,
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			),
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

/* =========================
   ACCIONES
   ========================= */

func moveSelectedTask(
	b *block,
	newStatus string,
	refresh func(),
) {

	id := b.list.Selected
	if id < 0 || id >= len(b.tasks) {
		return
	}

	_, _ = db.DB.Exec(
		"UPDATE tasks SET status=? WHERE id=?",
		newStatus,
		b.tasks[id].ID,
	)
	refresh()
}

func archiveSelectedTask(
	b *block,
	refresh func(),
) {

	id := b.list.Selected
	if id < 0 || id >= len(b.tasks) {
		return
	}

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
			refresh()
		},
		w,
	)
}

/* =========================
   ALTA DE TAREAS
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
   QUERIES
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
