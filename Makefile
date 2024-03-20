.PHONY: help
help:
	@echo "==> describe make commands"
	@echo ""
	@echo "build  ==> build bgps for current GOOS and GOARCH"
	@echo "test   ==> run tests"

.PHONY: build
build:
	# TODO: temporary hack return 0 for first release
	@goreleaser build --single-target --clean --snapshot && exit 0

.PHONY: test
test:
	@go clean --testcache && go test -v ./...

