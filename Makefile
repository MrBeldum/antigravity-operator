BINARY_NAME=agyo
BUILD_DIR=bin
VERSION=0.1.0

.PHONY: all build test clean build-linux build-mac release

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/agyo
	@ln -sf $(BINARY_NAME) $(BUILD_DIR)/antigravity-operator
	@echo "✅ Binário compilado em $(BUILD_DIR)/$(BINARY_NAME) (e alias $(BUILD_DIR)/antigravity-operator)"

test:
	go test -v ./...

# Cross-compilação estática (CGO_ENABLED=0)
build-mac:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/agyo
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/agyo
	@echo "✅ Binários macOS gerados em $(BUILD_DIR)"

build-linux:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/agyo
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/agyo
	@echo "✅ Binários Linux estáticos gerados em $(BUILD_DIR)"

release: build-mac build-linux
	@echo "🚀 Todos os binários de release foram gerados com sucesso em $(BUILD_DIR)!"

install: build
	(cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME) && ln -sf /usr/local/bin/$(BINARY_NAME) /usr/local/bin/antigravity-operator) || \
	(mkdir -p $(HOME)/.local/bin && cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME) && ln -sf $(HOME)/.local/bin/$(BINARY_NAME) $(HOME)/.local/bin/antigravity-operator)
	@echo "✅ Instalado no PATH como 'agyo' e 'antigravity-operator'!"

clean:
	rm -rf $(BUILD_DIR)
