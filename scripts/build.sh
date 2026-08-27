#!/usr/bin/env bash
set -euo pipefail
VERSION="${1:-dev}"
LDFLAGS="-ldflags=-s -w -X main.version=${VERSION}"
mkdir -p dist
TARGETS=(
	"windows amd64 .exe"
	"windows arm64 .exe"
	"linux amd64 "
	"linux arm64 "
	"darwin amd64 "
	"darwin arm64 "
)
for t in "${TARGETS[@]}"; do
	read -r GOOS GOARCH EXT <<<"$t"
	OUT="dist/docmost-mcp-${GOOS}-${GOARCH}${EXT}"
	echo "Building $OUT..."
	CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
		go build ${LDFLAGS} -o "$OUT" ./cmd/docmost-mcp
done
echo "Done. Artifacts:"
ls -lh dist/
