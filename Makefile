.PHONY: run run-cli test build clean

# Default target - run the API server
run:
	go run cmd/api/main.go

# Run the CLI version
run-cli:
	go run cmd/cli/main.go

# Run all tests
test:
	go test ./...

# Build the binaries
build:
	mkdir -p bin
	go build -o bin/scoundrel-api cmd/api/main.go
	go build -o bin/scoundrel-cli cmd/cli/main.go

# Clean built binaries
clean:
	rm -rf bin/