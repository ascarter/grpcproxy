// Package main implements a server for chat service.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/ascarter/grpcproxy/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	certFile string
	keyFile  string
	address  string
)

func init() {
	flag.StringVar(&address, "address", ":50051", "listen address")
	flag.StringVar(&certFile, "cert", "server.crt", "certificate file")
	flag.StringVar(&keyFile, "key", "server.key", "key file")
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

// LotsOfReplices repeaths greeting 5 times
func (s *greeterServer) LotsOfReplies(in *chat.HelloRequest, streamRes chat.Greeter_LotsOfRepliesServer) error {
	log.Printf("Received LotsOfReplies for HelloRequest: %v", in.GetName())

	// Return greeting 5 times
	for i := 0; i < 5; i++ {
		log.Printf("--> sending reply %d for %v", i, in.GetName())
		res := chat.HelloReply{Message: fmt.Sprintf("Hello %v - %d", in.GetName(), i)}
		streamRes.Send(&res)
	}
	return nil
}

// ManyHellos sends hello to many names
func (s *greeterServer) ManyHellos(in chat.Greeter_ManyHellosServer) error {
	log.Printf("Received ManyHellos for HelloRequest")
	for {
		req, err := in.Recv()
		if err == io.EOF {
			// end of receiving requests
			log.Print("--> end receiving requests and send response")
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("--> sending reply for %v", req.GetName())
		in.Send(&chat.HelloReply{Message: "Hello " + req.GetName()})
	}
	return nil
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

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	chat.RegisterGreeterServer(s, &greeterServer{})
	chat.RegisterEchoServer(s, &echoServer{})

	// Enable reflection for grpcurl
	reflection.Register(s)

	log.Printf("Listening on %v", address)
	log.Fatal(s.Serve(lis))
}
