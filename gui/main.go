package main

import (
	ruleEntry "gohost/gui/component/rule-entry"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	appInstant := app.New()
	window := appInstant.NewWindow("gohost")

	ruleEntry := ruleEntry.MakeRuleEntry()
	certLable := widget.NewLabel("cert")
	tabs := container.NewAppTabs(
		container.NewTabItem("rule", ruleEntry),
		container.NewTabItem("cert", certLable),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	window.Resize(fyne.NewSize(800, 600))
	window.SetContent(tabs)
	window.ShowAndRun()

	
}