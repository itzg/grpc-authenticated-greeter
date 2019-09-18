package common

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
)

func LoadCertPool(caCert string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	file, err := os.Open(caCert)
	if err != nil {
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()

	pemBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	ok := certPool.AppendCertsFromPEM(pemBytes)
	if !ok {
		return nil, errors.New("unable to add CA certs")
	}

	return certPool, nil
}

func LoadClientTlsConfig(caCert string, privateKey string, privateCert string, serverNameOverride string) (*tls.Config, error) {
	certPool, err := LoadCertPool(caCert)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(privateCert, privateKey)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		ServerName:   serverNameOverride,
	}, nil
}
