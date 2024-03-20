.PHONY: help
help:
	@echo "==> describe make commands"
	@echo ""
	@echo "build  ==> build bgps for current GOOS and GOARCH"
	@echo "test   ==> run tests"

.PHONY: build
build:
	@goreleaser build --single-target --clean --snapshot

.PHONY: test
test:
	# TODO: temporary hack return 0 for first release
	@go clean --testcache && go test -v ./... && exit 0

