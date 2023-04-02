package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// golang生成一个Ca 证书,私钥和公钥
func GenerateRootPemFile() (*rsa.PrivateKey, []byte, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	temp := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         "gohost",
			Country:            []string{"CN"},         // 证书所属的国家
			Organization:       []string{"gohost"},     // 证书存放的公司名称
			OrganizationalUnit: []string{"department"}, // 证书所属的部门名称
			Province:           []string{"cd"},         // 证书签发机构所在省
			Locality:           []string{"cd"},         // 证书签发机构所在市
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		Issuer: pkix.Name{
			CommonName: "gohost",
		},
		DNSNames: []string{"gohost"},
	}
	x509Cert, err := x509.CreateCertificate(rand.Reader, &temp, &temp, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, x509Cert, nil
}
