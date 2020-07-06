// Package main implements a server for chat service.
package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/ascarter/grpcproxy/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	certFile string
	keyFile  string
	address  string
)

func init() {
	flag.StringVar(&address, "address", ":50051", "listen address")
	flag.StringVar(&certFile, "cert", "server_cert.pem", "certificate file")
	flag.StringVar(&keyFile, "key", "server_key.pem", "key file")
	flag.Parse()
}

// server is used to implement chat.GreeterServer.
type greeterServer struct {
	chat.UnimplementedGreeterServer
}

// echoServer is used to implement a chat.EchoServer
type echoServer struct {
	chat.UnimplementedEchoServer
}

// SayHello implements chat.GreeterServer
func (s *greeterServer) SayHello(ctx context.Context, in *chat.HelloRequest) (*chat.HelloReply, error) {
	log.Printf("Received HelloRequest: %v", in.GetName())
	return &chat.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// Replay impelments chat.EchoServer
func (s *echoServer) Replay(ctx context.Context, in *chat.EchoRequest) (*chat.EchoReply, error) {
	log.Printf("Received EchoRequest: %v", in.GetMessage())
	return &chat.EchoReply{Message: in.GetMessage()}, nil
}

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tlsConfig, err := chat.NewServerTLSConfig(certFile, keyFile, "localhost")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	chat.RegisterGreeterServer(s, &greeterServer{})
	chat.RegisterEchoServer(s, &echoServer{})

	log.Printf("Listening on %v", address)
	log.Fatal(s.Serve(lis))
}
