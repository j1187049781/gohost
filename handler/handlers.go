package handler

import (
	"gohost/config"
	req "gohost/handler/request"
	"net/http"
	"sync"
)

type Handlers struct{
	// RequestHandlers is a slice of RequestHandler
	requestHandlers []req.RequestHandler

	lock sync.RWMutex
}


func (h *Handlers) LoadFromConfig(conf *config.Config) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for _, um := range conf.GetMapping() {
		pair := req.NewUrlHandler(um.Pattern,um.Target)
		if pair != nil {
			h.requestHandlers = append(h.requestHandlers, pair)
		}
	}
}

func (h *Handlers) HandleRequest(req *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for _, h := range h.requestHandlers {
		if h.Match(req) {
			h.Handle(req) 
			return
		}
	}
}
