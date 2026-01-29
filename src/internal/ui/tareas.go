package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"minicrm/internal/db"
	"minicrm/internal/models"
)

func tareasView() fyne.CanvasObject {
	clients := loadClients()
	clientNames := []string{}
	clientMap := map[string]int64{}

	for _, c := range clients {
		clientNames = append(clientNames, c.Name)
		clientMap[c.Name] = c.ID
	}

	clientSelect := widget.NewSelect(clientNames, nil)
	list := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i int, o fyne.CanvasObject) {},
	)

	clientSelect.OnChanged = func(name string) {
		clientID := clientMap[name]
		tasks := loadTasks(clientID)
		list.Length = func() int { return len(tasks) }
		list.UpdateItem = func(i int, o fyne.CanvasObject) {
			t := tasks[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%s (%s %d%%)", t.Title, t.Status, t.Progress))
		}
		list.Refresh()
	}

	addBtn := widget.NewButton("Nueva Tarea", func() {
		if clientSelect.Selected == "" {
			dialog.ShowInformation("Atención", "Selecciona un cliente", w)
			return
		}
		showTaskForm(clientMap[clientSelect.Selected], func() {
			clientSelect.OnChanged(clientSelect.Selected)
		})
	})

	return container.NewBorder(
		container.NewHBox(widget.NewLabel("Cliente:"), clientSelect, addBtn),
		nil, nil, nil,
		list,
	)
}

func showTaskForm(clientID int64, onSave func()) {
	title := widget.NewEntry()
	status := widget.NewSelect([]string{"Por hacer", "En curso", "Hecho"}, nil)
	priority := widget.NewSelect([]string{"Baja", "Media", "Alta"}, nil)
	owner := widget.NewEntry()
	progress := widget.NewEntry()
	progress.SetText("0")

	items := []*widget.FormItem{
		{Text: "Título", Widget: title},
		{Text: "Estado", Widget: status},
		{Text: "Prioridad", Widget: priority},
		{Text: "Responsable", Widget: owner},
		{Text: "% Avance", Widget: progress},
	}

	dialog.ShowForm("Tarea", "Guardar", "Cancelar", items, func(ok bool) {
		if !ok {
			return
		}
		p, _ := strconv.Atoi(progress.Text)
		_, _ = db.DB.Exec(
			`INSERT INTO tasks (client_id, title, status, priority, owner, progress) VALUES (?,?,?,?,?,?)`,
			clientID, title.Text, status.Selected, priority.Selected, owner.Text, p,
		)
		onSave()
	}, w)
}

func loadTasks(clientID int64) []models.Task {
	rows, err := db.DB.Query(`SELECT id, title, status, priority, owner, progress FROM tasks WHERE client_id=?`, clientID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []models.Task
	for rows.Next() {
		var t models.Task
		rows.Scan(&t.ID, &t.Title, &t.Status, &t.Priority, &t.Owner, &t.Progress)
		out = append(out, t)
	}
	return out
}
