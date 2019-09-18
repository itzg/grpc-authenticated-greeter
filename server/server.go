package server

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/itzg/grpc-authenticated-greeter/protocol"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, req *protocol.HelloRequest) (*protocol.HelloResponse, error) {
	p, hasPeer := peer.FromContext(ctx)
	if hasPeer {

		clientSubject := "UNKNOWN"
		if tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			if tlsValue, ok := tlsInfo.GetSecurityValue().(*credentials.TLSChannelzSecurityValue); ok {
				remoteCert, err := x509.ParseCertificate(tlsValue.RemoteCertificate)
				if err != nil {
					logrus.WithError(err).WithField("addr", p.Addr).
						Warn("unable to parse remote certificate")
				} else {
					clientSubject = remoteCert.Subject.CommonName
				}
			}
		}

		return &protocol.HelloResponse{
			Reply: fmt.Sprintf("Hello, %s. You said '%s'", clientSubject, req.Greeting),
		}, nil
	} else {
		return &protocol.HelloResponse{
			Reply: fmt.Sprintf("Hello, mystery caller. You said '%s'", req.Greeting),
		}, nil
	}
}
