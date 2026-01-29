package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"minicrm/internal/db"
)

var w fyne.Window

func Run() {
	a := app.NewWithID("minicrm")

	if err := db.Init(); err != nil {
		panic(err)
	}

	w = a.NewWindow("Mini CRM")
	w.Resize(fyne.NewSize(1200, 800))

	tabs := container.NewAppTabs(
		container.NewTabItem("Dashboard", dashboardView()),
		container.NewTabItem("Clientes", clientesView()),
		container.NewTabItem("Tareas", tareasView()),
		container.NewTabItem("Tareas realizadas", tareasArchivadasView()),
		container.NewTabItem("Reportes", reportesView()),
		container.NewTabItem("Ajustes", ajustesView()),
	)

	w.SetContent(tabs)
	w.ShowAndRun()
}

func dashboardView() fyne.CanvasObject {
	return widget.NewLabel("Dashboard")
}

func reportesView() fyne.CanvasObject {
	return widget.NewLabel("Reportes")
}

func ajustesView() fyne.CanvasObject {
	return widget.NewLabel("Ajustes")
}
