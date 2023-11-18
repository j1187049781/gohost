package config_test

import (
	"gohost/config"
	"testing"
)

func TestSave(t *testing.T) {
	conf := config.InitConfig()
	conf.RequestCopyFileUrls = append(conf.RequestCopyFileUrls, config.FormFileMapping{
		Pattern: "/user/firconfig",
		FormFileKeys: []string{"file"},
	})
	conf.SaveConfig()
}