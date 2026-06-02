VERSION  := $(shell cat VERSION)
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE     := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS  := -s -w -X github.com/aiwen/aw-cli/internal/build.Version=$(VERSION) -X github.com/aiwen/aw-cli/internal/build.Commit=$(COMMIT) -X github.com/aiwen/aw-cli/internal/build.Date=$(DATE)

.PHONY: test fmt build build-all clean version patch minor major

test:
	go test ./...

fmt:
	gofmt -w .

build:
	go build -ldflags "$(LDFLAGS)" -o aw-cli .

build-all: build-linux build-windows build-darwin

build-linux:
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/aw-cli_$(VERSION)_linux_amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/aw-cli_$(VERSION)_linux_arm64 .

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/aw-cli_$(VERSION)_windows_amd64.exe .
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/aw-cli_$(VERSION)_windows_arm64.exe .

build-darwin:
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/aw-cli_$(VERSION)_darwin_amd64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/aw-cli_$(VERSION)_darwin_arm64 .

clean:
	rm -rf dist/ aw-cli aw-cli.exe

version:
	@echo $(VERSION)

patch:
	@bash scripts/version.sh patch

minor:
	@bash scripts/version.sh minor

major:
	@bash scripts/version.sh major
