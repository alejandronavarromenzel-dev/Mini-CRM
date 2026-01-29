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
	var tasks []models.Task
	var selectedTask *models.Task

	clientSelect := widget.NewSelect([]string{"Todos"}, nil)
	refreshClients := func() {
		clients := loadClients()
		opts := []string{"Todos"}
		for _, c := range clients {
			opts = append(opts, c.Name)
		}
		clientSelect.Options = opts
		clientSelect.SetSelected("Todos")
	}

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
		if id >= 0 && id < len(tasks) {
			selectedTask = &tasks[id]
		}
	}

	refreshTasks := func() {
		if clientSelect.Selected == "Todos" {
			tasks = loadAllTasks()
		} else {
			tasks = loadTasksByClientName(clientSelect.Selected)
		}
		list.Refresh()
	}

	clientSelect.OnChanged = func(string) {
		refreshTasks()
	}

	addBtn := widget.NewButton("Nueva Tarea", func() {
		showTaskForm(0, func() {
			refreshTasks()
		})
	})

	editBtn := widget.NewButton("Editar % Avance", func() {
		if selectedTask == nil {
			dialog.ShowInformation("Atención", "Selecciona una tarea", w)
			return
		}
		showProgressEdit(selectedTask, refreshTasks)
	})

	refreshClients()
	refreshTasks()

	top := container.NewHBox(
		widget.NewLabel("Cliente:"),
		clientSelect,
		addBtn,
		editBtn,
	)

	return container.NewBorder(top, nil, nil, nil, list)
}

func showTaskForm(_ int64, onSave func()) {
	clients := loadClients()
	clientNames := []string{}
	clientMap := map[string]int64{}

	for _, c := range clients {
		clientNames = append(clientNames, c.Name)
		clientMap[c.Name] = c.ID
	}

	clientSel := widget.NewSelect(clientNames, nil)
	title := widget.NewEntry()
	status := widget.NewSelect([]string{"Por hacer", "En curso", "Hecho"}, nil)
	priority := widget.NewSelect([]string{"Baja", "Media", "Alta"}, nil)
	progress := widget.NewEntry()
	progress.SetText("0")

	items := []*widget.FormItem{
		{Text: "Cliente", Widget: clientSel},
		{Text: "Título", Widget: title},
		{Text: "Estado", Widget: status},
		{Text: "Prioridad", Widget: priority},
		{Text: "% Avance", Widget: progress},
	}

	dialog.ShowForm("Nueva Tarea", "Guardar", "Cancelar", items, func(ok bool) {
		if !ok {
			return
		}

		p, _ := strconv.Atoi(progress.Text)
		_, _ = db.DB.Exec(
			`INSERT INTO tasks (client_id, title, status, priority, progress)
			 VALUES (?,?,?,?,?)`,
			clientMap[clientSel.Selected],
			title.Text,
			status.Selected,
			priority.Selected,
			p,
		)
		onSave()
	}, w)
}

func showProgressEdit(t *models.Task, onSave func()) {
	entry := widget.NewEntry()
	entry.SetText(strconv.Itoa(t.Progress))

	dialog.ShowForm(
		"Editar avance",
		"Guardar",
		"Ca
