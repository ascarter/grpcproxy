// Package main implements a proxy for chat service.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ascarter/grpcproxy/chat"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	certFile   string
	keyFile    string
	caCertFile string
	address    string
	origin     string
)

func init() {
	flag.StringVar(&address, "address", ":50050", "listen address")
	flag.StringVar(&origin, "origin", ":50051", "proxy origin")
	flag.StringVar(&certFile, "cert", "cert.pem", "certificate file")
	flag.StringVar(&keyFile, "key", "key.pem", "key file")
	flag.StringVar(&caCertFile, "cacert", "cert.pem", "ca certificate file")
	flag.Parse()
}

// newProxy returns a reverse proxy for addr
func newProxy(addr, cacert string) (*httputil.ReverseProxy, error) {
	// Origin URL
	o := url.URL{Scheme: "https", Host: addr}

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Proto", req.Proto)
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", o.Host)
		req.URL.Scheme = o.Scheme
		req.URL.Host = o.Host
	}

	// Add provided ca root ca's
	crt, err := ioutil.ReadFile(cacert)
	if err != nil {
		return nil, err
	}

	roots, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if !roots.AppendCertsFromPEM(crt) {
		return nil, errors.New("credentials: failed to append ca certificates")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         "localhost",
			RootCAs:            roots,
		},
	}
	http2.ConfigureTransport(transport)

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}, nil
}

// healthServer is used to implement chat.HealthServer
type healthServer struct {
	chat.UnimplementedHealthServer
}

// Status implements chat.HealthServer Status request
func (s *healthServer) Status(context.Context, *chat.StatusRequest) (*chat.StatusReply, error) {
	log.Printf("Received StatusRequest")
	return &chat.StatusReply{Code: http.StatusOK, Message: "OK"}, nil
}

func dumpRequest(r *http.Request) error {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	log.Print(string(dump))
	return nil
}

func main() {
	proxy, err := newProxy(origin, caCertFile)
	if err != nil {
		log.Fatal(err)
	}

	proxyHandle := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Forwarding request for %s", r.RequestURI)
		if err := dumpRequest(r); err != nil {
			log.Print(err)
		}
		proxy.ServeHTTP(w, r)
	}

	// Create a gRPC server for internal messages
	gs := grpc.NewServer()
	chat.RegisterHealthServer(gs, &healthServer{})

	// Enable reflection for grpcurl
	reflection.Register(gs)

	mux := http.NewServeMux()
	mux.HandleFunc("/chat.Greeter/", proxyHandle)
	mux.HandleFunc("/chat.Echo/", proxyHandle)
	mux.Handle("/grpc.reflection.v1alpha.ServerReflection/", gs)
	mux.Handle("/internal.Health/", gs)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			log.Printf("Not found: %v", req.URL)
			http.NotFound(w, req)
			return
		}
		fmt.Fprint(w, "index")
	})

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := &http.Server{Handler: mux}

	log.Printf("Proxy forwarding %v => %v", address, origin)
	log.Fatal(s.ServeTLS(lis, certFile, keyFile))
}
