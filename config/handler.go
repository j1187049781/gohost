package config


type UrlMapping struct {
	Pattern string `yaml:"pattern"`
	Target  string `yaml:"target"`
}

type FormFileMapping struct {
	Pattern string `yaml:"pattern"`
	FormFileKeys []string `yaml:"form_file_keys"`
}