BINARY=dat
VERSION=$(shell git describe --tags)
BUILD=$(shell date +%FT%T%z)

LDFLAGS=-ldflags "-w -s -X github.com/Setheck/dat/cmd.Application=${BINARY} -X github.com/Setheck/dat/cmd.Version=${VERSION} -X github.com/Setheck/dat/cmd.Build=${BUILD}"

build: clean
	go build ${LDFLAGS} -o ${BINARY} .

test: build
	go test ./... -race -cover

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean