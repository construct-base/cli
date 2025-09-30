.PHONY: build install clean test

# Build the CLI
build:
	@echo "Building construct CLI..."
	@go build -o construct-bin main.go
	@echo "✅ Built: ./construct-bin"

# Install globally
install:
	@echo "Installing construct CLI..."
	@go build -o construct-bin main.go
	@sudo mv construct-bin /usr/local/bin/construct
	@echo "✅ Installed to /usr/local/bin/construct"
	@echo "Run 'construct --help' to get started"

# Clean build artifacts
clean:
	@rm -f construct-bin construct/main
	@echo "✅ Cleaned build artifacts"

# Run tests
test:
	@go test -v ./...

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=darwin GOARCH=amd64 go build -o dist/construct-darwin-amd64 main.go
	@GOOS=darwin GOARCH=arm64 go build -o dist/construct-darwin-arm64 main.go
	@GOOS=linux GOARCH=amd64 go build -o dist/construct-linux-amd64 main.go
	@GOOS=windows GOARCH=amd64 go build -o dist/construct-windows-amd64.exe main.go
	@echo "✅ Built binaries in dist/"