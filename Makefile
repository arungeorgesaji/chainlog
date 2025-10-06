build-all: build-linux build-windows build-macos 

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/chainlog-cli-linux ./cmd/chainlog-cli/

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/chainlog-cli-windows.exe ./cmd/chainlog-cli/

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/chainlog-cli-macos-intel ./cmd/chainlog-cli/
	GOOS=darwin GOARCH=arm64 go build -o bin/chainlog-cli-macos-apple ./cmd/chainlog-cli/

build-cli:
	go build -o bin/chainlog-cli ./cmd/chainlog-cli/

clean:
	rm -rf bin/
