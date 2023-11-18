package main

import (
	"gohost/config"
	"gohost/server"
	"io"
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	conf := config.InitConfig()
	s := server.NewMixedServer(&conf)
	s.Setup()
}