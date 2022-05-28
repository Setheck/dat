BINARY=dat
VERSION=$(shell git describe --tags)
BUILT=$(shell date +%FT%T%z)
BUILD_PKG:=github.com/Setheck/dat/pkg/build


CGO_ENABLED?=0
LDFLAGS=-ldflags "-extldflags '-static' -w -s \
				-X ${BUILD_PKG}.Application=${BINARY} \
				-X ${BUILD_PKG}.Version=${VERSION} \
				-X ${BUILD_PKG}.Built=${BUILT}"

build:
	@echo "CGO_ENABLED=$(CGO_ENABLED)"
	@echo "building ${BINARY} version:${VERSION} built:${BUILT}"
	@cd cmd/dat && go build -a ${LDFLAGS} -o ../../bin/${BINARY} .

test:
	go test ./... -cover

install: test build
	@echo "installing ${BINARY}"
	@cd cmd/dat && go install ${LDFLAGS} .

release: RELEASE_VERSION=v$(shell docker run --rm alpine/semver semver -c -i patch $(VERSION))
release:
	@echo "Creating Release: $(RELEASE_VERSION)"
	@echo "tagging..." && git tag -a $(RELEASE_VERSION) -m $(RELEASE_VERSION)
	@echo "pushing..." && git push origin $(RELEASE_VERSION)


clean:
	@if [ -d bin ]; then echo "Removing bin/..."; rm -rf bin; fi

.PHONY: clean install test dbuild deploy
