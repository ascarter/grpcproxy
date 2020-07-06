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
|*build*|Generate pb.go files from proto, a server and proxy keypair, and build all clients, proxy, and server|
|*clean*|Remove generated files|
|*distclean*|Remove everything including keypairs|

## Running

The server and proxy should either run in separate shells or run in the background.

### Server

By default, server listens on `:50051` and looks for `server_cert.pem` and `server_key.pem` in the current directory.

```
./dist/server
```

| Flag | Use |
| ---- | --- |
|`-address string`|listen address (default ":50051")|
|`-cert string`|certificate file (default "server_cert.pem")
|`-key string`|key file (default "server_key.pem")|

### Proxy

By default, proxy listens on `:50050` and looks for `proxy_cert.pem` and `proxy_key.pem` in the current directory:

```
./dist/server
```

| Flag | Use |
| ---- | --- |
|`-address string`|listen address (default ":50050")|
|` -cert string`|certificate file (default "proxy_cert.pem")
|`-key string`|key file (default "proxy_key.pem")|
|`-origincert string`|origin certificate file (default "server_cert.pem")|
|`-origin string`|proxy origin (default ":50051")|

### hello

Send a name and receive a greeting. By default, it will send to server. Use `-address` to use proxy

```
# Send to server
./dist/hello bob

# Send to proxy
./dist/hello -address :50050 sally
```

| Flag | Use |
| ---- | --- |
|`-address string`|server address (default ":50051")|
|` -cert string`|certificate file (default "cert.pem")

### echo

Send a message and receive it back. By default, it will send to server. Use `-address` to use proxy

```
# Send to server
./dist/echo "Echo me"

# Send to proxy
./dist/echo -address :50050 "Echo me via proxy"
```

| Flag | Use |
| ---- | --- |
|`-address string`|server address (default ":50051")|
|` -cert string`|certificate file (default "cert.pem")

## Example

```
make
cd dist
./server &
./proxy &
./echo -cert proxy_cert.pem -address :50050
```

## References

* [protocol-buffers](https://developers.google.com/protocol-buffers/)
* [protobuf-go](https://github.com/protocolbuffers/protobuf-go)
* [grpc-go](https://github.com/grpc/grpc-go)
* [gRPC Go](https://grpc.io/docs/languages/go/)
* [Language guide proto3](https://developers.google.com/protocol-buffers/docs/proto3)
