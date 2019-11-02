# This how we want to name the binary output
BINARY=dat

# These are the values we want to pass for VERSION and BUILD
# git tag 1.0.1
# git commit -am "One more change after the tags"
VERSION=$(shell git describe --tags)
BUILD=$(shell date +%FT%T%z)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-w -s -X github.com/Setheck/dat/cmd.Version=${VERSION} -X github.com/Setheck/dat/cmd.Build=${BUILD}"

# Builds the project
build:
	go build ${LDFLAGS} -o ${BINARY} .

# Installs our project: copies binaries
install:
	go install ${LDFLAGS_f1}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install