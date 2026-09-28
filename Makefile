.PHONY: all tidy test vet fmt lint frontend daemon ci test-cluster package-headless run-daemon run-web screenshots

all: tidy test frontend

tidy:
	GOSUMDB=off go mod tidy

fmt:
	gofmt -w $$(find . -name '*.go' \
		-not -path './web/*' \
		-not -path './.gocache/*' \
		-not -path './vendor/*' \
		-not -path '*/node_modules/*')

vet:
	go vet ./...

lint:
	golangci-lint run ./...
	cd web && pnpm lint

test:
	go test ./...

frontend:
	cd web && pnpm install && pnpm test && pnpm build

daemon:
	go build -ldflags "-X github.com/yeixio/yggdrasil-core/internal/version.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo unknown)" -o bin/yggdrasil-daemon ./cmd/daemon
	go build -ldflags "-X github.com/yeixio/yggdrasil-core/internal/version.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo unknown)" -o bin/yggctl ./cmd/devctl

run-daemon: daemon
	./bin/yggdrasil-daemon

run-web:
	cd web && pnpm dev

test-cluster:
	chmod +x scripts/cluster-e2e.sh
	./scripts/cluster-e2e.sh

package-headless:
	chmod +x scripts/build/package-headless.sh
	./scripts/build/package-headless.sh

ci: fmt lint vet test frontend

screenshots:
	chmod +x scripts/capture-screenshots.sh
	./scripts/capture-screenshots.sh
