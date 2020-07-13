# grpcproxy

Example project with a gRPC client, server, and proxy. Demonstrates a reverse proxy for gRPC requests in Go.

The project includes a server that handles the chat.Greeter and chat.Echo services. A reverse grpc proxy supports forwarding the chat services. Two clients - echo and hello will send an echo request or a hello request.

## Requirements

Install [protobuf compiler](https://github.com/protocolbuffers/protobuf/releases)

+ Ubuntu: `apt install protobuf-compiler-grpc`
+ macOS: `brew install protobuf`

Install Go generator for protobuf and grpc:

```
go get -u google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc
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

By default, server listens on `:50051` and looks for `cert.pem` and `key.pem` in the current directory.

```
./dist/server
```

| Flag | Use |
| ---- | --- |
|`-address string`|listen address (default ":50051")|
|`-cert string`|certificate file (default "cert.pem")
|`-key string`|key file (default "key.pem")|

### Proxy

By default, proxy listens on `:50050` and looks for `cert.pem` and `key.pem` in the current directory for the server side of the proxy. For the orign connection, by default it looks for the same `cert.pem` as the ca cert. The origin cert needs to match the cert the server uses:

```
./dist/proxy
```

| Flag | Use |
| ---- | --- |
|`-address string`|listen address (default ":50050")|
|`-cert string`|certificate file (default "cert.pem")
|`-key string`|key file (default "key.pem")|
|`-cacert string`|ca certificate file (default "cert.pem")
|`-origin string`|proxy origin (default ":50051")|

### client

The client supports two commands:

hello <name>
echo <string>

By default, it will send to server on `:50051`. Use `-address` to use proxy. The client uses `cert.pem` as the ca cert by default.

```
# Send to server
./dist/client hello bob

# Send to proxy
./dist/client -address :50050 echo where are you sally?
```

| Flag | Use |
| ---- | --- |
|`-address string`|server address (default ":50051")|
|`-cacert string`|ca certificate file (default "cert.pem")


## Multiple Certificates

Running `make` will generate a self-signed keypair (`cert.pem` and `key.pem`). If mutliple certs are desired, use the openssl command in the makefile to generate other key pairs and pass them to client, server, or proxy. If you have a local ca, you can supply that cert and sign separate server certificates for the proxy and server.

## Ad-Hoc

Using gRPC server reflection, tools like [grpcurl](https://github.com/fullstorydev/grpcurl) can be used for making ad-hoc requests.

Example:

```
grpcurl -cacert ./cert.pem localhost:50051 list
grpcurl -cacert ./cert.pem -d '{"name": "Bob"}' localhost:50051 chat.Greeter/SayHello
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
