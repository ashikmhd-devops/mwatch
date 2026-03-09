BINARY=mwatch
VERSION=1.0.0
BUILD_DIR=./build

.PHONY: all build install clean run tidy

all: tidy build

tidy:
	go mod tidy

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY) .
	@echo "✓ Built $(BUILD_DIR)/$(BINARY)"

install: build
	cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/$(BINARY)
	@echo "✓ Installed to /usr/local/bin/$(BINARY)"

run:
	go run .

clean:
	rm -rf $(BUILD_DIR)

# Apple Silicon optimized build
build-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY)-arm64 .
	@echo "✓ Built Apple Silicon binary"
