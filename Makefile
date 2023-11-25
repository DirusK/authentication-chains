.PHONY: all proto build clean test

BIN_NAME = ac
GO_FILES = $(shell find . -name '*.go')
PROTO_FILES = $(shell find ./proto -name '*.proto')

all: build

proto: $(PROTO_FILES)
	protoc -I=proto --go_out=. --go_opt=paths=import \
        --go-grpc_out=. --go-grpc_opt=paths=import \
        ./proto/*.proto

build: proto $(GO_FILES)
	go build -o $(BIN_NAME) main.go

clean:
	rm -rf ./bin
	rm -rf ./pkg/api/*.pb.go

test:
	go test ./...

run:
	go run main.go

start-alice:
	go run . start -c configs/nodes/alice.yaml

start-bob:
	go run . start -c configs/nodes/bob.yaml

start-tom:
	go run . start -c configs/nodes/tom.yaml

keygen:
	go run . keygen finn

dar:
	go run . send-dar finn

