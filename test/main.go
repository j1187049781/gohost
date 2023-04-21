package test

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

func main() {
	// 创建HTTP/2代理服务器的TLS配置
	config := &tls.Config{
		NextProtos: []string{"h2", "http/1.1"},
	}

	// 创建HTTP/2代理服务器的TCP监听器
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	// 创建HTTP/2代理服务器的HTTP处理器
	handler := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// 在这里进行请求重定向
			req.URL.Scheme = "https"
			req.URL.Host = req.Host
		},
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}

	// 创建HTTP/2代理服务器并启动监听器
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
		TLSConfig: config,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	server.TLSNextProto["h2"] = func(srv *http.Server, conn *tls.Conn, handler http.Handler) {
		// 因为我们已经使用了ReverseProxy作为处理器，所以可以不做任何事情
	}

	log.Printf("HTTP/2 proxy listening on %s", listener.Addr())

	// 等待连接请求
	err = server.Serve(tls.NewListener(listener, config))
	if err != nil {
		log.Fatal(err)
	}
}



