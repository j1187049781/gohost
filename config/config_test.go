package config_test

import (
	"gohost/config"
	"testing"
)

func TestSave(t *testing.T) {
	conf := config.InitConfig()
	mappings := []config.UrlMapping {{Pattern: "//www.baidu.com", Target: "//14.119.104.254"}}
	conf.SetMapping(mappings)
	conf.SaveConfig()
}