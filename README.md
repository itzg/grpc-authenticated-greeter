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