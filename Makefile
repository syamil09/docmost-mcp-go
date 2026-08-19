VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: build build-all test lint clean release

build:
	go build $(LDFLAGS) -o dist/docmost-mcp-$$(go env GOOS)-$$(go env GOARCH)$$(go env GOEXE) ./cmd/docmost-mcp

build-all:
	@./scripts/build.sh $(VERSION)

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf dist/

release: build-all
	@cd dist && sha256sum docmost-mcp-* > SHA256SUMS
	@cat dist/SHA256SUMS
