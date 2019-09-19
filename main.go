package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/itzg/grpc-authenticated-greeter/certs"
	"github.com/itzg/grpc-authenticated-greeter/client"
	"github.com/itzg/grpc-authenticated-greeter/server"
	"os"
)

// injected by release build
var (
	version string
	commit  string
	date    string
)

type ClientServerArgs struct {
	Ca   string `arg:"required" help:"PEM file containing the CA cert shared by server and clients"`
	Key  string `arg:"required" help:"PEM file containing private key"`
	Cert string `arg:"required" help:"PEM file containing public certificate"`
}

type ServerCmd struct {
	ClientServerArgs
	Binding string `arg:"required" help:"host:port of server binding where host is optional"`
}

type ClientCmd struct {
	ClientServerArgs
	ServerAddress string `arg:"required" help:"host:port of the gRPC server"`
	ServerName    string `help:"SNI name to use when contacting the server. If not set, host from --serveraddress is used"`
	Message       string `arg:"required" help:"Any message you want to send to the server"`
}

type GenCerts struct {
	// no args needed
}

type args struct {
	Client   *ClientCmd `arg:"subcommand:client" help:"Runs the gRPC client and sends authenticated hello request"`
	Server   *ServerCmd `arg:"subcommand:server" help:"Runs the gRPC server"`
	GenCerts *GenCerts  `arg:"subcommand:gencerts" help:"Generate CA, server, and client certs for testing"`
}

func (args) Version() string {
	return fmt.Sprintf("grpc-authenticated-greeter %s (%s @ %s)", version, commit, date)
}

func main() {
	var args args
	parser := arg.MustParse(&args)

	switch {
	case args.Client != nil:
		client.Run(args.Client.Ca, args.Client.Key, args.Client.Cert, args.Client.ServerAddress,
			args.Client.ServerName, args.Client.Message)

	case args.Server != nil:
		server.Run(args.Server.Ca, args.Server.Key, args.Server.Cert, args.Server.Binding)

	case args.GenCerts != nil:
		certs.Generate()

	default:
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
}
