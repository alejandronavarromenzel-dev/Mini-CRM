package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	a := app.NewWithID("minicrm")
	w := a.NewWindow("Mini CRM")
	w.Resize(fyne.NewSize(1200, 800))

	tabs := container.NewAppTabs(
		container.NewTabItem("Dashboard", dashboardView()),
		container.NewTabItem("Clientes", clientesView()),
		container.NewTabItem("Tareas", tareasView()),
		container.NewTabItem("Kanban", kanbanView()),
		container.NewTabItem("Reportes", reportesView()),
		container.NewTabItem("Ajustes", ajustesView()),
	)

	w.SetContent(tabs)
	w.ShowAndRun()
}

func dashboardView() fyne.CanvasObject {
	return widget.NewLabel("Dashboard – KPIs y seguimiento general")
}

func clientesView() fyne.CanvasObject {
	return widget.NewLabel("Clientes – Alta, edición, estado, etiquetas, notas")
}

func tareasView() fyne.CanvasObject {
	return widget.NewLabel("Tareas – Por cliente, fechas, prioridad, avance")
}

func kanbanView() fyne.CanvasObject {
	return widget.NewLabel("Kanban – Por hacer / En curso / Hecho")
}

func reportesView() fyne.CanvasObject {
	return widget.NewLabel("Reportes – Exportación Excel / CSV / PDF")
}

func ajustesView() fyne.CanvasObject {
	return widget.NewLabel("Ajustes – Backups, uso local, configuración")
}
