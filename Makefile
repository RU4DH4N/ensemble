PROTO_DIR = pkg/protocol/v1
PROTO_SRC = $(PROTO_DIR)/ensemble.proto
PROTO_GEN = $(PROTO_DIR)/*.pb.go
BIN_DIR = bin

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

PROTOC = protoc
GO = go

.PHONY: all proto build clean help test lint install-tools check deps

all: proto build

proto:
	@echo "Compiling Protobuf..."
	@$(PROTOC) --proto_path=$(PROTO_DIR) \
	   --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
	   --go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
	   $(PROTO_SRC)

build: $(BIN_DIR)/ensembled

$(BIN_DIR)/ensembled: $(shell find cmd/ensembled -name '*.go' 2>/dev/null) $(PROTO_GEN)
	@echo "Building ensembled..."
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(LDFLAGS) -o $@ ./cmd/ensembled

test:
	@echo "Running tests..."
	@$(GO) test -v -race -coverprofile=coverage.out ./...

test-coverage: test
	@$(GO) tool cover -html=coverage.out

bench:
	@echo "Running benchmarks..."
	@$(GO) test -bench=. -benchmem ./...

lint:
	@echo "Running linters..."
	@golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@goimports -w .

vet:
	@echo "Running go vet..."
	@$(GO) vet ./...

check: fmt vet lint

deps:
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy

install-tools:
	@echo "Installing tools..."
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

run: build
	@./$(BIN_DIR)/ensembled

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@rm -f $(PROTO_GEN)
	@rm -f coverage.out

clean-all: clean
	@$(GO) clean -cache -modcache -testcache