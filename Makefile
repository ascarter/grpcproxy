GOROOT ?= $(shell go env GOROOT)

DIST    := dist
OBJS    := $(addprefix $(DIST)/,$(notdir $(wildcard cmd/*)))
PB_SRCS := chat/*.proto
PB_OBJS := chat/*.pb.go
CERTS   := cert.pem
KEYS    := $(CERTS:%cert.pem=%key.pem)

OPENSSL_ARGS := -newkey rsa:2048 -new -nodes -x509 -days 3650 -subj "/C=US/ST=Washington/L=Snoqualmie/O=$(USER)/OU=Development/CN=localhost"

$(DIST)/%: ./cmd/% cmd/%/*.go | $(DIST)
	go build -o $@ ./$<

$(DIST):
	mkdir -p $(DIST)

%.pem:
	openssl req $(OPENSSL_ARGS) -out $@ -keyout $(@:%cert.pem=%key.pem)

$(PB_OBJS): $(PB_SRCS)
	protoc -I chat/ $^ --go_out=chat --go_opt=paths=source_relative --go-grpc_out=chat --go-grpc_opt=paths=source_relative

.DEFAULT_GOAL = all

all: build

build: $(PB_OBJS) $(OBJS) $(CERTS)

test: $(PB_OBJS) $(PB_SRCS)
	go test ./...

clean:
	-rm $(PB_OBJS) $(OBJS)

distclean: clean
	-rm $(CERTS) $(KEYS)
	-rm -rf $(DIST)
