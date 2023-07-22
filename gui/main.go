package main

import (
	"gohost/config"
	ruleEntry "gohost/gui/component/rule-entry"
	"gohost/server"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)


func main() {
	conf := config.InitConfig()
	s := server.NewMixedServer(&conf)
	s.Setup()

	
	appInstant := app.New()
	window := appInstant.NewWindow("gohost")

	ruleEntry := ruleEntry.MakeRuleEntry(&conf)
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