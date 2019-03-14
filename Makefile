
GOPATH ?=	$(HOME)/go
pkgs := 	$(shell go list ./...)

.PHONY: test
test: lint ## Run tests with coverage
	go test -short -cover -coverprofile coverage.txt $(pkgs)
	go tool cover -html=coverage.txt -o coverage.html

.PHONY: lint
lint: $(GOPATH)/bin/golint ## Run the code linter
	@golint $(pkgs)

$(GOPATH)/bin/golint:
	go get -u golang.org/x/lint/golint

.PHONY: clean
clean: ## Clean all test files
	@rm -rf coverage*

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) |sort \
		|awk 'BEGIN{FS=":.*?## "};{printf "\033[36m%-30s\033[0m %s\n",$$1,$$2}'
