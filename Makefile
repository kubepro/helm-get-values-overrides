build: get-values-overrides

get-values-overrides: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o get-values-overrides

# linting
LINTER              := golangci-lint
LINTER_CONFIG       := .golangci.yaml

.PHONY: lint
lint:
lint:
	$(LINTER) run --config $(LINTER_CONFIG)
#	@git --no-pager show --check
