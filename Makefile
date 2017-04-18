PKG := github.com/dh1tw/gorigctl
COMMITID := $(shell git describe --always --long --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags)

PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
all: build

genproto:
	protoc --proto_path=./icd --gofast_out=./sb_radio ./icd/radio.proto
	protoc --proto_path=./icd --gofast_out=./sb_log ./icd/log.proto
	protoc --proto_path=./icd --gofast_out=./sb_ping ./icd/ping.proto
	protoc --proto_path=./icd --gofast_out=./sb_status ./icd/status.proto

build:	genproto	
	go build -v -ldflags="-X github.com/dh1tw/gorigctl/cmd.commitHash=${COMMIT} \
		-X github.com/dh1tw/gorigctl/cmd.version=${VERSION}"

# strip off dwraf table - used for travis CI
dist: genproto
	go build -v -ldflags="-w -X github.com/dh1tw/gorigctl/cmd.commitHash=${COMMIT} \
		-X github.com/dh1tw/gorigctl/cmd.version=${VERSION}"

# test:
# 	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

install: genproto 
	go install -v -ldflags="-w -X github.com/dh1tw/gorigctl/cmd.commitHash=${COMMIT} \
		-X github.com/dh1tw/gorigctl/cmd.version=${VERSION}"

install-deps:
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get -u ./...

# static: vet lint
# 	go build -i -v -o ${OUT}-v${VERSION} -tags netgo -ldflags="-extldflags \"-static\" -w -s -X main.version=${VERSION}" ${PKG}

clean:
	-@rm gorigctl gorigctl-v*

.PHONY: build install dist genproto vet lint clean install-deps