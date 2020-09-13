BINARY=dat
VERSION=$(shell git describe --tags)
BUILD=$(shell date +%FT%T%z)
BASE_PKG:=github.com/Setheck/dat/cmd
IMAGE:=setheck/dat

GOOS?=linux
CGO_ENABLED?=0
LDFLAGS=-ldflags "-extldflags '-static' -w -s \
				-X ${BASE_PKG}.Application=${BINARY} \
				-X ${BASE_PKG}.Version=${VERSION} \
				-X ${BASE_PKG}.Build=${BUILD}"

build: test
	echo "GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED)"
	go build -a ${LDFLAGS} -o ${BINARY} .

dbuild:
	# *Note, docker file calls `make build`
	docker build . -t ${IMAGE}:latest
	docker run --rm ${IMAGE}:latest --version

deploy: dbuild
	docker push ${IMAGE}:latest
	#TODO docker push ${IMAGE}:${VERSION}

test:
	go test ./... -cover

install: test
	go install ${LDFLAGS} .

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install test dbuild deploy
