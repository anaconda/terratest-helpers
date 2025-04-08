.PHONY: install
install:
	brew install pre-commit go
	pre-commit install

# Directly called by the CI
.PHONY: setup
setup:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2

# Called by users for local setup
.PHONY: init
init: install setup
