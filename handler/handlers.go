package handler

import (
	"gohost/config"
	req "gohost/handler/request"
	"net/http"
	"sync"
)

var (
	// RequestHandlers is a slice of RequestHandler
	requestHandlers []req.RequestHandler

	lock sync.RWMutex = sync.RWMutex{}
)


func LoadFromConfig(conf *config.Config) {
	lock.Lock()
	defer lock.Unlock()

	for _, um := range conf.GetMapping() {
		pair := req.NewUrlHandler(um.Pattern,um.Target)
		if pair != nil {
			requestHandlers = append(requestHandlers, pair)
		}
	}
}

func HandleRequest(req *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	for _, h := range requestHandlers {
		if h.Match(req) {
			h.Handle(req) 
			return
		}
	}
}
