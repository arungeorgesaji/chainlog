build-cli:
	go build -o bin/chainlog-cli ./cmd/chainlog-cli/

build-test:
	go build -o bin/chainlog-test ./cmd/chainlog-test/

build-all: build build-test

run-cli:
	go run ./cmd/chainlog-cli/

run-test:
	go run ./cmd/chainlog-test/

clean:
	rm -rf bin/
