.PHONY: all proto build clean test

BIN_NAME = ac
GO_FILES = $(shell find . -name '*.go')
PROTO_FILES = $(shell find ./proto -name '*.proto')

all: build

proto: $(PROTO_FILES)
	protoc -I=proto --go-opt=paths=source-relative --go-grpc-opt=paths=source-relative --go_out=. -go_grpc_out=. ./proto/*.proto

build: proto $(GO_FILES)
	go build -o $(BIN_NAME) main.go

clean:
	rm -rf ./bin
	rm -rf ./pkg/api/*.pb.go

test:
	go test ./...

run:
	go run main.go