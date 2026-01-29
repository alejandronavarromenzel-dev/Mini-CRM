package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"

	"minicrm/internal/db"
)

func tareasArchivadasView() fyne.CanvasObject {

	rows, _ := db.DB.Query(`
		SELECT t.tag, c.name, t.title, t.status, t.owner
		FROM tasks t
		JOIN clients c ON c.id = t.client_id
		WHERE t.archived = 1
		ORDER BY c.name, t.title
	`)

	data := [][]string{}
	for rows.Next() {
		var tag, client, title, status, owner string
		rows.Scan(&tag, &client, &title, &status, &owner)
		data = append(data, []string{tag, client, title, status, owner})
	}

	headers := []string{"Etiqueta", "Cliente", "Tarea", "Estado", "Responsable"}

	table := widget.NewTable(
		func() (int, int) { return len(data)+1, len(headers) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			lbl := o.(*widget.Label)
			if id.Row == 0 {
				lbl.SetText(headers[id.Col])
				return
			}
			lbl.SetText(data[id.Row-1][id.Col])
		},
	)

	return container.NewBorder(nil, nil, nil, nil, table)
}
