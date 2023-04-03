package test

import (
	"fmt"
	"gohost/common/cert"
	"testing"
)

func TestGenerateRootPemFile(t *testing.T) {
	// 生成证书
	privateKey, x509Cert, err := cert.GeneratePemFile("go-host", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("privateKey: %v\n", privateKey)
	fmt.Printf("x509Cert: %v\n", x509Cert)
	// 保存证书
	// err = cert.SaveCertFile("ca.pem", x509Cert)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// // 保存私钥
	// err = cert.SaveKeyFile("ca.key", privateKey)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
}
