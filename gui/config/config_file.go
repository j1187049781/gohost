package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	UrlHandlerConfig UrlHandlerConfig `yaml:"url_handler"`
}

var GlobalConfig Config
var lock = &sync.Mutex{}
var path = "conf/config.yaml"

func LoadConfig()  {
	lock.Lock()
	defer lock.Unlock()
	
	configYaml, err := os.ReadFile(path)
	if err != nil {
		log.Printf("read config file error: %v", err)
	}
	yaml.Unmarshal(configYaml, &GlobalConfig)
}

func SaveConfig(){
	lock.Lock()
	defer lock.Lock()

	configYaml, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		log.Printf("marshal config error: %v", err)
		return
	}
	os.WriteFile(path, configYaml, 0644)
	
	
}

func initDefaultConfig() (Config){
	return Config{}
}
