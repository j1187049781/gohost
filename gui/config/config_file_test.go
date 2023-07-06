package config_test

import (
	"gohost/gui/config"
	"testing"
)

func TestConfigMarshal(t *testing.T) {

	config.GlobalConfig.UrlHandlerConfig.HandlerType = "rule"
	config.GlobalConfig.UrlHandlerConfig.Mappings = []config.UrlMapping{
		{
			Pattern: "//www.baidu.com",
			Target:  "//14.119.104.189",
		},
		{
			Pattern: "//www.bing.com",
			Target:  "//202.89.233.100",
		},
	}
	config.SaveConfig()
	

}

func TestConfigUnMarshal(t *testing.T) {

	config.LoadConfig()
	t.Log(config.GlobalConfig.UrlHandlerConfig)

}
