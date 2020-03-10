
.PHONY: test
test: lint
	go test -short -cover -coverprofile coverage.txt ./... \
		&& go tool cover -html=coverage.txt -o coverage.html

.PHONY: lint
lint:
	go vet ./...

.PHONY: clean
clean:
	rm -rf coverage*
