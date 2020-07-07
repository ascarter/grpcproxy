// Package main implements a client for chat service.
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

const (
	defaultMessage = "Is anyone out there?"
	defaultName    = "World"
)

var (
	caCertFile string
	address    string
)

func init() {
	flag.StringVar(&address, "address", ":50051", "server address")
	flag.StringVar(&caCertFile, "cacert", "cert.pem", "ca certificate file")
	flag.Parse()
}

func echo(conn *grpc.ClientConn, args []string) error {
	c := chat.NewEchoClient(conn)

	// Contact the server and print out its response
	msg := defaultMessage
	if len(args) > 0 {
		msg = strings.Join(args, " ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Replay(ctx, &chat.EchoRequest{Message: msg})
	if err != nil {
		return err
	}

	log.Printf("Echo: %s", r.GetMessage())
	return nil
}

func hello(conn *grpc.ClientConn, args []string) error {
	c := chat.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(args) > 0 {
		name = strings.Join(args, " ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &chat.HelloRequest{Name: name})
	if err != nil {
		return err
	}

	log.Printf("Greeting: %s", r.GetMessage())
	return nil
}

func main() {
	// Verify subcommand provided
	if len(flag.Args()) < 1 {
		log.Fatal("command required (hello|echo)")
	}

	// Set up a connection to the server
	creds, err := credentials.NewClientTLSFromFile(caCertFile, "localhost")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connection %v", address)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Execute command
	switch flag.Args()[0] {
	case "hello":
		if err := hello(conn, flag.Args()[1:]); err != nil {
			log.Fatal(err)
		}
	case "echo":
		if err := echo(conn, flag.Args()[1:]); err != nil {
			log.Fatal(err)
		}
	}
}
