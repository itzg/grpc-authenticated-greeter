package server

import (
	"crypto/tls"
	"github.com/itzg/grpc-authenticated-greeter/certs"
	"github.com/itzg/grpc-authenticated-greeter/protocol"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

func Run(caCert string, privateKey string, privateCert string, binding string) {
	lis, err := net.Listen("tcp", binding)
	if err != nil {
		logrus.WithError(err).Fatal("binding listener")
	}

	cert, err := tls.LoadX509KeyPair(privateCert, privateKey)
	if err != nil {
		logrus.WithError(err).Fatal("loading key pair")
	}

	clientCertPool, err := certs.LoadCertPool(caCert)
	if err != nil {
		logrus.WithError(err).Fatal("loading CA cert")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	})
	if err != nil {
		logrus.WithError(err).Fatal("building transport creds")
	}
	s := grpc.NewServer(grpc.Creds(transportCreds))

	protocol.RegisterHelloServiceServer(s, &server{})

	logrus.Infof("Serving greetings at %s", binding)
	err = s.Serve(lis)
	if err != nil {
		logrus.WithError(err).Fatal("serving grpc")
	}
}
