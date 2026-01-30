package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"minicrm/internal/db"
)

var w fyne.Window

func Run() {
	a := app.NewWithID("minicrm")

	w = a.NewWindow("Mini CRM")
	w.Resize(fyne.NewSize(1200, 800))

	// Inicializar DB (si falla, muestra error en pantalla)
	if err := db.Init(); err != nil {
		dialog.ShowError(err, w)
		w.ShowAndRun()
		return
	}

	tabs := container.NewAppTabs(
		container.NewTabItem("Dashboard", widget.NewLabel("Dashboard")),
		container.NewTabItem("Clientes", clientesView()),
		container.NewTabItem("Tareas", tareasView()),
		container.NewTabItem("Tareas realizadas", tareasArchivadasView()),
		container.NewTabItem("Reportes", widget.NewLabel("Reportes")),
		container.NewTabItem("Ajustes", widget.NewLabel("Ajustes")),
	)

	w.SetContent(tabs)
	w.ShowAndRun()
}
