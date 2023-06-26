package handler

import (
	req "gohost/handler/request"
	"net/http"
)

var (
	// RequestHandlers is a slice of RequestHandler
	RequestHandlers []req.RequestHandler
)

func init() {
	RequestHandlers = make([]req.RequestHandler, 0)

	// add your handler here
	RequestHandlers = append(RequestHandlers, req.NewUrlHandler("//test-galaxy.hzins.com/api/rde","//127.0.0.1:8080/api/rde"))
}

func HandleRequest(req *http.Request) {
	for _, h := range RequestHandlers {
		if h.Match(req) {
			h.Handle(req)
			return
		}
	}
}
