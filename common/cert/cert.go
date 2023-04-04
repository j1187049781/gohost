package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"time"
)

// 利用根证书签发子证书
func SignCert(host string, rootCert *tls.Certificate) (*tls.Certificate, error) {

	parent, err := x509.ParseCertificate(rootCert.Certificate[0])
	if err != nil {
		return nil, err
	}
	parentPrivateKey, ok := rootCert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("root cert private key is not rsa")
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
