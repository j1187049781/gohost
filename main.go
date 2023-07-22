package main

import (
	"gohost/config"
	"gohost/server"
)

func main() {
	conf := config.InitConfig()
	server.Setup(&conf)
}