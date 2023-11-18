package test

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"testing"
)


func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("get form file error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	w.WriteHeader(http.StatusOK)
	
	fileByes, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	md5 := md5.Sum(fileByes)
	fileName := hex.EncodeToString(md5[:])
	log.Printf("md5 name: %s", fileName)
}

// 模拟文件保存服务器
func TestWebServer(t *testing.T) {
	http.HandleFunc("/user/firconfig", handleFileUpload)

	log.Fatal(http.ListenAndServe(":8080", nil))

}



