package rule_entry

import (
	"gohost/gui/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func MakeRuleEntry() fyne.CanvasObject {
	ruleEntry := widget.NewMultiLineEntry()
	RefreshRuleEntry(ruleEntry)
	return ruleEntry
}

func RefreshRuleEntry(ruleEntry *widget.Entry) {
	mappings := config.GlobalConfig.UrlHandlerConfig.Mappings
	for _, mapping := range mappings {
		ruleEntry.SetText(ruleEntry.Text + mapping.Pattern + "      " + mapping.Target + "\n")
	}
}