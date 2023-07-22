package rule_entry

import (
	"gohost/config"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func MakeRuleEntry(conf *config.Config) fyne.CanvasObject {
	ruleEntry := widget.NewMultiLineEntry()
	refreshRuleEntry(ruleEntry, conf)
	go saveRuleEntry(ruleEntry,conf)
	return ruleEntry
}

func refreshRuleEntry(ruleEntry *widget.Entry, conf *config.Config) {
	mappings := conf.GetMapping()
	for _, mapping := range mappings {
		ruleEntry.SetText(ruleEntry.Text + mapping.Pattern + "      " + mapping.Target + "\n")
	} 
}

func saveRuleEntry(ruleEntry *widget.Entry, conf *config.Config) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	
	for c := range ticker.C {
		mappings, err := sparseLines(ruleEntry.Text)
		if err != nil {
			log.Printf(" save rule err: [%v], time: %v", err.Error(), c)
			continue
		}

		conf.SetMapping(mappings)
		conf.SaveConfig()
	}
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