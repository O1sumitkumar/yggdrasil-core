.PHONY: all tidy test vet fmt frontend daemon ci test-cluster package-headless run-daemon run-web

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

test:
	go test ./...

frontend:
	cd web && pnpm install && pnpm test && pnpm build

daemon:
	go build -o bin/yggdrasil-daemon ./cmd/daemon

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

ci: fmt vet test frontend
