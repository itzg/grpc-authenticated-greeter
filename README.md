## Building

```shell script
go build
```

## Example usage

Generate CA, server, and client certs
```shell script
./grpc-authenticated-greeter gencerts
```

Start the server on port 7676:
```shell script
./grpc-authenticated-greeter server \
  --ca ca_cert.pem --cert server_cert.pem --key server_key.pem \
  --binding :7676
```

In another terminal, run a client:
```shell script
./grpc-authenticated-greeter client \
  --ca ca_cert.pem --cert client1_cert.pem --key client1_key.pem \
  --serveraddress 127.0.0.1:7676 --servername server \
  --message "Read me"
```

The client should log the response from the server, such as:
```
INFO[0000] got response   response="Hello, client1. You said 'Read me'"
```

## What I learned

- [go-arg](https://github.com/alexflint/go-arg) is very cool! Just feed it a struct and it'll parse command line arguments. It's very flexible and intuitive. Even embedding common arguments into a "command struct" did what I expected
```go
type ClientServerArgs struct {
	Ca   string `arg:"required" help:"PEM file containing the CA cert shared by server and clients"`
	Key  string `arg:"required" help:"PEM file containing private key"`
	Cert string `arg:"required" help:"PEM file containing public certificate"`
}

type ServerCmd struct {
	ClientServerArgs
	Binding string `arg:"required" help:"host:port of server binding where host is optional"`
}
```
- Generating self-signed certs and then signing client and server certs is quite doable in Go. Check out [certs/generate.go](certs/generate.go) to see how that's done
- The authenticating client info is a little bit buried when implementing a server handler, but [peer.FromContext](https://godoc.org/google.golang.org/grpc/peer#FromContext) was the key to cracking that open. [server/server.go](server/server.go) is where that is used.
- To do full [mTLS authentication](https://grpc.io/docs/guides/auth/, be sure to configure the server's TLS to require and verify the client cert:
```go
	transportCreds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	})
```