// Package main implements a client for chat service.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"strings"
	"sync"
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
	caCertFile     string
	clientCertFile string
	clientKeyFile  string
	address        string
	insecure       bool
	commands       map[string]Command
)

func init() {
	flag.StringVar(&address, "address", ":50051", "server address")
	flag.StringVar(&caCertFile, "cacert", "certificates/ca.crt", "ca certificate file")
	flag.StringVar(&clientCertFile, "clientcert", "certificates/client.crt", "client certificate file")
	flag.StringVar(&clientKeyFile, "clientkey", "certificates/client.key", "client key file")
	flag.BoolVar(&insecure, "insecure", false, "connect insecure")

	flag.Parse()

	commands = map[string]Command{
		"hello":  hello,
		"lots":   lots,
		"many":   many,
		"echo":   echo,
		"status": status,
	}
}

// A Command is the API for a sub-command
type Command func(*grpc.ClientConn, []string) error

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

func lots(conn *grpc.ClientConn, args []string) error {
	c := chat.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(args) > 0 {
		name = strings.Join(args, " ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.LotsOfReplies(ctx, &chat.HelloRequest{Name: name})
	if err != nil {
		return err
	}

	for {
		m, err := r.Recv()
		if err == io.EOF {
			// end of stream
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Greeting: %s", m.GetMessage())
	}

	return nil
}

func many(conn *grpc.ClientConn, args []string) error {
	c := chat.NewGreeterClient(conn)

	if len(args) < 1 {
		args = append(args, defaultName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.ManyHellos(ctx)
	if err != nil {
		return err
	}

	// Send requests and receive responses concurrently
	var wg sync.WaitGroup
	wg.Add(2)

	// Sending requests
	go func() {
		defer wg.Done()
		for _, name := range args {
			stream.Send(&chat.HelloRequest{Name: name})
		}
		stream.CloseSend()
	}()

	// Receive messages
	go func() {
		defer wg.Done()
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Greeting: %s", res.GetMessage())
		}
	}()

	wg.Wait()
	log.Print("Streaming terminated")
	return nil
}

func status(conn *grpc.ClientConn, args []string) error {
	c := chat.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Status(ctx, &chat.StatusRequest{})
	if err != nil {
		return err
	}

	log.Printf("Status: %d %s", r.GetCode(), r.GetMessage())
	return nil
}

func main() {
	// Verify subcommand provided
	if len(flag.Args()) < 1 {
		log.Fatal("command required (hello|echo|status)")
	}

	var (
		tlsConfig      *tls.Config
		grpcDialOption grpc.DialOption
	)

	if insecure {
		log.Print("running insecure!")
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "localhost",
		}
		grpcDialOption = grpc.WithInsecure()
	} else {
		// Add provided ca to root ca's
		caPool, err := chat.NewCAPool(caCertFile)
		if err != nil {
			log.Fatal(err)
		}

		// Configure client TLS
		tlsConfig = &tls.Config{
			RootCAs:    caPool,
			ServerName: "localhost",
		}

		// Check if client provided certificate/key
		if clientCertFile != "" && clientKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
			if err != nil {
				log.Fatal(err)
			}

			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		// Set up a connection to the server
		creds := credentials.NewTLS(tlsConfig)
		grpcDialOption = grpc.WithTransportCredentials(creds)
	}

	log.Printf("Connection %v", address)
	conn, err := grpc.Dial(address, grpcDialOption)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Execute command
	if cmd, ok := commands[flag.Args()[0]]; ok {
		if err := cmd(conn, flag.Args()[1:]); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("invalid command %s", flag.Args()[0])
	}
}
