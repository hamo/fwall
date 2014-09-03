GOPATH:=$(GOPATH):$(CURDIR)

.PHONY: all client server goget

all: client server

client: goget
	cd client; GOPATH=$(GOPATH) go build

server: goget
	cd server; GOPATH=$(GOPATH) go build

goget:
	GOPATH=$(GOPATH) go get -d ./...

test: all
	GOPATH=$(GOPATH) go test ./...
