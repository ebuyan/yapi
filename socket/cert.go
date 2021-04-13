package socket

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func GetCerts(certificate string) (*x509.CertPool, error) {
	certs := x509.NewCertPool()
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse certificate: " + err.Error())
	}
	certs.AddCert(cert)
	return certs, nil
}
