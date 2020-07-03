package chat

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

// NewClientTLSConfig returns TLS config for cert for client requests
func NewClientTLSConfig(certPath, serverName string) (*tls.Config, error) {
	crt, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(crt)

	return &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false,
		ServerName:         serverName,
	}, nil
}

// NewServerTLSConfig returns TLS config for cert and key for server
func NewServerTLSConfig(certPath, keyPath, serverName string) (*tls.Config, error) {
	crt, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
	}, nil
}
