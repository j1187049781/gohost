package server

import "net"

type ConnProxy struct {
	Conn      net.Conn
	Protocol  string
	KeepAlive bool
}
