package config

type Option func(s *Config)

func WithListenAddr(addr string) Option{
	return func(s *Config) {
		s.ServerConfig.ListenAddr = addr
	}
}

func WithListenPort(port int) Option {
	return func(s *Config) {
		s.ServerConfig.ListenPort = port
	}
}