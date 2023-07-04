package config

type UrlHandlerConfig struct {
	HandlerType string `yaml:"handler_type"`
	Mappings	[]UrlMapping `yaml:"mappings"`
}

type UrlMapping struct {
	Pattern string `yaml:"pattern"`
	Target  string `yaml:"target"`
}