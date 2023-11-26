.PHONY: all proto build clean test

BIN_NAME = authentication-chains
GO_FILES = $(shell find . -name '*.go')
PROTO_FILES = $(shell find ./proto -name '*.proto')

CLIENT_NAME = finn
NODE_NAME = alice

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
	go run . node start -c configs/nodes/alice.yaml

start-bob:
	go run . node start -c configs/nodes/bob.yaml

start-tom:
	go run . node start -c configs/nodes/tom.yaml

keygen:
	go run . client keygen -n $(CLIENT_NAME)

send-dar:
	go run . client send-dar -n $(CLIENT_NAME)

get-blocks:
	go run . client get-blocks 0 100 -n $(CLIENT_NAME)

get-auth-table:
	go run . client get-auth-table -n $(CLIENT_NAME)

send-message:
	go run . client send-message "Hello world!" -n $(CLIENT_NAME)