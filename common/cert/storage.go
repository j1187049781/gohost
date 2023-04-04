package cert

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"path"
	"sync"
)

const (
	dir      = "cert"
	rootHost = "go-host"
)

var (
	cache map[string]*tls.Certificate
	lock  sync.Mutex
)

func init() {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("mkdir  failed: %v", err)
	}
	cache = make(map[string]*tls.Certificate)
	cert, err := loadRootCert()
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

func loadRootCert() (*tls.Certificate, error) {
	// 判断文件是否存在
	keyPath := path.Join(dir, rootHost+".key")
	certPath := path.Join(dir, rootHost+".crt")
	existKey, err := isExist(keyPath)
	if err != nil {
		return nil, err
	}
	existCert, err := isExist(certPath)
	if err != nil {
		return nil, err
	}
	if !existKey || !existCert {
		privateKey, cert, err := GeneratePemFile(rootHost, nil, nil)
		if err != nil {
			return nil, err
		}
		// 保存
		err = save("PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey), keyPath)
		if err != nil {
			return nil, err
		}
		err = save("CERTIFICATE", cert, certPath)
		if err != nil {
			return nil, err
		}

	}

	// 读取证书
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// 判断文件是否存在
func isExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 保存文件
func save(typeInfo string, data []byte, fileName string) (err error) {

	keyFd, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer keyFd.Close()

	err = pem.Encode(keyFd, &pem.Block{
		Type:  typeInfo,
		Bytes: data,
	})
	if err != nil {
		return
	}
	return nil
}
