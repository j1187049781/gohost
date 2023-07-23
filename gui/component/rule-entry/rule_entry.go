package rule_entry

import (
	"gohost/config"
	"gohost/gui/tasks"
	"log"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func MakeRuleEntry(conf *config.Config) fyne.CanvasObject {
	ruleEntry := widget.NewMultiLineEntry()
	refreshRuleEntry(ruleEntry, conf)

	task := tasks.AutoSaveTask{
		NeedSave: make(chan struct{}),
		TaskFunc: func() {
			saveRuleEntry(ruleEntry,conf)
		},
	}
	ruleEntry.OnChanged = func(s string) {
		if task.NeedSave != nil && len(task.NeedSave) < 1{
			task.NeedSave <- struct{}{}
		}
		
	}
	task.RunBgTask()

	return ruleEntry
}

func refreshRuleEntry(ruleEntry *widget.Entry, conf *config.Config) {
	mappings := conf.GetMapping()
	for _, mapping := range mappings {
		ruleEntry.SetText(ruleEntry.Text + mapping.Pattern + "      " + mapping.Target + "\n")
	} 
}

func saveRuleEntry(ruleEntry *widget.Entry, conf *config.Config) {
	
	mappings, err := sparseLines(ruleEntry.Text)
	if err != nil {
		log.Printf(" save rule err: [%v]", err.Error())
		return
	}

	conf.SetMapping(mappings)
	conf.SaveConfig()
}

func sparseLines(text string) ([]config.UrlMapping , error){
	mappings := []config.UrlMapping{}
	lines := strings.Split(text,"\n")
	for _, line := range lines{
		tokens := strings.Fields(line)

		if len(tokens) !=2 {
			log.Printf("skip line: %v", line)
			continue
		}

		m := config.UrlMapping{}
		m.Pattern = tokens[0]
		m.Target = tokens[1]

		mappings = append(mappings, m)
	}
	return mappings, nil
}