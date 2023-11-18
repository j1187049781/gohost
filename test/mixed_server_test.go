package test

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestMixedServerProxyHttps(t *testing.T) {
	proxyStr := "http://127.0.0.1:8888"
    proxyURL, err := url.Parse(proxyStr)
    if err != nil {
        panic(err)
    }

    // 设置代理
    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
    }
	req, _ := http.NewRequest("GET", "https://blog.csdn.net/weixin_43314519/article/details/119899888", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	t.Log(resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b) 
	s = s[len(s)-100:]
	log.Println(s)
	
}