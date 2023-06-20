# grpcproxy

Example project with a gRPC client, server, and proxy. Demonstrates a reverse proxy for gRPC requests in Go.

The project includes a server that handles the chat.Greeter and chat.Echo services. A reverse grpc proxy supports forwarding the chat services. Two clients - echo and hello will send an echo request or a hello request.

This project demonstrates using mTLS for all connections.

## Requirements

Install [protobuf compiler](https://github.com/protocolbuffers/protobuf/releases)

+ Ubuntu: `apt install protobuf-compiler-grpc`
+ macOS: `brew install protobuf`

Install Go generator for protobuf and grpc:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Build

A Makefile provides the following targets:

| Target | Use |
| ------ | --- |
|*build*|Generate pb.go files from proto, a server and proxy keypair, and builds the client, proxy, and server|
|*clean*|Remove generated files|
|*distclean*|Remove everything including keypairs|

## Running

The server and proxy should either run in separate shells or run in the background.

### Server

By default, server listens on `:8051` and looks for certificates in the current directory.

```
./dist/server
```

### Proxy

By default, proxy listens on `:8050`, forwards to `:8051` and looks for certficates in the current directory:

```
./dist/proxy
```

### client

The client supports two commands:

hello <name>
echo <string>

By default, it will send to server on `:8051`. Use `-address` to use proxy. The client uses certificates in the current directory by default.

```
# Send to server
./dist/client hello bob

# Send to proxy
./dist/client -address :8050 echo where are you sally?
```

## Certificates

Running `make keys` will generate a self-signed keypair (`ca.crt` and `ca.key`) for a self-signed CA. Keypairs are created and signed by the CA for server, proxy, and client.

## Ad-Hoc

Using gRPC server reflection, tools like [grpcurl](https://github.com/fullstorydev/grpcurl) can be used for making ad-hoc requests.

Example:

```
grpcurl -cacert ca.crt -cert client.crt -key client.key localhost:8051 list
grpcurl -cacert ca.crt -cert client.crt -key client.key -d '{"name": "Bob"}' localhost:8051 chat.Greeter/SayHello
```

Response is in JSON:

```
{
  "message": "Hello Bob"
}
```

## References

* [protocol-buffers](https://developers.google.com/protocol-buffers/)
* [protobuf-go](https://github.com/protocolbuffers/protobuf-go)
* [grpc-go](https://github.com/grpc/grpc-go)
* [gRPC Go](https://grpc.io/docs/languages/go/)
* [Language guide proto3](https://developers.google.com/protocol-buffers/docs/proto3)
* [RFC 5246 - TLS](https://www.ietf.org/rfc/rfc5246.txt)
* [Azure Self-Signed Certificates](https://docs.microsoft.com/en-us/azure/application-gateway/self-signed-certificates)
* [Azure Configure Mutual TLS Auth](https://docs.microsoft.com/en-us/azure/app-service/app-service-web-configure-tls-mutual-auth)
