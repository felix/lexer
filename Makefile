
pkgs	:= $(shell go list ./...)

.PHONY: lint test clean

ifdef GOPATH
GO111MODULE=on
endif

test: lint ## Run tests with coverage
	go test -short -cover -coverprofile coverage.out $(pkgs)
	go tool cover -html=coverage.out -o coverage.html

lint:
	golint $(pkgs)

clean: ## Clean all test files
	rm -rf coverage*
