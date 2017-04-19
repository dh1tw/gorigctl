#!/bin/bash

env GOOS=$GIMME_OS GOARCH=$GIMME_ARCH go get github.com/golang/protobuf/proto
env GOOS=$GIMME_OS GOARCH=$GIMME_ARCH go get github.com/gogo/protobuf/protoc-gen-gofast
env GOOS=$GIMME_OS GOARCH=$GIMME_ARCH go get -d ./...

if [[ $GIMME_OS == 'darwin' ]]; then
    brew update
    brew outdated pkg-config || brew upgrade pkg-config
    brew install hamlib
    brew install protobuf
else #Linux
    # Ubuntu 16.04 comes with an old version of protobuf. 
    # We have to download and install a newer one
    ./ci/install-protobuf.sh
fi