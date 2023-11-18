package config

import (
	"log"
	"os"
	"path"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerConfig *ServerConfig `yaml:"server_config"`
	Mappings     []UrlMapping  `yaml:"mappings"`
	RequestCopyFileUrls []FormFileMapping `yaml:"request_copy_file_urls"`
	lock         sync.RWMutex  `yaml:"-"`
}

type ServerConfig struct {
	Network    string `yaml:"network"`
	ListenAddr string `yaml:"listen_addr"`
	ListenPort int    `yaml:"listen_port"`
}

func InitConfig(options ...Option) (conf Config) {
	serveConf := &ServerConfig{
		Network:    "tcp",
		ListenAddr: "127.0.0.1",
		ListenPort: 8888,
	}
	conf = Config{
		ServerConfig: serveConf,
	}

	conf.loadConfig()

	for _, o := range options {
		o(&conf)
	}
	return
}

func (c *Config) SetMapping(m []UrlMapping){
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Mappings = m
}

func (c *Config) GetMapping() ([]UrlMapping){
	return c.Mappings
}

func (c *Config) loadConfig() {
	c.lock.RLock()
	defer c.lock.RUnlock()

	configPath := makeConfPath()

	configYaml, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("read config file error: %v", err)
		return
	}
	yaml.Unmarshal(configYaml, &c)
}


func (c *Config) SaveConfig(){
	c.lock.Lock()
	defer c.lock.Unlock()

	configYaml, err := yaml.Marshal(&c)
	if err != nil {
		log.Printf("marshal config error: %v", err)
		return
	}
	
	configPath := makeConfPath()
	os.WriteFile(configPath, configYaml, 0644)
	
	
}



func makeConfPath() string {
	var configPath = ".gohost/conf/config.yaml"

	
	os.MkdirAll(path.Dir(configPath), os.ModeDir|os.ModePerm)
	return configPath
}

