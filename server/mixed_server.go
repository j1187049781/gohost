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
	connProxy chan ConnProxy
)

func init() {
	connProxy = make(chan ConnProxy, 2000)
}

func Setup(conf *config.Config) {
	addrPort := fmt.Sprintf("%s:%d", conf.ServerConfig.ListenAddr, conf.ServerConfig.ListenPort)
	ln, err := net.Listen(conf.ServerConfig.Network, addrPort)
	if err != nil {
		log.Fatalf("启动监听失败: %s", err.Error())

	}
	log.Printf("服务启动成功： %s", addrPort)

	go proxy()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("接受请求失败: %s", err.Error())
		}
		go handleConnWithProtocol(conn)
	}
}

func proxy() {
	client := &http.Client{}

	for c := range connProxy {
		p := c
		go func() {
			defer p.Conn.Close()

			for {
				reader := N.NewReader(p.Conn)
				req, err := http.ReadRequest(reader.Reader())
				if err != nil {
					fmt.Printf("req: %v\n, %s", req, err.Error())
					return
				}

				// 读取Header中的keepAlive，判断是是否需要keepAlive
				if p.Protocol == "http" {
					p.KeepAlive = strings.EqualFold(req.Header.Get("Connection"), "keep-alive")
				}

				// 移除http请求头中的hop-by-hop headers
				removeHopByHopHeaders(req.Header)

				// reuse existing http.Request
				req.RequestURI = ""
				req.URL.Scheme = p.Protocol
				req.URL.Host = req.Host

				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("resp: %v, %s", resp, err)
					return
				}

				removeHopByHopHeaders(resp.Header)
				if p.KeepAlive {
					resp.Header.Set("Connection", "keep-alive")
					resp.Header.Set("Proxy-Connection", "keep-alive")
					resp.Header.Set("Keep-Alive", "timeout=5, max=1000")
				}
				// 不需要resp.Body.Close
				// The Response Body is closed after it is sent
				if err := resp.Write(p.Conn); err != nil {
					fmt.Printf("resp reply error: %v, %s", resp, err)
					return
				}
				if !p.KeepAlive {
					return
				}
			}
		}()
	}
}

func handleConnWithProtocol(conn net.Conn) {
	//todo: 判断代理协议,目前支持Http

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
		log.Printf("接受https代理请求: %s --> %s", bConn.RemoteAddr().String(), bConn.LocalAddr().String())

		tlsConn := tls.Server(bConn, &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return cert.GetSignedCert(info.ServerName)
			},
		})
		err = tlsConn.Handshake()
		if err != nil {
			fmt.Printf("tls handshake error: %s", err.Error())
			return
		}
		log.Printf("https连接建立成功: %s --> %s", bConn.RemoteAddr().String(), bConn.LocalAddr().String())

		keepAlive := strings.EqualFold(req.Header.Get("Proxy-Connection"), "keep-alive")
		connProxy <- ConnProxy{tlsConn, "https", keepAlive}

	} else {
		// http 请求
		connProxy <- ConnProxy{bConn, "https", true}
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
