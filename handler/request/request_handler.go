package handler

import "net/http"

type RequestHandler interface {

	Match(req *http.Request) bool

	Handle(req *http.Request)
}