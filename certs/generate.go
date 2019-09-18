package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/big"
	"os"
	"time"
)

func Generate() {
	serialNumber := big.NewInt(1)

	caCert, caKey, err := generateCaCert(serialNumber)
	if err != nil {
		logrus.WithError(err).Fatal("generating CA key and cert")
	}

	serialNumber.Add(serialNumber, big.NewInt(1))
	err = generateCert(caCert, caKey, serialNumber, x509.ExtKeyUsageServerAuth, "server")
	if err != nil {
		logrus.WithError(err).Fatal("generating server key and cert")
	}

	serialNumber.Add(serialNumber, big.NewInt(1))
	err = generateCert(caCert, caKey, serialNumber, x509.ExtKeyUsageClientAuth, "client1")
	if err != nil {
		logrus.WithError(err).Fatal("generating client key and cert")
	}

	serialNumber.Add(serialNumber, big.NewInt(1))
	err = generateCert(caCert, caKey, serialNumber, x509.ExtKeyUsageClientAuth, "client2")
	if err != nil {
		logrus.WithError(err).Fatal("generating client key and cert")
	}
}

func generateCert(caCert *x509.Certificate, caKey interface{}, serialNumber *big.Int, usage x509.ExtKeyUsage, subjectCn string) error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating client key: %w", err)
	}

	certTemplate := &x509.Certificate{
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{usage},
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(time.Hour * 24 * 365),
		AuthorityKeyId:     caCert.SubjectKeyId,
		SignatureAlgorithm: x509.SHA512WithRSA,
		SerialNumber:       serialNumber,
		Subject: pkix.Name{
			CommonName: subjectCn,
		},
	}

	clientCertDer, err := x509.CreateCertificate(rand.Reader, certTemplate, caCert, key.Public(), caKey)
	if err != nil {
		return fmt.Errorf("creating client cert: %w", err)
	}

	err = writeDerToPemFile(clientCertDer, "CERTIFICATE", fmt.Sprintf("%s_cert.pem", subjectCn))
	if err != nil {
		return fmt.Errorf("writing client cert: %w", err)
	}

	keyDer, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return fmt.Errorf("marshaling private key: %w", err)
	}
	err = writeDerToPemFile(keyDer, "PRIVATE KEY", fmt.Sprintf("%s_key.pem", subjectCn))
	if err != nil {
		return fmt.Errorf("writing private key file: %w", err)
	}

	return nil
}

func generateCaCert(serialNumber *big.Int) (*x509.Certificate, interface{}, error) {
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, nil, fmt.Errorf("generating CA key: %w", err)
	}
	certTemplate := &x509.Certificate{
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		SignatureAlgorithm:    x509.SHA512WithRSA,
		SerialNumber:          serialNumber,
		Subject: pkix.Name{
			CommonName: "ca",
		},
	}
	caCertDer, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, caKey.Public(), caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("creating CA cert: %w", err)
	}
	err = writeDerToPemFile(caCertDer, "CERTIFICATE", "ca_cert.pem")
	if err != nil {
		return nil, nil, fmt.Errorf("writing CA cert file: %w", err)
	}
	caKeyDer, err := x509.MarshalPKCS8PrivateKey(caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("marshaling CA key: %w", err)
	}
	err = writeDerToPemFile(caKeyDer, "PRIVATE KEY", "ca_key.pem")
	if err != nil {
		return nil, nil, fmt.Errorf("writing CA key file: %w", err)
	}

	caCert, err := x509.ParseCertificate(caCertDer)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing CA cert: %w", err)
	}

	return caCert, caKey, nil
}

func writeDerToPemFile(derBytes []byte, pemType string, filename string) error {
	logrus.Infof("Writing to %s", filename)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()
	err = pem.Encode(file, &pem.Block{
		Type:  pemType,
		Bytes: derBytes,
	})
	if err != nil {
		return fmt.Errorf("encoding pem: %w", err)
	}

	return nil
}
