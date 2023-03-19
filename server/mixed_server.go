package server

import (
	"bufio"
	"fmt"
	"gohost/config"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	inHttp chan net.Conn
	inHttps chan net.Conn
)

func init(){
	inHttp = make(chan net.Conn)
	inHttps = make(chan net.Conn)
}

func proxyHttps(){
	for c := range inHttps {
		reader := bufio.NewReader(c)
		req, err:= http.ReadRequest(reader)
		fmt.Printf("req: %v\n, %s", req, err.Error())
		c.SetDeadline(time.Now())
	}
}

func proxyHttp(){
	for c := range inHttp {
		reader := bufio.NewReader(c)
		req, err:= http.ReadRequest(reader)
		fmt.Printf("req: %v\n, %s", req, err.Error())
		c.SetDeadline(time.Now())
	}
}


func Setup(conf *config.Config){
	addrPort := fmt.Sprintf("%s:%d",conf.ServerConfig.ListenAddr,conf.ServerConfig.ListenPort)
	ln, err := net.Listen(conf.ServerConfig.Network, addrPort);
	if  err != nil {
		log.Fatalf("启动监听失败: %s", err.Error())
		
	}
	log.Printf("服务启动成功： %s",addrPort)

	go proxyHttp()
	go proxyHttps()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("接受请求失败: %s", err.Error())
		}
		go handleConnWithProtocol(conn)
	}
}

func handleConnWithProtocol(conn net.Conn){
	//todo: 判断代理协议
	
	reader := bufio.NewReader(conn)
	head, err := reader.Peek(1)
	if err != nil {
		log.Printf("接受代理请求失败：%s", err.Error())
		return
	}

	// https ws wss的代理请求：请求方法是CONNECT
	if head[0] == byte('C') {
		req, err := http.ReadRequest(reader)
		if err != nil {
			log.Printf("接受CONNECT请求失败：%s", err.Error())
		}

		if _, err := fmt.Fprintf(conn, "HTTP/%d.%d %03d %s\r\n\r\n", req.ProtoMajor, req.ProtoMinor, http.StatusOK, "Connection established"); err != nil{
			log.Printf("发送接受代理请求失败：%s", err.Error())
			return
		}

		// inHttps <- conn
	}else{
		// http 请求
		// inHttp <- conn
		req, err:= http.ReadRequest(reader)
		fmt.Printf("req: %v\n, %s", req, err.Error())
	}


}