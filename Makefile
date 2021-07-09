HAS_GOLANGCI := $(shell command -v golangci-lint;)
HAS_GODOC := $(shell command -v godoc;)

lint:
ifndef HAS_GOLANGCI
	$(error You must install github.com/golangci/golangci-lint)
endif
	@golangci-lint run -v -c .golangci.yml && echo "Lint OK"

test:
	@go test -timeout 30s -short -v -race -cover -coverprofile=coverage.out && echo "Test OK"

test-integration:
	@SYPL_TEST_MODE=integration go test -timeout 30s -v -race -cover -coverprofile=coverage.out -run Integration && echo "Test OK"

coverage:
	@go tool cover -func=coverage.out

doc:
ifndef HAS_GODOC
	$(error You must install godoc, run "go get golang.org/x/tools/cmd/godoc")
endif
	@echo "Open localhost:6060/pkg/github.com/saucelabs/sypl/ in your browser\n"
	@godoc -http :6060

ci: lint test coverage cobertura
ci-integration: lint test-integration coverage cobertura

.PHONY: lint test coverage cobertura ci ci-integration
