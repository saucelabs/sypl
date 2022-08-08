HAS_GOLANGCI := $(shell command -v golangci-lint;)
HAS_GODOC := $(shell command -v godoc;)

define PRINT_HELP_PYSCRIPT
import re, sys

for line in sys.stdin:
	match = re.match(r'^([0-9a-zA-Z_-]+):.*?## (.*)$$', line)
	if match:
		target, help = match.groups()
		print("%-20s %s" % (target, help))
endef
export PRINT_HELP_PYSCRIPT

help:  ## prints (this) help message
	@python -c "$$PRINT_HELP_PYSCRIPT" < $(MAKEFILE_LIST)

lint:  ## lint the code
ifndef HAS_GOLANGCI
	$(error You must install github.com/golangci/golangci-lint)
endif
	@golangci-lint run -v -c .golangci.yml && echo "Lint OK"

test:  ## run unit tests
	@go test -timeout 30s -short -v -race -cover -coverprofile=coverage.out ./...

test-integration:  ## run integration tests
	@SYPL_TEST_MODE=integration go test -timeout 30s -v -race -cover -coverprofile=coverage.out -run Integration && echo "Test OK"

coverage:  ## generate test coverage report at coverage.out
	@go tool cover -func=coverage.out

doc:  ## run godoc server
ifndef HAS_GODOC
	$(error You must install godoc, run "go install golang.org/x/tools/cmd/godoc@latest")
endif
	@echo "open http://localhost:6060/pkg/github.com/saucelabs/sypl/ in your browser\n"
	@godoc -http :6060

ci: lint test coverage  ## lint, test and coverage combined
ci-integration: lint test-integration coverage  ## lint, test-integration and coverage combined

.PHONY: lint test coverage ci ci-integration
