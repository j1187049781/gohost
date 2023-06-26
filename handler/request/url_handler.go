package handler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)
/**
 如果Request的URL与Pattern匹配，修改Request的URL的Scheme，Host，Path前缀为Target的Scheme，Host， Path前缀
*/
type UrlHandler struct {
	Pattern string
	Target  string
	patternUrl *url.URL
	targetUrl  *url.URL
}

func NewUrlHandler(pattern, target string) *UrlHandler {

	pUrl, err := url.Parse(pattern)
	if err != nil {
		log.Printf("parse pattern url error: %v", err)
	}

	tUrl, err := url.Parse(target)
	if err != nil {
		log.Printf("parse target url error: %v", err)
	}

	return &UrlHandler{
		Pattern: pattern,
		Target:  target,
		patternUrl: pUrl,
		targetUrl:  tUrl,
	}
}

func (h *UrlHandler) Match(req *http.Request) bool {
	if h.patternUrl == nil {
		log.Printf("pattern url is nil")
		return false
	}

	if req.URL == nil {
		log.Printf("request url is nil")
		return false
	}

	if req.URL.Host != h.patternUrl.Host {
		return false
	}

	if !strings.HasPrefix(req.URL.Path, h.patternUrl.Path) {
		return false
	}

	return true
}

func (h *UrlHandler) Handle(req *http.Request) {
	if h.patternUrl == nil || h.targetUrl == nil {
		log.Printf("pattern url or target url is nil")
		return
	}

	//url: [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	url := h.targetUrl
	if url.Scheme != "" {
		req.URL.Scheme = url.Scheme
	}
	
	if url.Host != "" {
		req.URL.Host = url.Host
	}

	if url.Path != "" {
		req.URL.Path = strings.Replace(req.URL.Path, h.patternUrl.Path, url.Path, 1)
	}
	
}