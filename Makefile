VERSION  := $(shell cat VERSION)
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE     := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS  := -s -w -X github.com/aiwen/aw-cli/internal/build.Version=$(VERSION) -X github.com/aiwen/aw-cli/internal/build.Commit=$(COMMIT) -X github.com/aiwen/aw-cli/internal/build.Date=$(DATE)

.PHONY: test fmt build build-all archives release clean version patch minor major

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

# Package binaries into archives for npm install script (npx @aiwen/aw-cli)
archives: build-all
	cd dist && \
	rm -f checksums.txt && \
	for f in aw-cli_$(VERSION)_linux_* aw-cli_$(VERSION)_darwin_*; do \
		if [ -f "$$f" ]; then \
			base="$${f#aw-cli_$(VERSION)_}"; \
			os="$${base%_*}"; \
			arch="$${base##*_}"; \
			archive="aw-cli-$(VERSION)-$${os}-$${arch}.tar.gz"; \
			tmpdir="$$(mktemp -d)"; \
			cp "$$f" "$$tmpdir/aw-cli"; \
			tar -C "$$tmpdir" -czf "$$archive" aw-cli; \
			rm -rf "$$tmpdir"; \
		fi; \
	done && \
	for f in aw-cli_$(VERSION)_windows_*.exe; do \
		if [ -f "$$f" ]; then \
			base="$${f#aw-cli_$(VERSION)_}"; \
			base="$${base%.exe}"; \
			os="$${base%_*}"; \
			arch="$${base##*_}"; \
			archive="aw-cli-$(VERSION)-$${os}-$${arch}.zip"; \
			archive_path="$$PWD/$$archive"; \
			tmpdir="$$(mktemp -d)"; \
			cp "$$f" "$$tmpdir/aw-cli.exe"; \
			(cd "$$tmpdir" && zip "$$archive_path" aw-cli.exe); \
			rm -rf "$$tmpdir"; \
		fi; \
	done && \
	for archive in aw-cli-$(VERSION)-*.tar.gz aw-cli-$(VERSION)-*.zip; do \
		if [ -f "$$archive" ]; then \
			shasum -a 256 "$$archive" >> checksums.txt; \
		fi; \
	done

# Full release build (binaries + archives + checksums)
release: build-all archives

clean:
	rm -rf dist/ aw-cli aw-cli.exe bin/ node_modules/

version:
	@echo $(VERSION)

patch:
	@bash scripts/version.sh patch

minor:
	@bash scripts/version.sh minor

major:
	@bash scripts/version.sh major
