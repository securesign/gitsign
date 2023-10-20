
GIT_VERSION ?= $(shell git describe --tags --always --dirty)

LDFLAGS=-buildid= -X github.com/sigstore/gitsign/pkg/version.gitVersion=$(GIT_VERSION)

.PHONY: 
cross-platform: gitsign-cli-darwin-arm64 gitsign-cli-darwin-amd64 gitsign-cli-linux-amd64 gitsign-cli-linux-arm64 gitsign-cli-linux-ppc64le gitsign-cli-linux-s390x gitsign-cli-windows ## Build all distributable (cross-platform) binaries

.PHONY:	gitsign-cli-darwin-arm64
gitsign-cli-darwin-arm64: ## Build for mac M1
	env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o gitsign_cli_darwin_arm64 -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-darwin-amd64
gitsign-cli-darwin-amd64:  ## Build for Darwin (macOS)
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o gitsign_cli_darwin_amd64 -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-linux-amd64 
gitsign-cli-linux-amd64: ## Build for Linux amd64
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gitsign_cli_linux_amd64 -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-linux-arm64
gitsign-cli-linux-arm64: ## Build for Linux arm64
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o gitsign_cli_linux_arm64 -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-linux-ppc64le
gitsign-cli-linux-ppc64le: ## Build for Linux ppc64le
	env CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -o gitsign_cli_linux_ppc64le -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-linux-s390x
gitsign-cli-linux-s390x:  ## Build for Linux s390x
	env CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -o gitsign_cli_linux_s390x -trimpath -ldflags "$(LDFLAGS) -w -s" .

.PHONY: gitsign-cli-windows
gitsign-cli-windows: ## Build for Windows
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gitsign_cli_windows_amd64.exe -trimpath -ldflags "$(LDFLAGS) -w -s" .