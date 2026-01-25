# Variáveis
BINARY_NAME=mj-cli
BUILD_DIR=build
MAIN_FILE=main.go

# Versão (pode ser sobrescrita via variável de ambiente)
VERSION?=dev

# Flags de build
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION)"

.PHONY: all build build-linux build-windows test clean

# Build padrão (sistema atual)
all: build

# Build para o sistema atual
build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Build para Linux (amd64)
build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)

# Build para Linux (arm64)
build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)

# Build para Windows (amd64)
build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)

# Build para Windows (arm64)
build-windows-arm64:
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_FILE)

# Build para todas as plataformas
build-all: build-linux build-linux-arm64 build-windows build-windows-arm64

# Executar testes unitários
test:
	go test ./... -v

# Executar testes com cobertura
test-coverage:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Limpar artefatos de build
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Instalar dependências
deps:
	go mod download
	go mod tidy

# Verificar código (lint)
lint:
	go vet ./...

# Executar a aplicação
run:
	go run $(MAIN_FILE)
