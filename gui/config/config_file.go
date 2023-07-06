package config

import (
	"log"
	"os"
	"sync"
	"path"
	"gopkg.in/yaml.v3"
)

type Config struct {
	UrlHandlerConfig UrlHandlerConfig `yaml:"url_handler"`
}

var GlobalConfig Config = initDefaultConfig()
var lock = &sync.Mutex{}
var configPath = "conf/config.yaml"

func init() {
	userPath, err := os.UserHomeDir()
	if err != nil {
		log.Printf("get user home dir error: %v", err)
	}else{
		configPath = path.Join(userPath, configPath)
	}
	os.MkdirAll(path.Dir(configPath), 0644)
	
	LoadConfig()
}

func LoadConfig()  {
	lock.Lock()
	defer lock.Unlock()
	
	configYaml, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("read config file error: %v", err)
	}
	yaml.Unmarshal(configYaml, &GlobalConfig)
}

func SaveConfig(){
	lock.Lock()
	defer lock.Unlock()

	configYaml, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		log.Printf("marshal config error: %v", err)
		return
	}
	os.WriteFile(configPath, configYaml, 0644)
	
	
}

func initDefaultConfig() (Config){
	return Config{}
}
