package cert

import (
	"crypto/tls"
	"log"
	"sync"
)

var (
	cache map[string]*tls.Certificate
	lock  sync.Mutex
)

func init() {
	cache = make(map[string]*tls.Certificate)
	cert, err := LoadRootCert()
	if err != nil {
		log.Fatalf("load root cert failed: %v", err)
	}
	cache[rootHost] = cert
}

func GetRootCert() *tls.Certificate {
	lock.Lock()
	defer lock.Unlock()
	return cache[rootHost]
}

func GetSignedCert(host string) (*tls.Certificate, error) {
	lock.Lock()
	defer lock.Unlock()
	cert, ok := cache[host]
	if ok {
		return cert, nil
	}

	cert, err := SignCert(host, cache[rootHost])
	if err != nil {
		return nil, err
	}
	cache[host] = cert
	return cert, nil
}


