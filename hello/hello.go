// Package main implements a greeter client for chat service.
package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/ascarter/grpcproxy/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const defaultName = "world"

var (
	certFile string
	keyFile  string
	address  string
)

func init() {
	flag.StringVar(&address, "address", ":50051", "server address")
	flag.StringVar(&certFile, "cert", "cert.pem", "certificate file")
	flag.Parse()
}

func main() {
	tlsConfig, err := chat.NewClientTLSConfig(certFile, "localhost")
	if err != nil {
		log.Fatal(err)
	}

	// Set up a connection to the server.
	log.Printf("Connection %v", address)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := chat.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(flag.Args()) > 0 {
		name = strings.Join(flag.Args(), " ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &chat.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}
