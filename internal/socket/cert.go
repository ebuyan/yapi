package socket

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func GetCerts(certificate string) (certs *x509.CertPool, err error) {
	certs = x509.NewCertPool()
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		err = errors.New("failed to parse certificate PEM")
		return
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return
	}
	certs.AddCert(cert)
	return
}
