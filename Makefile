GOROOT ?= $(shell go env GOROOT)

DIST    := dist
OBJS    := $(DIST)/echo \
		   $(DIST)/hello \
		   $(DIST)/proxy \
		   $(DIST)/server
PB_SRCS := chat/*.proto
PB_OBJS := chat/*.pb.go
CERT    := cert.pem
KEYFILE := key.pem
MK_CERT := $(GOROOT)/src/crypto/tls/generate_cert.go

$(DIST)/%: %/*.go | $(DIST)
	go build -o $@ ./$<

$(DIST):
	mkdir -p $(DIST)

$(CERT):
	go run $(MK_CERT) -host localhost

$(PB_OBJS): $(PB_SRCS)
	protoc -I chat/ $^ --go_out=chat --go_opt=paths=source_relative --go-grpc_out=chat --go-grpc_opt=paths=source_relative

.DEFAULT_GOAL = all

all: build

build: $(PB_OBJS) $(OBJS) $(CERT)

clean:
	-rm $(PB_OBJS) $(OBJS)

distclean: clean
	-rm $(CERT) $(KEYFILE)
	-rm -rf $(DIST)
