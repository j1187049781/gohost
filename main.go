package main

import (
	"gohost/config"
	"gohost/server"
)

func main() {
	conf := config.InitConfig()
	s := server.NewMixedServer(&conf)
	s.Setup()
}