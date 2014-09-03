.PHONY: all client server gopath

all: client server

client: gopath
	cd client; go get -d; go build

server: gopath
	cd server; go get -d; go build

gopath:
	export GOPATH=`pwd`
