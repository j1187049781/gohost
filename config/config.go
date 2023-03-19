package config

type Config struct {
	ServerConfig *ServerConfig
}

type ServerConfig struct {
	Network    string
	ListenAddr string
	ListenPort int
}

func InitDefaultConfig(options ...Option) (conf Config) {
	serveConf := &ServerConfig{
		Network:    "tcp",
		ListenAddr: "127.0.0.1",
		ListenPort: 8888,
	}
	conf = Config{
		ServerConfig: serveConf,
	}

	for _, o := range options {
		o(&conf)
	}
	return
}
