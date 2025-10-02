BINARY_NAME=chainlog
CLI_NAME=chainlog-cli
MINER_NAME=chainlog-miner

build:
	go build -o bin/$(BINARY_NAME) cmd/chainlog-node/main.go

build-cli:
	go build -o bin/$(CLI_NAME) cmd/chainlog-cli/main.go

build-miner:
	go build -o bin/$(MINER_NAME) cmd/chainlog-miner/main.go

build-all: build build-cli build-miner

run:
	go run cmd/chainlog-node/main.go

test:
	go test ./...

clean:
	rm -rf bin/
