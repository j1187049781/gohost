package config


type UrlMapping struct {
	Pattern string `yaml:"pattern"`
	Target  string `yaml:"target"`
}