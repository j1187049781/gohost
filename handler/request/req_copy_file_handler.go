package handler

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/**
如果Request的URL与Pattern匹配，保存这个Requst中的文件到本地
*/
type ReqCopyFileHandler struct {
	Pattern string
	FormFileKeys []string
	patternUrl *url.URL
}

func NewReqCopyFileHandler(pattern string, formFileKeys []string) *ReqCopyFileHandler {
	
	pUrl, err := url.Parse(pattern)
	if err != nil {
		log.Printf("parse pattern url error: %v", err)
		return nil
	}

	return &ReqCopyFileHandler{
		Pattern: pattern,
		patternUrl: pUrl,
		FormFileKeys: formFileKeys,
	}
}

func (h *ReqCopyFileHandler) Match(req *http.Request) bool {
	if h.patternUrl == nil {
		log.Printf("pattern url is nil")
		return false
	}

	if req.URL == nil {
		log.Printf("request url is nil")
		return false
	}

	if !strings.HasPrefix(req.URL.Path, h.patternUrl.Path) {
		return false
	}

	return true
}

func (h *ReqCopyFileHandler) Handle(req *http.Request) {
	if h.patternUrl == nil {
		log.Printf("pattern url url is nil")
		return
	}

	if req.Body == nil {
		log.Printf("request body is nil")
		return
	}

	body := req.Body

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		log.Printf("read request body error: %v", err)
		return
	}
	defer body.Close()

	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	colonedReq := req.Clone(req.Context())
	colonedReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	h.saveFile(colonedReq)
}

func (h *ReqCopyFileHandler) saveFile(req *http.Request) {
	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("parse multipart form error: %v", err)
		return
	}

	for _, key := range h.FormFileKeys {
		file, fileHeader, err := req.FormFile(key)
		if err != nil {
			log.Printf("get form file error: %v", err)
			return
		}
		defer file.Close()
		// get the file name
		fileName := fileHeader.Filename

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			log.Printf("read file error: %v", err)
			return
		}

		md5 := md5.Sum(fileBytes)
		fileName = hex.EncodeToString(md5[:]) + "_" + fileName
		writeFile(fileName, fileBytes)
	}
}

func writeFile(path string, data []byte) {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("create file error: %v", err)
		return
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Printf("write file error: %v", err)
		return
	}
}