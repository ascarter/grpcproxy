GOROOT ?= $(shell go env GOROOT)

DIST      := dist
OBJS      := $(addprefix $(DIST)/,$(notdir $(wildcard cmd/*)))
PB_SRCS   := chat/*.proto
PB_OBJS   := chat/*.pb.go

CERT_ROOT := $(DIST)/certificates
CSR_CONF  := csr.conf
SRLS      := $(CERT_ROOT)/ca.srl
CA_CRT	  := $(CERT_ROOT)/ca.crt
CA_KEY	  := $(CA_CRT:%.crt=%.key)
CRTS      := $(patsubst %.crt,$(CERT_ROOT)/%.crt,client.crt server.crt proxy.crt)
KEYS      := $(CRTS:%.crt=%.key)
PFXS      := $(CRTS:%.crt=%.pfx)

.PRECIOUS: %.key

.DEFAULT_GOAL = all

all: build

build: certificates $(PB_OBJS) $(OBJS)

certificates: $(CA_CRT) $(CRTS) $(PFXS)

test: $(PB_OBJS) $(PB_SRCS)
	go test ./...

clean:
	-rm $(PB_OBJS) $(OBJS)

distclean: clean
	-rm $(PFXS) $(CRTS) $(CSRS) $(KEYS) $(SRLS)
	-rm -rf $(DIST)

# Rules

$(CERT_ROOT):
	mkdir -p $(CERT_ROOT)

$(DIST):
	mkdir -p $(DIST)

$(DIST)/%: ./cmd/% cmd/%/*.go | $(DIST)
	go build -o $@ ./$<

$(PB_OBJS): $(PB_SRCS)
	protoc -I chat/ $^ --go_out=chat --go_opt=paths=source_relative --go-grpc_out=chat --go-grpc_opt=paths=source_relative

%.key: | $(CERT_ROOT)
	openssl ecparam -out $@ -name prime256v1 -genkey

%.csr: %.key | $(CERT_ROOT)
	openssl req -new -sha256 -key $^ -out $@ -subj '/CN=localhost'

$(CA_CRT): $(CA_KEY) | $(CERT_ROOT)
	openssl req -x509 -new -nodes -key $^ -sha256 -days 365 -out $@ -subj '/CN=localhost'
	openssl x509 -in $@ -text -noout

%.crt: %.csr %.key $(CA_KEY) $(CA_CRT) | $(CERT_ROOT)
	openssl x509 -req -in $< -out $@ -sha256 -days 365 -CA $(CA_CRT) -CAkey $(CA_KEY) -CAcreateserial -extfile $(CSR_CONF)
	openssl x509 -in $@ -text -noout

%.pfx: | $(CERT_ROOT)
	openssl pkcs12 -inkey $(@:%.pfx=%.key) -in $(@:%.pfx=%.crt) -export -nodes -passout pass: -out $@
