package server

import (
	"crypto/tls"
	"fmt"
	"gohost/common/cert"
	N "gohost/common/net"
	"gohost/config"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	inHttp  chan net.Conn
	inHttps chan net.Conn
)

func init() {
	inHttp = make(chan net.Conn, 2000)
	inHttps = make(chan net.Conn, 2000)
}

func proxyHttps() {
	//todo: reuse client
	client := &http.Client{}
	for c := range inHttps {

		tlsConn := tls.Server(c, &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return cert.GetSignedCert(info.ServerName)
			},
		})
		err := tlsConn.Handshake()
		if err != nil {
			fmt.Printf("tls handshake error: %s", err.Error())
			continue
		}
		
		go func() {
			for {
				reader := N.NewReader(tlsConn)
				req, err := http.ReadRequest(reader.Reader())
				if err != nil {
					fmt.Printf("req: %v\n, %s", req, err.Error())
					return
				}
				fmt.Printf("connection: %s --> %s",  tlsConn.RemoteAddr(), req.Host)
				
				
				// 移除http请求头中的hop-by-hop headers
				removeHopByHopHeaders(req.Header)
	
				// 解析req请求体，获取请求头和请求体
	
				// reuse existing http.Request
				req.RequestURI = ""
				req.URL.Scheme = "https"
				req.URL.Host = req.Host
				// req.Header.Set("Connection", "false")
	
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("resp: %v, %s\n", resp, err)
					return
				}
	
				removeHopByHopHeaders(resp.Header)
	
				// 不需要resp.Body.Close
				// The Response Body is closed after it is sent
				resp.Close = false
				if err := resp.Write(tlsConn); err != nil {
					fmt.Printf("resp reply error: %v, %s", resp, err)
					return
				}
			}
		}()
		//fixme: close tlsConn
	}
}

func proxyHttp() {
	//todo: reuse client
	client := &http.Client{}

	for c := range inHttp {
		n := c
		go func ()  {
			for {
				
				reader := N.NewReader(n)
				req, err := http.ReadRequest(reader.Reader())
				if err != nil {
					fmt.Printf("req: %v\n, %s", req, err.Error())
					return
				}
				// 读取Header中的keepAlive，判断是是否需要keepAlive
	
				// 移除http请求头中的hop-by-hop headers
				removeHopByHopHeaders(req.Header)
	
				// 解析req请求体，获取请求头和请求体
	
				// reuse existing http.Request
				req.RequestURI = ""
	
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("resp: %v, %s", resp, err)
					return
				}
	
				removeHopByHopHeaders(resp.Header)
	
				// 不需要resp.Body.Close
				// The Response Body is closed after it is sent
				if err := resp.Write(n); err != nil {
					fmt.Printf("resp reply error: %v, %s", resp, err)
					return
				}
			}
		}()

	}
}

func Setup(conf *config.Config) {
	addrPort := fmt.Sprintf("%s:%d", conf.ServerConfig.ListenAddr, conf.ServerConfig.ListenPort)
	ln, err := net.Listen(conf.ServerConfig.Network, addrPort)
	if err != nil {
		log.Fatalf("启动监听失败: %s", err.Error())

	}
	log.Printf("服务启动成功： %s", addrPort)

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

func handleConnWithProtocol(conn net.Conn) {
	//todo: 判断代理协议

	bConn := N.NewReader(conn)
	head, err := bConn.Peek(1)
	if err != nil {
		log.Printf("接受代理请求失败：%s", err.Error())
		return
	}

	// https ws wss的代理请求：请求方法是CONNECT
	if head[0] == byte('C') {
		req, err := http.ReadRequest(bConn.Reader())
		if err != nil {
			log.Printf("接受CONNECT请求失败：%s", err.Error())
		}

		// 发送接受代理响应
		if _, err := fmt.Fprintf(conn, "HTTP/%d.%d %03d %s\r\n\r\n", req.ProtoMajor, req.ProtoMinor, http.StatusOK, "Connection established"); err != nil {
			log.Printf("发送接受代理请求失败：%s", err.Error())
			return
		}

		inHttps <- bConn
		log.Printf("接受https代理请求: %s --> %s", bConn.RemoteAddr().String(), bConn.LocalAddr().String())
	} else {
		// http 请求
		inHttp <- bConn
		log.Printf("接受http代理请求: %s --> %s", bConn.RemoteAddr().String(), bConn.LocalAddr().String())
	}

}

// 移除http请求头中的hop-by-hop headers
// Hop-by-hop headers 是 HTTP 协议中的一种头部，用于控制在传输过程中每个单独的代理或网关所需执行的操作。这些头部只对当前经过的单个节点有效，并且不会被转发到下一个节点。常见的 hop-by-hop headers 包括 Connection、Keep-Alive、TE（用于传输编码）、Trailer（用于指定消息尾部的字段列表）等
func removeHopByHopHeaders(header http.Header) {
	hopHeaders := []string{
		"TE",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Proxy-Connection",
	}
	for _, hk := range hopHeaders {
		header.Del(hk)
	}
	connHeader := header.Get("Connection")
	for _, h := range strings.Split(connHeader, ",") {
		h = strings.TrimSpace(h)
		header.Del(h)
	}
	header.Del("Connection")
}
