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

func clientesView() fyne.CanvasObject {
	clients := loadClients()

	list := widget.NewList(
		func() int { return len(clients) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i int, o fyne.CanvasObject) {
			c := clients[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%s (%s)", c.Name, c.Status))
		},
	)

	addBtn := widget.NewButton("Nuevo Cliente", func() {
		showClientForm(nil, func() {
			clients = loadClients()
			list.Refresh()
		})
	})

	return container.NewBorder(
		container.NewHBox(addBtn),
		nil, nil, nil,
		list,
	)
}

func showClientForm(client *models.Client, onSave func()) {
	name := widget.NewEntry()
	status := widget.NewSelect([]string{"Activo", "Prospecto", "Inactivo"}, nil)
	owner := widget.NewEntry()
	tags := widget.NewEntry()
	notes := widget.NewMultiLineEntry()

	if client != nil {
		name.SetText(client.Name)
		status.SetSelected(client.Status)
		owner.SetText(client.Owner)
		tags.SetText(client.Tags)
		notes.SetText(client.Notes)
	}

	items := []*widget.FormItem{
		{Text: "Nombre", Widget: name},
		{Text: "Estado", Widget: status},
		{Text: "Responsable", Widget: owner},
		{Text: "Etiquetas", Widget: tags},
		{Text: "Notas", Widget: notes},
	}

	dialog.ShowForm("Cliente", "Guardar", "Cancelar", items, func(ok bool) {
		if !ok {
			return
		}
		if name.Text == "" {
			dialog.ShowError(fmt.Errorf("el nombre es obligatorio"), w)
			return
		}

		saveClient(&models.Client{
			Name:   name.Text,
			Status: status.Selected,
			Owner:  owner.Text,
			Tags:   tags.Text,
			Notes:  notes.Text,
		})
		onSave()
	}, w)
}

func saveClient(c *models.Client) {
	_, _ = db.DB.Exec(
		`INSERT INTO clients (name, status, owner, tags, notes) VALUES (?,?,?,?,?)`,
		c.Name, c.Status, c.Owner, c.Tags, c.Notes,
	)
}

func loadClients() []models.Client {
	rows, err := db.DB.Query(`SELECT id, name, status, owner, tags, notes FROM clients ORDER BY name`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []models.Client
	for rows.Next() {
		var c models.Client
		rows.Scan(&c.ID, &c.Name, &c.Status, &c.Owner, &c.Tags, &c.Notes)
		out = append(out, c)
	}
	return out
}
