# grpcproxy

Example project with a gRPC client, server, and proxy. Demonstrates a reverse proxy for gRPC requests in Go.


## Requirements

Install [protobuf compiler](https://github.com/protocolbuffers/protobuf/releases)

+ Ubuntu: `apt install protobuf-compiler-grpc`
+ macOS: `brew install protobuf`

Install Go generator for protobuf and grpc:

```
go get -u google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc
```