
pkgs	:= $(shell go list ./...)

.PHONY: test
test: lint ## Run tests with coverage
	go test -race -short -cover -coverprofile coverage.out $(pkgs)
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint:
	golint $(pkgs)

.PHONY: clean
clean: ## Clean all test files
	rm -rf coverage*
