// Package main implements a proxy for chat service.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ascarter/grpcproxy/chat"
	"golang.org/x/net/http2"
)

var (
	certFile string
	keyFile  string
	address  string
	origin   string
)

func init() {
	flag.StringVar(&address, "address", ":50050", "listen address")
	flag.StringVar(&origin, "origin", ":50051", "proxy origin")
	flag.StringVar(&certFile, "cert", "cert.pem", "certificate file")
	flag.StringVar(&keyFile, "key", "key.pem", "key file")
	flag.Parse()
}

// newProxy returns a reverse proxy for addr
func newProxy(addr, cert, serverName string) (*httputil.ReverseProxy, error) {
	o := url.URL{Scheme: "https", Host: addr}

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", o.Host)
		req.URL.Scheme = o.Scheme
		req.URL.Host = o.Host
	}

	tlsConfig, err := chat.NewClientTLSConfig(cert, serverName)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	http2.ConfigureTransport(transport)

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}, nil
}

func main() {
	proxy, err := newProxy(origin, certFile, "localhost")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Printf("%q", dump)
		proxy.ServeHTTP(w, r)
	})

	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tlsConfig, err := chat.NewServerTLSConfig(certFile, keyFile, "localhost")
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Handler:   nil,
		TLSConfig: tlsConfig,
	}

	log.Printf("Proxy listening on %v => sending to %v", address, origin)
	log.Fatal(s.ServeTLS(l, "", ""))
	// log.Fatal(http.ListenAndServeTLS(address, certFile, keyFile, nil))
}