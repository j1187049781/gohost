package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	appInstant := app.New()
	window := appInstant.NewWindow("gohost")

	tabs := container.NewAppTabs(
		container.NewTabItem("rule", widget.NewLabel("Hello")),
		container.NewTabItem("cert", widget.NewLabel("cert")),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	window.Resize(fyne.NewSize(800, 600))
	window.SetContent(tabs)
	window.ShowAndRun()

	
}