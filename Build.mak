
GIT_VERSION ?= $(shell git describe --tags --always --dirty)

GIT_HASH ?= $(shell git rev-parse HEAD)
DATE_FMT = +%Y-%m-%dT%H:%M:%SZ
SOURCE_DATE_EPOCH ?= $(shell git log -1 --no-show-signature --pretty=%ct)
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
GIT_TREESTATE = "clean"
DIFF = $(shell git diff --quiet >/dev/null 2>&1; if [ $$? -eq 1 ]; then echo "1"; fi)
ifeq ($(DIFF), 1)
    GIT_TREESTATE = "dirty"
endif

LDFLAGS=-buildid= -X github.com/sigstore/gitsign/pkg/version.gitVersion=$(GIT_VERSION)

.PHONY: 
cross-platform: gitsign-cli-darwin-arm64 gitsign-cli-darwin-amd64 gitsign-cli-linux-amd64 gitsign-cli-linux-arm64 gitsign-cli-linux-ppc64le gitsign-cli-linux-s390x gitsign-cli-windows ## Build all distributable (cross-platform) binaries

.PHONY:	gitsign-cli-darwin-arm64
gitsign-cli-darwin-arm64: ## Build for mac M1
	env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -mod=readonly -o gitsign_cli_darwin_arm64 -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-darwin-amd64
gitsign-cli-darwin-amd64:  ## Build for Darwin (macOS)
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod=readonly -o gitsign_cli_darwin_amd64 -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-linux-amd64 
gitsign-cli-linux-amd64: ## Build for Linux amd64
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -o gitsign_cli_linux_amd64 -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-linux-arm64
gitsign-cli-linux-arm64: ## Build for Linux arm64
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod=readonly -o gitsign_cli_linux_arm64 -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-linux-ppc64le
gitsign-cli-linux-ppc64le: ## Build for Linux ppc64le
	env CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -mod=readonly -o gitsign_cli_linux_ppc64le -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-linux-s390x
gitsign-cli-linux-s390x:  ## Build for Linux s390x
	env CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -mod=readonly -o gitsign_cli_linux_s390x -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-windows
gitsign-cli-windows: ## Build for Windows
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -mod=readonly -o gitsign_cli_windows_amd64.exe -trimpath -ldflags "$(LDFLAGS) -w -s" .