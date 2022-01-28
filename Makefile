default: test lint
ci: ci-test lint
.PHONY: default ci


test:
	go test -v ./...
.PHONY: test

ci-test:
	go test ./...
.PHONY: ci-test

lint:
	go vet ./...
.PHONY: lint
