BINARY=dat
VERSION=$(shell git describe --tags)
BUILD=$(shell date +%FT%T%z)
BASE_PKG:=main

CGO_ENABLED?=0
LDFLAGS=-ldflags "-extldflags '-static' -w -s \
				-X ${BASE_PKG}.Application=${BINARY} \
				-X ${BASE_PKG}.Version=${VERSION} \
				-X ${BASE_PKG}.Build=${BUILD}"

build:
	echo "CGO_ENABLED=$(CGO_ENABLED)"
	@echo "building ${BINARY} version:${VERSION} build:${BUILD}"
	@go build -a ${LDFLAGS} -o ${BINARY} .

test:
	go test ./... -cover

install: test build
	go install ${LDFLAGS} .

release: RELEASE_VERSION=v$(shell docker run --rm alpine/semver semver -c -i patch $(VERSION))
release:
	@echo "Creating Release: $(RELEASE_VERSION)"
	@echo "tagging..." && git tag -a $(RELEASE_VERSION) -m $(RELEASE_VERSION)
	@echo "pushing..." && git push origin $(RELEASE_VERSION)


clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install test dbuild deploy
