BINARY=dat
VERSION=$(shell git describe --tags)
BUILD=$(shell date +%FT%T%z)

LDFLAGS=-ldflags "-w -s -X github.com/Setheck/dat/cmd.Application=${BINARY} -X github.com/Setheck/dat/cmd.Version=${VERSION} -X github.com/Setheck/dat/cmd.Build=${BUILD}"

build: test
	go build ${LDFLAGS} -o ${BINARY}

test:
	go test ./... -race -cover

install: test
	go install ${LDFLAGS} .

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean