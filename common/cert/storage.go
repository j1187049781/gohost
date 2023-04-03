package cert

import (
	"crypto/tls"
	"log"
	"path"
)

const (
	dir      = "cert"
	rootHost = "go-host"
)

var (
	cache map[string]*tls.Certificate
)

func init() {
	cache = make(map[string]*tls.Certificate)
	cert, err := LoadRootCert()
	if err != nil {
		log.Fatalf("load root cert failed: %v", err)
	}
	cache[host] = cert
}

func GetSignedCert(host string) (*tls.Certificate, error) {
	cert, ok := cache[host]
	if ok {
		return cert, nil
	}

	// todo: 加锁

	// 判断文件是否存在
	keyPath := path.Join(dir, host+".key")
	certPath := path.Join(dir, host+".crt")
	existKey, err := isExist(keyPath)
	if err != nil {
		return nil, err
	}
	existCert, err := isExist(certPath)
	if err != nil {
		return nil, err
	}
	if !existKey || !existCert {
		err = genAndSave(host, false)
		if err != nil {
			return nil, err
		}
	}
	cert, err := SignCert(host)
	if err != nil {
		return nil, err
	}
	cache[host] = cert
	return cert, nil
}
