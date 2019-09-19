package client

import (
	"context"
	"github.com/itzg/grpc-authenticated-greeter/certs"
	"github.com/itzg/grpc-authenticated-greeter/protocol"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"time"
)

func Run(caCert string, privateKey string, privateCert string, serverAddress string, serverName string, message string) {

	if serverName == "" {
		host, _, err := net.SplitHostPort(serverAddress)
		if err != nil {
			logrus.WithError(err).Fatal("Unable to split server address")
		}

		serverName = host
	}

	tlsConfig, err := certs.LoadClientTlsConfig(caCert, privateKey, privateCert, serverName)
	if err != nil {
		logrus.WithError(err).Fatal("loading client tls config")
	}

	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		logrus.WithError(err).Fatal("connecting")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client := protocol.NewHelloServiceClient(conn)
	response, err := client.SayHello(ctx, &protocol.HelloRequest{
		Greeting: message,
	})
	if err != nil {
		logrus.WithError(err).Warn("hello failed")
	} else {
		logrus.WithField("response", response.Reply).Info("got response")
	}
}
