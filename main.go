package main

import (
	"gohost/config"
	"gohost/server"
)

func main() {
	conf := config.InitDefaultConfig()
	server.Setup(&conf)
}