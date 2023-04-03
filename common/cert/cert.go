package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path"
	"time"
)

const (
	dir      = "cert"
	rootHost = "go-host"
)

// 利用根证书签发子证书
func SignCert(host string, rootCert *tls.Certificate) (*tls.Certificate, error) {

	parent, err := x509.ParseCertificate(rootCert.Certificate[0])
	if err != nil {
		return nil, err
	}
	parentPrivateKey, err := x509.ParsePKCS1PrivateKey(rootCert.PrivateKey.(*rsa.PrivateKey).D.Bytes())
	if err != nil {
		return nil, err
	}

	privateKey, cert, err := GeneratePemFile(host, parent, parentPrivateKey)
	if err != nil {
		return nil, err
	}
	ret := &tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  privateKey,
	}
	return ret, nil
}

func LoadRootCert() (*tls.Certificate, error) {
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

// GeneratePemFile 生成证书,私钥和公钥
func GeneratePemFile(host string, parent *x509.Certificate, parentPrivateKey *rsa.PrivateKey) (*rsa.PrivateKey, []byte, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	temp := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         host,
			Country:            []string{"CN"},         // 证书所属的国家
			Organization:       []string{host},         // 证书存放的公司名称
			OrganizationalUnit: []string{"department"}, // 证书所属的部门名称
			Province:           []string{"cd"},         // 证书签发机构所在省
			Locality:           []string{"cd"},         // 证书签发机构所在市
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(100 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		Issuer: pkix.Name{
			CommonName: host,
		},
		DNSNames: []string{host},
	}
	if parent == nil && parentPrivateKey == nil {
		parent = &temp
		parentPrivateKey = privateKey
	}

	x509Cert, err := x509.CreateCertificate(rand.Reader, &temp, parent, &privateKey.PublicKey, parentPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, x509Cert, nil
}
