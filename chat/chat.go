package chat

import (
	"crypto/x509"
	"errors"
	"io/ioutil"
)

// NewCAPool creates a certificate pool with provided ca cert appended
func NewCAPool(cacert string) (*x509.CertPool, error) {
	caCrt, err := ioutil.ReadFile(cacert)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCrt) {
		return nil, errors.New("credentials: failed to append ca certificates")
	}

	return caPool, nil
}
