.PHONY: help
help:
	@echo "==> describe make commands"
	@echo ""
	@echo "build  ==> build binary for current GOOS and GOARCH"
	@echo "test   ==> run tests"

.PHONY: build
build:
	@goreleaser build --single-target --clean --snapshot

.PHONY: test
test:
	@go clean --testcache && go test -v ./...

