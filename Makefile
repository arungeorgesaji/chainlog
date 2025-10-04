build-cli:
	go build -o bin/chainlog-cli cmd/chainlog-cli/main.go

build-test:
	go build -o bin/chainlog-test cmd/chainlog-test/main.go

build-all: build build-test

run-cli:
	go run cmd/chainlog-cli/main.go

run-test:
	go run cmd/chainlog-test/main.go

clean:
	rm -rf bin/
