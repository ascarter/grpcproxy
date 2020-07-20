GOROOT ?= $(shell go env GOROOT)

DIST    := dist
OBJS    := $(addprefix $(DIST)/,$(notdir $(wildcard cmd/*)))
PB_SRCS := chat/*.proto
PB_OBJS := chat/*.pb.go
CERTS   := ca.crt client.crt server.crt proxy.crt
KEYS    := $(CERTS:%.crt=%.key)
CSRS    := $(CERTS:%.crt=%.csr)
SRLS    := ca.srl

# OPENSSL_ARGS := -newkey rsa:2048 -new -nodes -x509 -days 3650 -subj "/C=US/ST=Washington/L=Snoqualmie/O=$(USER)/OU=Development/CN=localhost"

$(DIST)/%: ./cmd/% cmd/%/*.go | $(DIST)
	go build -o $@ ./$<

$(DIST):
	mkdir -p $(DIST)

%.pem:
	openssl req $(OPENSSL_ARGS) -out $@ -keyout $(@:%cert.pem=%key.pem)
	openssl x509 -noout -text -in $@

$(PB_OBJS): $(PB_SRCS)
	protoc -I chat/ $^ --go_out=chat --go_opt=paths=source_relative --go-grpc_out=chat --go-grpc_opt=paths=source_relative

.DEFAULT_GOAL = all

all: build

build: keys $(PB_OBJS) $(OBJS)

%.key:
	openssl ecparam -out $@ -name prime256v1 -genkey

%.csr: %.key
	openssl req -new -sha256 -key $^ -out $@ -subj "/C=US/ST=Washington/L=Snoqualmie/O=$(USER)/OU=Development/CN=localhost"

%.crt: %.csr
	openssl x509 -req -sha256 -days 365 -in $^ -CA ca.crt -CAkey ca.key -CAcreateserial -out $@
	openssl x509 -in $@ -text -noout

ca.crt: ca.csr
	openssl x509 -req -sha256 -days 365 -in ca.csr -signkey ca.key -out ca.crt

keys: $(KEYS) $(CSRS) $(CERTS)

# openssl ecparam -out ca.key -name prime256v1 -genkey
# openssl req -new -sha256 -key ca.key -out ca.csr -subj "/C=US/ST=Washington/L=Snoqualmie/O=$(USER)/OU=Development/CN=localhost"
# openssl x509 -req -sha256 -days 365 -in ca.csr -signkey ca.key -out ca.crt

# openssl ecparam -out server.key -name prime256v1 -genkey
# openssl req -new -sha256 -key server.key -out server.csr -subj "/C=US/ST=Washington/L=Snoqualmie/O=$(USER)/OU=Development/CN=localhost"
# openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256
# openssl x509 -in server.crt -text -noout

test: $(PB_OBJS) $(PB_SRCS)
	go test ./...

clean:
	-rm $(PB_OBJS) $(OBJS)

distclean: clean
	-rm $(CERTS) $(CSRS) $(KEYS) $(SRLS)
	-rm -rf $(DIST)
