.PHONY: all proto build clean test

BIN_NAME = ac
GO_FILES = $(shell find . -name '*.go')
PROTO_FILES = $(shell find ./proto -name '*.proto')

all: build

proto: $(PROTO_FILES)
	protoc -I=./proto --go_out=. ./proto/*.proto

build: proto $(GO_FILES)
	go build -o $(BIN_NAME) main.go

clean:
	rm -rf ./bin
	rm -rf ./pkg/api/*.pb.go

test:
	go test ./...