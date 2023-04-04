package test

import (
	"crypto/x509"
	"gohost/common/cert"
	"testing"
	"time"
)

func TestGenerateRootPemFile(t *testing.T) {
	certificate := cert.GetRootCert()
	rootCert, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		t.Fatalf("parse root cert failed: %v", err)
	}
	hostName := "test.com"
	cert, err := cert.SignCert(hostName, certificate)
	if err != nil {
		t.Fatalf("sign cert failed: %v", err)
	}
	parseCertificate, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("parse cert failed: %v", err)
	}

	// 1. 验证证书是否有效
	// 2. 验证证书是否是由rootCert签发的
	// 3. 验证证书是否是给hostName签发的
	roots := x509.NewCertPool()
	roots.AddCert(rootCert)
	options := x509.VerifyOptions{
		Roots:       roots,
		CurrentTime: time.Now(),
		DNSName:     hostName,
		KeyUsages:   []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	_, err = parseCertificate.Verify(options)
	if err != nil {
		t.Fatalf("verify cert failed: %v", err)
	}
	t.Log("verify cert success")
}
