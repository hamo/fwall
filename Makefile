GOPATH:=$(GOPATH):$(CURDIR)

.PHONY: all client server goget

all: client server

client: goget
	cd client; GOPATH=$(GOPATH) CGO_ENABLED=0 go build -a -installsuffix cgo

server: goget
	cd server; GOPATH=$(GOPATH) CGO_ENABLED=0 go build -a -installsuffix cgo

goget:
	GOPATH=$(GOPATH) go get -d ./...

test: all
	GOPATH=$(GOPATH) go test ./...

fmt:
	go fmt ./...
